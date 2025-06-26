#!/bin/bash

# PostgreSQL Setup Script for Nomado Real Estate App
# This script helps set up PostgreSQL for the Nomado application

echo "=== Nomado Real Estate App - Database Setup ==="

# Check if PostgreSQL is installed
if ! command -v psql &> /dev/null; then
    echo "PostgreSQL is not installed. Please install it first:"
    echo "Ubuntu/Debian: sudo apt-get install postgresql postgresql-contrib"
    echo "CentOS/RHEL: sudo yum install postgresql-server postgresql-contrib"
    echo "macOS: brew install postgresql"
    exit 1
fi

# Default database configuration
DB_NAME="nomado"
DB_USER="postgres"
DB_HOST="localhost"
DB_PORT="5432"

echo "Setting up database: $DB_NAME"
echo "Database user: $DB_USER"
echo "Host: $DB_HOST"
echo "Port: $DB_PORT"

# Create database if it doesn't exist
echo "Creating database if it doesn't exist..."
sudo -u postgres createdb $DB_NAME 2>/dev/null || echo "Database $DB_NAME may already exist"

# Test connection
echo "Testing database connection..."
if sudo -u postgres psql -d $DB_NAME -c "SELECT 1;" &> /dev/null; then
    echo "✓ Database connection successful"
else
    echo "✗ Database connection failed"
    echo "Please check your PostgreSQL installation and configuration"
    exit 1
fi

echo ""
echo "Database setup complete!"
echo ""
echo "To run the application:"
echo "1. Make sure PostgreSQL is running"
echo "2. Update the .env file with your database credentials if needed"
echo "3. Run: go run main.go"
echo ""
echo "The application will automatically create tables and seed sample data."
echo "Access the application at: http://localhost:8080"
