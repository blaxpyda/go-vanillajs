package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase() (*Database, error) {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file")
	}

	// Get database connection parameters from environment variables
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "")
	dbname := getEnv("DB_NAME", "nomado")
	sslmode := getEnv("DB_SSLMODE", "disable")

	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Successfully connected to database: %s:%s/%s", host, port, dbname)

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}

func (d *Database) CreateTables() error {
	schema := `
	-- Create house_types table
	CREATE TABLE IF NOT EXISTS house_types (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL UNIQUE
	);

	-- Create agents table
	CREATE TABLE IF NOT EXISTS agents (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		image_url TEXT
	);

	-- Create houses table
	CREATE TABLE IF NOT EXISTS houses (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		house_type_id INTEGER REFERENCES house_types(id) ON DELETE SET NULL,
		price DECIMAL(12,2) NOT NULL,
		tags TEXT, -- comma-separated tags
		image_url TEXT,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW(),
		agent_id INTEGER REFERENCES agents(id) ON DELETE SET NULL
	);

	-- Create indexes for better performance
	CREATE INDEX IF NOT EXISTS idx_houses_price ON houses(price);
	CREATE INDEX IF NOT EXISTS idx_houses_house_type_id ON houses(house_type_id);
	CREATE INDEX IF NOT EXISTS idx_houses_agent_id ON houses(agent_id);
	CREATE INDEX IF NOT EXISTS idx_houses_created_at ON houses(created_at);
	`

	_, err := d.DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("Database tables created successfully")
	return nil
}

func (d *Database) SeedData() error {
	// Insert sample house types
	houseTypesQuery := `
	INSERT INTO house_types (name) VALUES 
		('Villa'), 
		('Apartment'), 
		('House'), 
		('Townhouse'),
		('Condo')
	ON CONFLICT (name) DO NOTHING;
	`

	_, err := d.DB.Exec(houseTypesQuery)
	if err != nil {
		return fmt.Errorf("failed to seed house types: %w", err)
	}

	// Insert sample agents
	agentsQuery := `
	INSERT INTO agents (first_name, last_name, image_url) VALUES 
		('John', 'Smith', '/images/generic_actor.jpg'),
		('Sarah', 'Johnson', '/images/generic_actor.jpg'),
		('Michael', 'Davis', '/images/generic_actor.jpg'),
		('Emily', 'Wilson', '/images/generic_actor.jpg'),
		('David', 'Brown', '/images/generic_actor.jpg')
	ON CONFLICT DO NOTHING;
	`

	_, err = d.DB.Exec(agentsQuery)
	if err != nil {
		return fmt.Errorf("failed to seed agents: %w", err)
	}

	// Insert sample houses
	housesQuery := `
	INSERT INTO houses (name, description, house_type_id, price, tags, image_url, agent_id) VALUES 
		('Luxury Villa Downtown', 'Beautiful 4-bedroom villa in the heart of the city with stunning views and modern amenities.', 1, 850000.00, 'luxury,downtown,4-bedroom,modern', '/images/logo.png', 1),
		('Modern Apartment Complex', 'Contemporary 2-bedroom apartment with all modern conveniences and great location.', 2, 320000.00, 'modern,2-bedroom,apartment,convenient', '/images/logo.png', 2),
		('Family House Suburbia', 'Spacious 3-bedroom house perfect for families, with a large garden and quiet neighborhood.', 3, 450000.00, 'family,3-bedroom,garden,quiet', '/images/logo.png', 3),
		('Executive Townhouse', 'Elegant 3-bedroom townhouse with premium finishes and close to business district.', 4, 620000.00, 'executive,3-bedroom,premium,business', '/images/logo.png', 4),
		('City Center Condo', 'Stylish 1-bedroom condo in the city center with great amenities and city views.', 5, 280000.00, 'stylish,1-bedroom,city-center,views', '/images/logo.png', 5),
		('Waterfront Villa', 'Stunning waterfront villa with private beach access and panoramic ocean views.', 1, 1200000.00, 'waterfront,luxury,beach,ocean-views', '/images/logo.png', 1),
		('Garden Apartment', 'Charming 2-bedroom apartment with private garden and peaceful surroundings.', 2, 380000.00, 'charming,2-bedroom,garden,peaceful', '/images/logo.png', 2)
	ON CONFLICT DO NOTHING;
	`

	_, err = d.DB.Exec(housesQuery)
	if err != nil {
		return fmt.Errorf("failed to seed houses: %w", err)
	}

	log.Println("Database seeded successfully")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
