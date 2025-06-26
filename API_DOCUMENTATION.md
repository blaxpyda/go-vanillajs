# Nomado Real Estate API Documentation

## Overview

The Nomado Real Estate API is a RESTful backend service for managing real estate properties, agents, and property types. It provides comprehensive endpoints for CRUD operations and is designed to be consumed by web applications and mobile apps.

## Base URL

```
http://localhost:8080
```

## Response Format

All API responses follow a consistent format:

### Success Response
```json
{
  "success": true,
  "data": {...},
  "message": "Operation completed successfully"
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error description"
}
```

## Authentication

Currently, the API does not require authentication. This can be added later for production use.

## CORS

The API supports Cross-Origin Resource Sharing (CORS) with the following headers:
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type, Authorization`

## Endpoints

### Health Check

#### GET /api/health
Check if the API is running.

**Response:**
```json
{
  "status": "healthy",
  "service": "nomado-api"
}
```

### API Information

#### GET /api
Get information about available API endpoints.

**Response:**
```json
{
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
}
```

## Houses Endpoints

### GET /api/houses
Get all houses with their associated agent and house type information.

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Luxury Villa Downtown",
      "description": "Beautiful 4-bedroom villa...",
      "house_type_id": 1,
      "price": 850000.00,
      "tags": ["luxury", "downtown", "4-bedroom"],
      "image_url": "/images/logo.png",
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z",
      "agent_id": 1,
      "agent": {
        "id": 1,
        "first_name": "John",
        "last_name": "Smith",
        "image_url": "/images/generic_actor.jpg"
      },
      "house_type": {
        "id": 1,
        "name": "Villa"
      }
    }
  ],
  "message": "Houses retrieved successfully"
}
```

### GET /api/houses/top
Get top houses by price.

**Query Parameters:**
- `limit` (optional): Number of houses to return (default: 10)

**Example:** `/api/houses/top?limit=5`

**Response:** Same format as GET /api/houses

### GET /api/houses/{id}
Get a specific house by ID.

**Path Parameters:**
- `id`: House ID (integer)

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Luxury Villa Downtown",
    "description": "Beautiful 4-bedroom villa...",
    "house_type_id": 1,
    "price": 850000.00,
    "tags": ["luxury", "downtown", "4-bedroom"],
    "image_url": "/images/logo.png",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z",
    "agent_id": 1,
    "agent": {
      "id": 1,
      "first_name": "John",
      "last_name": "Smith",
      "image_url": "/images/generic_actor.jpg"
    },
    "house_type": {
      "id": 1,
      "name": "Villa"
    }
  },
  "message": "House retrieved successfully"
}
```

### POST /api/houses
Create a new house.

**Request Body:**
```json
{
  "name": "New Property",
  "description": "Property description",
  "house_type_id": 1,
  "price": 500000.00,
  "tags": ["modern", "spacious"],
  "image_url": "http://example.com/image.jpg",
  "agent_id": 1
}
```

**Validation Rules:**
- `name`: Required, non-empty string
- `price`: Required, must be greater than 0
- `house_type_id`: Optional, must reference existing house type
- `agent_id`: Optional, must reference existing agent

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 8,
    "name": "New Property",
    "description": "Property description",
    "house_type_id": 1,
    "price": 500000.00,
    "tags": ["modern", "spacious"],
    "image_url": "http://example.com/image.jpg",
    "created_at": "2025-06-26T10:30:00Z",
    "updated_at": "2025-06-26T10:30:00Z",
    "agent_id": 1
  },
  "message": "House created successfully"
}
```

### PUT /api/houses/{id}
Update an existing house.

**Path Parameters:**
- `id`: House ID (integer)

**Request Body:** Same as POST /api/houses

**Response:** Same format as POST response

### DELETE /api/houses/{id}
Delete a house.

**Path Parameters:**
- `id`: House ID (integer)

**Response:**
```json
{
  "success": true,
  "message": "House deleted successfully"
}
```

## Agents Endpoints

### GET /api/agents
Get all real estate agents.

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "first_name": "John",
      "last_name": "Smith",
      "image_url": "/images/generic_actor.jpg"
    },
    {
      "id": 2,
      "first_name": "Sarah",
      "last_name": "Johnson",
      "image_url": "/images/generic_actor.jpg"
    }
  ],
  "message": "Agents retrieved successfully"
}
```

## House Types Endpoints

### GET /api/house-types
Get all house types.

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Villa"
    },
    {
      "id": 2,
      "name": "Apartment"
    },
    {
      "id": 3,
      "name": "House"
    }
  ],
  "message": "House types retrieved successfully"
}
```

## Error Codes

The API uses standard HTTP status codes:

- `200 OK`: Successful GET request
- `201 Created`: Successful POST request
- `400 Bad Request`: Invalid request data
- `404 Not Found`: Resource not found
- `405 Method Not Allowed`: HTTP method not supported
- `500 Internal Server Error`: Server error

## Database Schema

### Houses Table
```sql
CREATE TABLE houses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    house_type_id INTEGER REFERENCES house_types(id),
    price DECIMAL(12,2) NOT NULL,
    tags TEXT, -- comma-separated tags
    image_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    agent_id INTEGER REFERENCES agents(id)
);
```

### Agents Table
```sql
CREATE TABLE agents (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    image_url TEXT
);
```

### House Types Table
```sql
CREATE TABLE house_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);
```

## Sample Data

The API comes with pre-seeded sample data:

### House Types
- Villa
- Apartment
- House
- Townhouse
- Condo

### Sample Agents
- John Smith
- Sarah Johnson
- Michael Davis
- Emily Wilson
- David Brown

### Sample Properties
- Luxury Villa Downtown ($850,000)
- Modern Apartment Complex ($320,000)
- Family House Suburbia ($450,000)
- Executive Townhouse ($620,000)
- City Center Condo ($280,000)
- Waterfront Villa ($1,200,000)
- Garden Apartment ($380,000)

## Usage Examples

### Get all houses
```bash
curl -X GET http://localhost:8080/api/houses
```

### Get top 3 houses
```bash
curl -X GET "http://localhost:8080/api/houses/top?limit=3"
```

### Get a specific house
```bash
curl -X GET http://localhost:8080/api/houses/1
```

### Create a new house
```bash
curl -X POST http://localhost:8080/api/houses \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Modern Condo",
    "description": "Stylish 2-bedroom condo",
    "house_type_id": 5,
    "price": 400000,
    "tags": ["modern", "condo", "2-bedroom"],
    "agent_id": 1
  }'
```

### Update a house
```bash
curl -X PUT http://localhost:8080/api/houses/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Villa Name",
    "description": "Updated description",
    "house_type_id": 1,
    "price": 900000,
    "tags": ["luxury", "updated"],
    "agent_id": 1
  }'
```

### Delete a house
```bash
curl -X DELETE http://localhost:8080/api/houses/1
```

## Development

### Prerequisites
- Go 1.19+
- PostgreSQL 12+

### Setup
1. Clone the repository
2. Configure database connection in `.env`
3. Run `go mod tidy`
4. Run `go run main.go`

The API will start on `http://localhost:8080`

## Future Enhancements

1. **Authentication & Authorization**: JWT-based authentication
2. **Pagination**: Add pagination support for large datasets
3. **Search & Filtering**: Advanced search capabilities
4. **Image Upload**: Support for property image uploads
5. **Rate Limiting**: API rate limiting for production
6. **Logging**: Enhanced structured logging
7. **Metrics**: API metrics and monitoring
8. **Documentation**: Interactive API documentation with Swagger
9. **Testing**: Comprehensive test suite
10. **Caching**: Redis caching for improved performance
