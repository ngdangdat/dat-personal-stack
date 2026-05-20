package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/ngdangdat/dat-personal-stack/backend/services"
)

type HealthResponse struct {
	Status string `json:"status"`
}

// SetupRouter configures all HTTP routes on a ServeMux
func SetupRouter(githubService *services.GitHubService) *http.ServeMux {
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

	// PRs dashboard data fetching
	mux.HandleFunc("/api/github/prs", enableCors(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		token := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else if strings.HasPrefix(authHeader, "token ") {
			token = strings.TrimPrefix(authHeader, "token ")
		}
		token = strings.TrimSpace(token)

		if token == "" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Read query parameters
		queryType := r.URL.Query().Get("type") // "reviewing" or "mine"
		username := r.URL.Query().Get("username")

		if username == "" {
			http.Error(w, "Missing 'username' query parameter", http.StatusBadRequest)
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
	mux := SetupRouter(githubService)

	log.Println("Starting backend server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
