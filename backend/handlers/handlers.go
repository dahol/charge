package main

import (
	"encoding/json"
	"net/http"
	"sync"

	_ "github.com/lib/pq"
)

var mu sync.Mutex

type Client struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetClients(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	rows, err := db.Query("SELECT id, name FROM clients")
	if err != nil {
		http.Error(w, "Failed to fetch clients", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var clients []Client
	for rows.Next() {
		var client Client
		if err := rows.Scan(&client.ID, &client.Name); err != nil {
			http.Error(w, "Failed to scan client", http.StatusInternalServerError)
			return
		}
		clients = append(clients, client)
	}

	json.NewEncoder(w).Encode(clients)
}

func GetChargers(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	rows, err := db.Query("SELECT id, name FROM chargers")
	if err != nil {
		http.Error(w, "Failed to fetch chargers", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var chargers []Client
	for rows.Next() {
		var charger Client
		if err := rows.Scan(&charger.ID, &charger.Name); err != nil {
			http.Error(w, "Failed to scan charger", http.StatusInternalServerError)
			return
		}
		chargers = append(chargers, charger)
	}

	json.NewEncoder(w).Encode(chargers)
}

func RegisterClient(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var client Client
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := db.QueryRow("INSERT INTO clients(name) VALUES($1) RETURNING id", client.Name).Scan(&client.ID)
	if err != nil {
		http.Error(w, "Failed to register client", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(client)
}

func RegisterCharger(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var charger Client
	if err := json.NewDecoder(r.Body).Decode(&charger); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := db.QueryRow("INSERT INTO chargers(name) VALUES($1) RETURNING id", charger.Name).Scan(&charger.ID)
	if err != nil {
		http.Error(w, "Failed to register charger", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(charger)
}
