package main

import (
	"log"
	"net/http"
	"os"
	"pull-up-mania/internal/handlers"
	"pull-up-mania/internal/storage"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func main() {
	// Check if running as Cloud Function
	if os.Getenv("FUNCTION_TARGET") != "" {
		// Register the HTTP function for Cloud Functions
		functions.HTTP("PullUpMania", httpHandler)
		// Framework will handle the rest automatically
		return
	} else {
		// Run as regular HTTP server for local development
		server := setupServer()
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		log.Printf("Server starting on port %s", port)
		log.Printf("Visit: http://localhost:%s", port)
		log.Fatal(http.ListenAndServe(":"+port, server))
	}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	server := setupServer()
	server.ServeHTTP(w, r)
}

func setupServer() *http.ServeMux {
	// Initialize storage
	var store storage.Storage
	var err error

	// Try to use Cloud Firestore in production (when project ID is set)
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID != "" {
		// TODO: Implement Firestore storage for production
		log.Println("Warning: Firestore not implemented yet, falling back to SQLite")
		store, err = storage.NewSQLiteStorage("workouts.db")
	} else {
		// Use SQLite for local development
		store, err = storage.NewSQLiteStorage("workouts.db")
	}

	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize handlers
	h := handlers.NewHandlers(store)
	if err := h.LoadTemplates(); err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}

	// Setup routes
	mux := http.NewServeMux()

	// Static files (for local development)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Admin routes (register these FIRST to avoid conflict with "/" catch-all)
	mux.HandleFunc("/admin/login", h.AdminLoginHandler)
	mux.HandleFunc("/admin/logout", h.AdminLogoutHandler)
	mux.HandleFunc("/admin/settings", h.AdminSettingsHandler)
	mux.HandleFunc("/admin/workout-plan", h.AdminWorkoutPlanHandler)
	mux.HandleFunc("/admin", h.AdminDashboardHandler)

	// Main routes
	mux.HandleFunc("/history", h.HistoryHandler)
	mux.HandleFunc("/save-workout", h.SaveWorkoutHandler)

	// Home route (register LAST as it's a catch-all)
	mux.HandleFunc("/", h.HomeHandler)

	// Health check for Cloud Functions
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return mux
}
