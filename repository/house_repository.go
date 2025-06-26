package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"thugcorp.io/nomado/models"
)

type HouseRepository struct {
	db *sql.DB
}

func NewHouseRepository(db *sql.DB) *HouseRepository {
	return &HouseRepository{db: db}
}

func (hr *HouseRepository) GetAllHouses() ([]models.House, error) {
	query := `
		SELECT h.id, h.name, h.description, h.house_type_id, h.price, 
			   h.tags, h.image_url, h.created_at, h.updated_at, h.agent_id
		FROM houses h
		ORDER BY h.created_at DESC
	`

	rows, err := hr.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query houses: %w", err)
	}
	defer rows.Close()

	var houses []models.House
	for rows.Next() {
		var house models.House
		var tagsStr string

		err := rows.Scan(
			&house.ID, &house.Name, &house.Description, &house.HouseTypeID,
			&house.Price, &tagsStr, &house.ImageURL, &house.CreatedAt,
			&house.UpdatedAt, &house.AgentID,
		)
		if err != nil {
			log.Printf("Error scanning house: %v", err)
			continue
		}

		// Parse tags from comma-separated string
		if tagsStr != "" {
			house.Tags = strings.Split(tagsStr, ",")
		}

		houses = append(houses, house)
	}

	return houses, nil
}

func (hr *HouseRepository) GetTopHouses(limit int) ([]models.House, error) {
	query := `
		SELECT h.id, h.name, h.description, h.house_type_id, h.price, 
			   h.tags, h.image_url, h.created_at, h.updated_at, h.agent_id
		FROM houses h
		ORDER BY h.price DESC
		LIMIT $1
	`

	rows, err := hr.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query top houses: %w", err)
	}
	defer rows.Close()

	var houses []models.House
	for rows.Next() {
		var house models.House
		var tagsStr string

		err := rows.Scan(
			&house.ID, &house.Name, &house.Description, &house.HouseTypeID,
			&house.Price, &tagsStr, &house.ImageURL, &house.CreatedAt,
			&house.UpdatedAt, &house.AgentID,
		)
		if err != nil {
			log.Printf("Error scanning house: %v", err)
			continue
		}

		// Parse tags from comma-separated string
		if tagsStr != "" {
			house.Tags = strings.Split(tagsStr, ",")
		}

		houses = append(houses, house)
	}

	return houses, nil
}

func (hr *HouseRepository) GetHouseByID(id int) (*models.House, error) {
	query := `
		SELECT h.id, h.name, h.description, h.house_type_id, h.price, 
			   h.tags, h.image_url, h.created_at, h.updated_at, h.agent_id
		FROM houses h
		WHERE h.id = $1
	`

	var house models.House
	var tagsStr string

	err := hr.db.QueryRow(query, id).Scan(
		&house.ID, &house.Name, &house.Description, &house.HouseTypeID,
		&house.Price, &tagsStr, &house.ImageURL, &house.CreatedAt,
		&house.UpdatedAt, &house.AgentID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("house with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to query house: %w", err)
	}

	// Parse tags from comma-separated string
	if tagsStr != "" {
		house.Tags = strings.Split(tagsStr, ",")
	}

	return &house, nil
}

func (hr *HouseRepository) CreateHouse(house *models.House) error {
	query := `
		INSERT INTO houses (name, description, house_type_id, price, tags, image_url, agent_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	tagsStr := strings.Join(house.Tags, ",")

	err := hr.db.QueryRow(
		query, house.Name, house.Description, house.HouseTypeID,
		house.Price, tagsStr, house.ImageURL, house.AgentID,
	).Scan(&house.ID, &house.CreatedAt, &house.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create house: %w", err)
	}

	return nil
}

func (hr *HouseRepository) UpdateHouse(house *models.House) error {
	query := `
		UPDATE houses 
		SET name = $1, description = $2, house_type_id = $3, price = $4, 
			tags = $5, image_url = $6, agent_id = $7, updated_at = NOW()
		WHERE id = $8
		RETURNING updated_at
	`

	tagsStr := strings.Join(house.Tags, ",")

	err := hr.db.QueryRow(
		query, house.Name, house.Description, house.HouseTypeID,
		house.Price, tagsStr, house.ImageURL, house.AgentID, house.ID,
	).Scan(&house.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update house: %w", err)
	}

	return nil
}

func (hr *HouseRepository) DeleteHouse(id int) error {
	query := `DELETE FROM houses WHERE id = $1`

	result, err := hr.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete house: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("house with id %d not found", id)
	}

	return nil
}
