package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error
	connStr := os.Getenv("DATABASE_URL")
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/clients", GetClients).Methods("GET")
	r.HandleFunc("/chargers", GetChargers).Methods("GET")
	r.HandleFunc("/register-client", RegisterClient).Methods("POST")
	r.HandleFunc("/register-charger", RegisterCharger).Methods("POST")

	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Backend starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
