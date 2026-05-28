package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ngdangdat/dat-personal-stack/backend/services"
)

type HealthResponse struct {
	Status string `json:"status"`
}

// SetupRouter configures all HTTP routes on a ServeMux
func SetupRouter(githubService *services.GitHubService, dbService *services.DBService) *http.ServeMux {
	mux := http.NewServeMux()

	// CORS Middleware Helper
	enableCors := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			h(w, r)
		}
	}

	// Health check endpoint
	mux.HandleFunc("/api/health", enableCors(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
	}))

	// GET Settings endpoint
	mux.HandleFunc("/api/settings", enableCors(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if dbService == nil {
			json.NewEncoder(w).Encode(map[string]string{"username": "", "token": ""})
			return
		}

		cfg, err := dbService.GetLatestConfig()
		if err != nil {
			log.Printf("Error fetching latest config: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if cfg == nil {
			json.NewEncoder(w).Encode(map[string]string{"username": "", "token": ""})
			return
		}

		json.NewEncoder(w).Encode(cfg)
	}))

	// POST Settings endpoint
	mux.HandleFunc("/api/settings/save", enableCors(func(w http.ResponseWriter, r *http.Request) {
		// support both GET/POST just in case or keep POST for saving
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if dbService == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Database service is not initialized"})
			return
		}

		var req struct {
			Username string `json:"username"`
			Token    string `json:"token"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON body"})
			return
		}

		req.Username = strings.TrimSpace(req.Username)
		req.Token = strings.TrimSpace(req.Token)

		if req.Username == "" || req.Token == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Username and token are required"})
			return
		}

		if err := dbService.SaveConfig(req.Username, req.Token); err != nil {
			log.Printf("Error saving config to DB: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))

	// PRs dashboard data fetching
	mux.HandleFunc("/api/github/prs", enableCors(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read token from Authorization header
		authHeader := r.Header.Get("Authorization")
		var token string
		if authHeader != "" {
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			} else if strings.HasPrefix(authHeader, "token ") {
				token = strings.TrimPrefix(authHeader, "token ")
			}
			token = strings.TrimSpace(token)
		}

		// Read query parameters
		queryType := r.URL.Query().Get("type") // "reviewing" or "mine"
		username := r.URL.Query().Get("username")

		if username == "" {
			http.Error(w, "Missing 'username' query parameter", http.StatusBadRequest)
			return
		}

		// If token is missing, lookup config from database
		if token == "" && dbService != nil {
			cfg, err := dbService.GetConfig(username)
			if err != nil {
				log.Printf("Error fetching config from DB for %s: %v", username, err)
			} else if cfg != nil {
				token = cfg.Token
			}
		}

		if token == "" {
			http.Error(w, "Invalid token or not configured in backend database", http.StatusUnauthorized)
			return
		}

		if queryType != "reviewing" && queryType != "mine" {
			queryType = "reviewing" // default
		}

		data, err := githubService.FetchPRs(token, queryType, username)
		if err != nil {
			log.Printf("Error fetching PRs: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}))

	return mux
}

func main() {
	githubService := services.NewGitHubService()

	databaseURL := os.Getenv("DATABASE_URL")
	var dbService *services.DBService
	var err error

	if databaseURL != "" {
		log.Println("Initializing database connection...")
		dbService, err = services.NewDBService(databaseURL)
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}
		defer dbService.Close()
	} else {
		log.Println("DATABASE_URL is not set. Database features are disabled.")
	}

	mux := SetupRouter(githubService, dbService)

	log.Println("Starting backend server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
