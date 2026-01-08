package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv" // Load the .env file
	_ "github.com/lib/pq"      // Database driver
)

func main() {

	// Load the .env file that contains our secret configuration
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: Fichier .env non trouvé, utilisation des variables système")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Format the connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	// Open the connection with database/sql + driver
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		// Stop the application if the configuration is invalid
		log.Fatal("Impossible d'ouvrir la connexion :", err)
	}

	// Automatically close resources when the program end
	defer db.Close()

	// We create a context with a 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if the database connection is actually alive
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal("La base de données ne répond pas :", err)
	}
	fmt.Println("Database connection verified!")

	// EXISTING ROUTES
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status": "OK", "database": "connected"}`)
	})

	fmt.Printf("Le serveur tourne sur le port 8080 et tente de joindre %s\n", dbHost)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", r))
}
