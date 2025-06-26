package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"thugcorp.io/nomado/logger"
	"thugcorp.io/nomado/models"
	"thugcorp.io/nomado/repository"
)

// Response structures for API responses
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
	Error      string      `json:"error,omitempty"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

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

// Helper methods for consistent API responses
func (h *HouseHandler) sendSuccessResponse(w http.ResponseWriter, data interface{}, message string) {
	response := APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
	h.sendJSONResponse(w, http.StatusOK, response)
}

func (h *HouseHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, errorMsg string) {
	response := APIResponse{
		Success: false,
		Error:   errorMsg,
	}
	h.sendJSONResponse(w, statusCode, response)
}

func (h *HouseHandler) sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", err)
	}
}

func (h *HouseHandler) GetTopHouses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
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
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve top houses")
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

	h.sendSuccessResponse(w, housesWithDetails, "Top houses retrieved successfully")
}

func (h *HouseHandler) GetAllHouses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	houses, err := h.houseRepo.GetAllHouses()
	if err != nil {
		h.logger.Error("Failed to get all houses", err)
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve houses")
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

	h.sendSuccessResponse(w, housesWithDetails, "Houses retrieved successfully")
}

func (h *HouseHandler) GetHouseByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract ID from URL path
	path := r.URL.Path
	idStr := path[len("/api/houses/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid house ID")
		return
	}

	house, err := h.houseRepo.GetHouseByID(id)
	if err != nil {
		h.logger.Error("Failed to get house by ID", err)
		h.sendErrorResponse(w, http.StatusNotFound, "House not found")
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

	h.sendSuccessResponse(w, houseWithDetails, "House retrieved successfully")
}

func (h *HouseHandler) CreateHouse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var house models.House
	if err := json.NewDecoder(r.Body).Decode(&house); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// Basic validation
	if house.Name == "" {
		h.sendErrorResponse(w, http.StatusBadRequest, "House name is required")
		return
	}
	if house.Price <= 0 {
		h.sendErrorResponse(w, http.StatusBadRequest, "House price must be greater than 0")
		return
	}

	if err := h.houseRepo.CreateHouse(&house); err != nil {
		h.logger.Error("Failed to create house", err)
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to create house")
		return
	}

	h.sendJSONResponse(w, http.StatusCreated, APIResponse{
		Success: true,
		Data:    house,
		Message: "House created successfully",
	})
}

func (h *HouseHandler) UpdateHouse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract ID from URL path
	path := r.URL.Path
	idStr := path[len("/api/houses/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid house ID")
		return
	}

	var house models.House
	if err := json.NewDecoder(r.Body).Decode(&house); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// Basic validation
	if house.Name == "" {
		h.sendErrorResponse(w, http.StatusBadRequest, "House name is required")
		return
	}
	if house.Price <= 0 {
		h.sendErrorResponse(w, http.StatusBadRequest, "House price must be greater than 0")
		return
	}

	house.ID = id
	if err := h.houseRepo.UpdateHouse(&house); err != nil {
		h.logger.Error("Failed to update house", err)
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to update house")
		return
	}

	h.sendSuccessResponse(w, house, "House updated successfully")
}

func (h *HouseHandler) DeleteHouse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract ID from URL path
	path := r.URL.Path
	idStr := path[len("/api/houses/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid house ID")
		return
	}

	if err := h.houseRepo.DeleteHouse(id); err != nil {
		h.logger.Error("Failed to delete house", err)
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to delete house")
		return
	}

	h.sendJSONResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "House deleted successfully",
	})
}

func (h *HouseHandler) GetAgents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	agents, err := h.agentRepo.GetAllAgents()
	if err != nil {
		h.logger.Error("Failed to get agents", err)
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve agents")
		return
	}

	h.sendSuccessResponse(w, agents, "Agents retrieved successfully")
}

func (h *HouseHandler) GetHouseTypes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	houseTypes, err := h.houseTypeRepo.GetAllHouseTypes()
	if err != nil {
		h.logger.Error("Failed to get house types", err)
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve house types")
		return
	}

	h.sendSuccessResponse(w, houseTypes, "House types retrieved successfully")
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
			h.sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	default:
		h.sendErrorResponse(w, http.StatusNotFound, "Endpoint not found")
	}
}
