package repository

import (
	"database/sql"
	"fmt"

	"thugcorp.io/nomado/models"
)

type AgentRepository struct {
	db *sql.DB
}

func NewAgentRepository(db *sql.DB) *AgentRepository {
	return &AgentRepository{db: db}
}

func (ar *AgentRepository) GetAllAgents() ([]models.Agent, error) {
	query := `
		SELECT id, first_name, last_name, image_url
		FROM agents
		ORDER BY first_name, last_name
	`

	rows, err := ar.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query agents: %w", err)
	}
	defer rows.Close()

	var agents []models.Agent
	for rows.Next() {
		var agent models.Agent
		err := rows.Scan(&agent.ID, &agent.FirstName, &agent.LastName, &agent.ImageURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan agent: %w", err)
		}
		agents = append(agents, agent)
	}

	return agents, nil
}

func (ar *AgentRepository) GetAgentByID(id int) (*models.Agent, error) {
	query := `
		SELECT id, first_name, last_name, image_url
		FROM agents
		WHERE id = $1
	`

	var agent models.Agent
	err := ar.db.QueryRow(query, id).Scan(
		&agent.ID, &agent.FirstName, &agent.LastName, &agent.ImageURL,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("agent with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to query agent: %w", err)
	}

	return &agent, nil
}

func (ar *AgentRepository) CreateAgent(agent *models.Agent) error {
	query := `
		INSERT INTO agents (first_name, last_name, image_url)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := ar.db.QueryRow(
		query, agent.FirstName, agent.LastName, agent.ImageURL,
	).Scan(&agent.ID)

	if err != nil {
		return fmt.Errorf("failed to create agent: %w", err)
	}

	return nil
}

func (ar *AgentRepository) UpdateAgent(agent *models.Agent) error {
	query := `
		UPDATE agents 
		SET first_name = $1, last_name = $2, image_url = $3
		WHERE id = $4
	`

	result, err := ar.db.Exec(
		query, agent.FirstName, agent.LastName, agent.ImageURL, agent.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update agent: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("agent with id %d not found", agent.ID)
	}

	return nil
}

func (ar *AgentRepository) DeleteAgent(id int) error {
	query := `DELETE FROM agents WHERE id = $1`

	result, err := ar.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete agent: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("agent with id %d not found", id)
	}

	return nil
}
