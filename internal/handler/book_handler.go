package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"ProjetAPI-GO-CRUD/internal/model"
	"ProjetAPI-GO-CRUD/internal/repository"

	"github.com/go-chi/chi/v5"
)

// BookHandler links API routes to the database logic
type BookHandler struct {
	repo repository.BookRepository
}

// NewBookHandler creates a new handler with its repository dependency
func NewBookHandler(repo repository.BookRepository) *BookHandler {
	return &BookHandler{repo: repo}
}

// CreateBook handles POST /books
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req model.CreateBookRequest

	// Parse JSON input
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Data Validation
	if err := req.Validate(); err != nil {
		h.sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call repository to save in DB
	id, err := h.repo.Create(r.Context(), &req)
	if err != nil {
		h.sendJSONError(w, "Internal server error during creation", http.StatusInternalServerError)
		return
	}

	// Send 201 Created response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// GetBook handles GET /books/{id}
func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	// Get ID from chi URL parameter
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
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

	// Success response
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

	// This ensures an empty table [] is returned instead of null when no books exist
	if books == nil {
		books = []model.Book{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// UpdateBook handles PUT /books/{id}
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req model.UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		h.sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.repo.Update(r.Context(), id, &req)
	if err != nil {
		h.sendJSONError(w, "Could not update book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Book updated successfully"})
}

// DeleteBook handles DELETE /books/{id}
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.repo.Delete(r.Context(), id)
	if err != nil {
		h.sendJSONError(w, "Could not delete book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// sendJSONError helper for uniform error responses
func (h *BookHandler) sendJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
