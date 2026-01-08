package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"ProjetAPI-GO-CRUD/internal/repository"
)

// BookHandler link  API routes to the database logic
type BookHandler struct {
	repo repository.BookRepository
}

// NewBookHandler creates a new handler with its repository dependency
func NewBookHandler(repo repository.BookRepository) *BookHandler {
	return &BookHandler{repo: repo}
}

// GetBook handles GET /books/{id} and manages the 404 error
func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		h.sendJSONError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		h.sendJSONError(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	// Fetch data from repository
	book, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		h.sendJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Handle sql.ErrNoRows
	if book == nil {
		h.sendJSONError(w, "Book not found", http.StatusNotFound)
		return
	}

	// Send the book as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// ListBooks handles GET /books
func (h *BookHandler) ListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.repo.List(r.Context())
	if err != nil {
		h.sendJSONError(w, "Could not fetch books", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// helper function to send errors in JSON format
func (h *BookHandler) sendJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
