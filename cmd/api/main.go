package main

import (
	"bake_backend/internal/api"
	"bake_backend/internal/config"
	"bake_backend/internal/repository"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	connStr := fmt.Sprintf("postgres://%s:%s@dpg-d1c7auje5dus73f9n0bg-a.oregon-postgres.render.com:5432/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	log.Println("Migrations applied successfully")

	repo := repository.NewPostgresRepository(db)
	handler := api.NewHandler(repo)

	r := mux.NewRouter()

	// Updated routes - removed report dates endpoints
	r.HandleFunc("/api/reports/available-dates", handler.GetAvailableDates).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/reports/marketing-sources", handler.GetMarketingSources).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/reports/sales-teams", handler.GetSalesTeams).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/reports/marketing-data", handler.GetMarketingData).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/reports/sales-data", handler.GetSalesData).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/reports/marketing-data", handler.SaveMarketingData).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/reports/sales-data", handler.SaveSalesData).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/reports/marketing-data/{id}", handler.UpdateMarketingData).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/reports/sales-data/{id}", handler.UpdateSalesData).Methods("PUT", "OPTIONS")

	loggedRouter := withLogging(r)
	corsRouter := withCORS(loggedRouter)

	log.Printf("Server starting on :8080")
	if err := http.ListenAndServe(":8080", corsRouter); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// Middleware to allow CORS from any origin
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
