package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"pull-up-mania/internal/storage"
	"time"
)

type Handlers struct {
	storage   storage.Storage
	templates *template.Template
}

func NewHandlers(storage storage.Storage) *Handlers {
	return &Handlers{
		storage: storage,
	}
}

func (h *Handlers) LoadTemplates() error {
	var err error
	h.templates, err = template.ParseGlob("internal/templates/*.html")
	return err
}

// generateID creates a random hex ID
func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// getDayOfWeek returns the day of week (0=Sunday) for today
func getDayOfWeek() int {
	return int(time.Now().Weekday())
}

// parseFormFloat safely parses a float from form values
func parseFormFloat(value string) float64 {
	if value == "" {
		return 0
	}
	// Simple parsing - in production you'd want better error handling
	var result float64
	if _, err := fmt.Sscanf(value, "%f", &result); err != nil {
		return 0
	}
	return result
}

// parseFormInt safely parses an int from form values
func parseFormInt(value string) int {
	if value == "" {
		return 0
	}
	var result int
	if _, err := fmt.Sscanf(value, "%d", &result); err != nil {
		return 0
	}
	return result
}
