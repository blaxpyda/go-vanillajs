package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"thugcorp.io/nomado/logger"
	"thugcorp.io/nomado/models"
	"thugcorp.io/nomado/repository"
)

type HouseHandler struct {
	houseRepo     *repository.HouseRepository
	agentRepo     *repository.AgentRepository
	houseTypeRepo *repository.HouseTypeRepository
	logger        *logger.Logger
}

type HouseWithDetails struct {
	models.House
	Agent     *models.Agent     `json:"agent,omitempty"`
	HouseType *models.HouseType `json:"house_type,omitempty"`
}

func NewHouseHandler(houseRepo *repository.HouseRepository, agentRepo *repository.AgentRepository, houseTypeRepo *repository.HouseTypeRepository, logger *logger.Logger) *HouseHandler {
	return &HouseHandler{
		houseRepo:     houseRepo,
		agentRepo:     agentRepo,
		houseTypeRepo: houseTypeRepo,
		logger:        logger,
	}
}

func (h *HouseHandler) GetTopHouses(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get limit from query parameter, default to 10
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	houses, err := h.houseRepo.GetTopHouses(limit)
	if err != nil {
		h.logger.Error("Failed to get top houses", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Enrich houses with agent and house type details
	var housesWithDetails []HouseWithDetails
	for _, house := range houses {
		houseWithDetails := HouseWithDetails{House: house}

		// Get agent details
		if agent, err := h.agentRepo.GetAgentByID(house.AgentID); err == nil {
			houseWithDetails.Agent = agent
		}

		// Get house type details
		if houseType, err := h.houseTypeRepo.GetHouseTypeByID(house.HouseTypeID); err == nil {
			houseWithDetails.HouseType = houseType
		}

		housesWithDetails = append(housesWithDetails, houseWithDetails)
	}

	if err := json.NewEncoder(w).Encode(housesWithDetails); err != nil {
		h.logger.Error("Failed to encode response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *HouseHandler) GetAllHouses(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	houses, err := h.houseRepo.GetAllHouses()
	if err != nil {
		h.logger.Error("Failed to get all houses", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Enrich houses with agent and house type details
	var housesWithDetails []HouseWithDetails
	for _, house := range houses {
		houseWithDetails := HouseWithDetails{House: house}

		// Get agent details
		if agent, err := h.agentRepo.GetAgentByID(house.AgentID); err == nil {
			houseWithDetails.Agent = agent
		}

		// Get house type details
		if houseType, err := h.houseTypeRepo.GetHouseTypeByID(house.HouseTypeID); err == nil {
			houseWithDetails.HouseType = houseType
		}

		housesWithDetails = append(housesWithDetails, houseWithDetails)
	}

	if err := json.NewEncoder(w).Encode(housesWithDetails); err != nil {
		h.logger.Error("Failed to encode response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *HouseHandler) GetHouseByID(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	path := r.URL.Path
	idStr := path[len("/api/houses/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid house ID", http.StatusBadRequest)
		return
	}

	house, err := h.houseRepo.GetHouseByID(id)
	if err != nil {
		h.logger.Error("Failed to get house by ID", err)
		http.Error(w, "House not found", http.StatusNotFound)
		return
	}

	houseWithDetails := HouseWithDetails{House: *house}

	// Get agent details
	if agent, err := h.agentRepo.GetAgentByID(house.AgentID); err == nil {
		houseWithDetails.Agent = agent
	}

	// Get house type details
	if houseType, err := h.houseTypeRepo.GetHouseTypeByID(house.HouseTypeID); err == nil {
		houseWithDetails.HouseType = houseType
	}

	if err := json.NewEncoder(w).Encode(houseWithDetails); err != nil {
		h.logger.Error("Failed to encode response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *HouseHandler) CreateHouse(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var house models.House
	if err := json.NewDecoder(r.Body).Decode(&house); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := h.houseRepo.CreateHouse(&house); err != nil {
		h.logger.Error("Failed to create house", err)
		http.Error(w, "Failed to create house", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(house); err != nil {
		h.logger.Error("Failed to encode response", err)
	}
}

func (h *HouseHandler) UpdateHouse(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	path := r.URL.Path
	idStr := path[len("/api/houses/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid house ID", http.StatusBadRequest)
		return
	}

	var house models.House
	if err := json.NewDecoder(r.Body).Decode(&house); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	house.ID = id
	if err := h.houseRepo.UpdateHouse(&house); err != nil {
		h.logger.Error("Failed to update house", err)
		http.Error(w, "Failed to update house", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(house); err != nil {
		h.logger.Error("Failed to encode response", err)
	}
}

func (h *HouseHandler) DeleteHouse(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	path := r.URL.Path
	idStr := path[len("/api/houses/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid house ID", http.StatusBadRequest)
		return
	}

	if err := h.houseRepo.DeleteHouse(id); err != nil {
		h.logger.Error("Failed to delete house", err)
		http.Error(w, "Failed to delete house", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *HouseHandler) GetAgents(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	agents, err := h.agentRepo.GetAllAgents()
	if err != nil {
		h.logger.Error("Failed to get agents", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(agents); err != nil {
		h.logger.Error("Failed to encode response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *HouseHandler) GetHouseTypes(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	houseTypes, err := h.houseTypeRepo.GetAllHouseTypes()
	if err != nil {
		h.logger.Error("Failed to get house types", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(houseTypes); err != nil {
		h.logger.Error("Failed to encode response", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *HouseHandler) HandleHousesRoute(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case path == "/api/houses" && r.Method == http.MethodGet:
		h.GetAllHouses(w, r)
	case path == "/api/houses" && r.Method == http.MethodPost:
		h.CreateHouse(w, r)
	case path == "/api/houses/top":
		h.GetTopHouses(w, r)
	case len(path) > len("/api/houses/") && path[:len("/api/houses/")] == "/api/houses/":
		switch r.Method {
		case http.MethodGet:
			h.GetHouseByID(w, r)
		case http.MethodPut:
			h.UpdateHouse(w, r)
		case http.MethodDelete:
			h.DeleteHouse(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}
