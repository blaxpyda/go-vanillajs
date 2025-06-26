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

	// Setup routes
	// API routes
	http.HandleFunc("/api/houses/top", houseHandler.GetTopHouses)
	http.HandleFunc("/api/houses/", houseHandler.HandleHousesRoute)
	http.HandleFunc("/api/houses", houseHandler.HandleHousesRoute)
	http.HandleFunc("/api/agents", houseHandler.GetAgents)
	http.HandleFunc("/api/house-types", houseHandler.GetHouseTypes)

	// Static file handler (this should be last)
	http.Handle("/", http.FileServer(http.Dir("public")))

	fmt.Println("Serving the files")
	logInstance.Info("Server starting on :8080")

	const addr = ":8080"
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Server has failed: %v", err)
		logInstance.Error("Server has failed", err)
	}
}
