package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Client struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	backendURL := os.Getenv("BACKEND_URL") + "/register-charger"

	charger := Client{
		Name: "Charger1",
	}

	data, err := json.Marshal(charger)
	if err != nil {
		log.Fatalf("Error marshalling charger: %v", err)
	}

	resp, err := http.Post(backendURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatalf("Error registering charger: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to register charger: received status code %d", resp.StatusCode)
	}

	fmt.Println("Charger registered successfully")
}
