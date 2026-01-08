package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"ProjetAPI-GO-CRUD/internal/handler"
	"ProjetAPI-GO-CRUD/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv" // Load the .env file
	_ "github.com/lib/pq"      // Database driver
)

func main() {
	// Load the .env file that contains our secret configuration
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found, using system variables")
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
		log.Fatal("Could not open database connection:", err)
	}
	defer db.Close()

	// Verify database connection with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal("Database is not responding:", err)
	}
	fmt.Println("Database connection verified!")

	// INITIALIZE LAYERS
	bookRepo := repository.NewPostgresRepository(db)
	bookHandler := handler.NewBookHandler(bookRepo)

	// ROUTER SETUP
	r := chi.NewRouter()

	// Health check route
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status": "OK", "database": "connected"}`)
	})

	// BOOK ROUTES
	r.Route("/books", func(r chi.Router) {
		r.Get("/", bookHandler.ListBooks)   // GET /books
		r.Post("/", bookHandler.CreateBook) // POST /books

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", bookHandler.GetBook)       // GET /books/{id}
			r.Put("/", bookHandler.UpdateBook)    // Update (Point 19 & 20)
			r.Delete("/", bookHandler.DeleteBook) // Delete (Point 19 & 20)
		})
	})

	fmt.Printf("Server is running on port 8080. Connecting to host: %s\n", dbHost)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", r))
}
