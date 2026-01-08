package model

import (
	"errors"
	"strings"
)

// Book represents the main data structure for a book
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

// CreateBookRequest defines what we expect when creating a book
type CreateBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

// validate checks if the data for a new book is correct
func (r CreateBookRequest) Validate() error {
	// check if title is empty
	if strings.TrimSpace(r.Title) == "" {
		return errors.New("le titre est obligatoire")
	}

	// check if author is empty
	if strings.TrimSpace(r.Author) == "" {
		return errors.New("l'auteur est obligatoire")
	}

	// check if the year is in a coherent range
	//
	if r.Year < 0 || r.Year > 2026 {
		return errors.New("l'annee doit être comprise entre 0 et 2026")
	}

	return nil
}

// defines what can be updated
type UpdateBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

// validate checks the update data
func (r UpdateBookRequest) Validate() error {
	// for an update we still want valid data if provided
	if strings.TrimSpace(r.Title) == "" {
		return errors.New("le titre ne peut pas être vide lors de la modification")
	}
	if strings.TrimSpace(r.Author) == "" {
		return errors.New("l'auteur ne peut pas être vide lors de la modification")
	}
	return nil
}
