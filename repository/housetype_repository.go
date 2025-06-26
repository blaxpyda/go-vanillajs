package repository

import (
	"database/sql"
	"fmt"

	"thugcorp.io/nomado/models"
)

type HouseTypeRepository struct {
	db *sql.DB
}

func NewHouseTypeRepository(db *sql.DB) *HouseTypeRepository {
	return &HouseTypeRepository{db: db}
}

func (htr *HouseTypeRepository) GetAllHouseTypes() ([]models.HouseType, error) {
	query := `
		SELECT id, name
		FROM house_types
		ORDER BY name
	`

	rows, err := htr.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query house types: %w", err)
	}
	defer rows.Close()

	var houseTypes []models.HouseType
	for rows.Next() {
		var houseType models.HouseType
		err := rows.Scan(&houseType.ID, &houseType.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan house type: %w", err)
		}
		houseTypes = append(houseTypes, houseType)
	}

	return houseTypes, nil
}

func (htr *HouseTypeRepository) GetHouseTypeByID(id int) (*models.HouseType, error) {
	query := `
		SELECT id, name
		FROM house_types
		WHERE id = $1
	`

	var houseType models.HouseType
	err := htr.db.QueryRow(query, id).Scan(&houseType.ID, &houseType.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("house type with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to query house type: %w", err)
	}

	return &houseType, nil
}

func (htr *HouseTypeRepository) CreateHouseType(houseType *models.HouseType) error {
	query := `
		INSERT INTO house_types (name)
		VALUES ($1)
		RETURNING id
	`

	err := htr.db.QueryRow(query, houseType.Name).Scan(&houseType.ID)

	if err != nil {
		return fmt.Errorf("failed to create house type: %w", err)
	}

	return nil
}

func (htr *HouseTypeRepository) UpdateHouseType(houseType *models.HouseType) error {
	query := `
		UPDATE house_types 
		SET name = $1
		WHERE id = $2
	`

	result, err := htr.db.Exec(query, houseType.Name, houseType.ID)
	if err != nil {
		return fmt.Errorf("failed to update house type: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("house type with id %d not found", houseType.ID)
	}

	return nil
}

func (htr *HouseTypeRepository) DeleteHouseType(id int) error {
	query := `DELETE FROM house_types WHERE id = $1`

	result, err := htr.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete house type: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("house type with id %d not found", id)
	}

	return nil
}
