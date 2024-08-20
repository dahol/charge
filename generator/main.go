package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Charger struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

var db *sql.DB

func main() {
	var err error
	connStr := os.Getenv("DATABASE_URL")
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error creating Kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	for {
		checkAndManageChargers(clientset)
		time.Sleep(30 * time.Second) // Poll every 30 seconds
	}
}

func checkAndManageChargers(clientset *kubernetes.Clientset) {
	rows, err := db.Query("SELECT id, name, status FROM charger")
	if err != nil {
		log.Printf("Error querying database: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var charger Charger
		if err := rows.Scan(&charger.ID, &charger.Name, &charger.Status); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		switch charger.Status {
		case "active":
			ensurePodRunning(clientset, charger)
		case "error":
			respawnPod(clientset, charger)
		case "disabled":
			deletePod(clientset, charger)
		case "inactive":
			// Do nothing
		default:
			log.Printf("Unknown status for charger %s: %s", charger.Name, charger.Status)
		}
	}
}

func ensurePodRunning(clientset *kubernetes.Clientset, charger Charger) {
	podName := fmt.Sprintf("charger-%d", charger.ID)
	_, err := clientset.CoreV1().Pods("golang-app").Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		if metav1.IsNoMatchError(err) {
			log.Printf("Creating pod for active charger: %s", charger.Name)
			createChargerPod(clientset, charger)
		} else {
			log.Printf("Error checking pod for charger %s: %v", charger.Name, err)
		}
	}
}

func respawnPod(clientset *kubernetes.Clientset, charger Charger) {
	log.Printf("Respawning pod for charger in error state: %s", charger.Name)
	deletePod(clientset, charger)
	createChargerPod(clientset, charger)
}

func deletePod(clientset *kubernetes.Clientset, charger Charger) {
	podName := fmt.Sprintf("charger-%d", charger.ID)
	err := clientset.CoreV1().Pods("golang-app").Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil && !metav1.IsNoMatchError(err) {
		log.Printf("Error deleting pod for charger %s: %v", charger.Name, err)
	} else {
		log.Printf("Deleted pod for charger: %s", charger.Name)
	}
}

func createChargerPod(clientset *kubernetes.Clientset, charger Charger) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("charger-%d", charger.ID),
			Labels: map[string]string{
				"app": "charger",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "charger",
					Image: "charger-image:latest",
					Env: []corev1.EnvVar{
						{
							Name:  "CHARGER_ID",
							Value: fmt.Sprintf("%d", charger.ID),
						},
						{
							Name:  "CHARGER_NAME",
							Value: charger.Name,
						},
					},
				},
			},
		},
	}

	_, err := clientset.CoreV1().Pods("golang-app").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Error creating pod for charger %s: %v", charger.Name, err)
	} else {
		log.Printf("Created pod for charger: %s", charger.Name)
	}
}
