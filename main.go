package main

import (
	"fmt"
	"log"
	"net/http"

	"thugcorp.io/nomado/db"
	"thugcorp.io/nomado/handlers"
	"thugcorp.io/nomado/logger"
	"thugcorp.io/nomado/repository"
)

func initializeLogger() *logger.Logger {
	logInstance, err := logger.NewLogger("nomado.log")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	return logInstance
}

func enableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r)
		if r.Method == "OPTIONS" {
			return
		}
		next(w, r)
	}
}

func main() {
	// Initialize the logger instance
	logInstance := initializeLogger()
	defer logInstance.Close()

	// Initialize database
	database, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
		logInstance.Error("Failed to initialize database", err)
	}
	defer database.Close()

	// Create tables and seed data
	if err := database.CreateTables(); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
		logInstance.Error("Failed to create tables", err)
	}

	if err := database.SeedData(); err != nil {
		log.Printf("Warning: Failed to seed data: %v", err)
		logInstance.Error("Failed to seed data", err)
	}

	// Initialize repositories
	houseRepo := repository.NewHouseRepository(database.DB)
	agentRepo := repository.NewAgentRepository(database.DB)
	houseTypeRepo := repository.NewHouseTypeRepository(database.DB)

	// Initialize handlers
	houseHandler := handlers.NewHouseHandler(houseRepo, agentRepo, houseTypeRepo, logInstance)

	// Setup API routes with CORS middleware
	http.HandleFunc("/api/houses/top", corsMiddleware(houseHandler.GetTopHouses))
	http.HandleFunc("/api/houses/", corsMiddleware(houseHandler.HandleHousesRoute))
	http.HandleFunc("/api/houses", corsMiddleware(houseHandler.HandleHousesRoute))
	http.HandleFunc("/api/agents", corsMiddleware(houseHandler.GetAgents))
	http.HandleFunc("/api/house-types", corsMiddleware(houseHandler.GetHouseTypes))

	// Health check endpoint
	http.HandleFunc("/api/health", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"nomado-api"}`))
	}))

	// API info endpoint
	http.HandleFunc("/api", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		apiInfo := `{
			"service": "Nomado Real Estate API",
			"version": "1.0.0",
			"endpoints": {
				"houses": "/api/houses",
				"top_houses": "/api/houses/top",
				"house_detail": "/api/houses/{id}",
				"agents": "/api/agents",
				"house_types": "/api/house-types",
				"health": "/api/health"
			}
		}`
		w.Write([]byte(apiInfo))
	}))

	fmt.Println("🚀 Nomado Real Estate API Server starting...")
	logInstance.Info("API Server starting on :8080")
	fmt.Printf("📡 API endpoints available at: http://localhost:8080/api\n")
	fmt.Printf("🔍 Health check: http://localhost:8080/api/health\n")

	const addr = ":8080"
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Server has failed: %v", err)
		logInstance.Error("Server has failed", err)
	}
}
