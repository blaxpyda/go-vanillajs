#!/bin/bash

# Nomado Real Estate API Test Script
# This script tests all the API endpoints

API_BASE="http://localhost:8080"

echo "=== Nomado Real Estate API Test Suite ==="
echo ""

# Function to test an endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo "Testing: $description"
    echo "Method: $method $API_BASE$endpoint"
    
    if [ -n "$data" ]; then
        echo "Data: $data"
        response=$(curl -s -X $method "$API_BASE$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data" \
            -w "\nHTTP_STATUS:%{http_code}")
    else
        response=$(curl -s -X $method "$API_BASE$endpoint" \
            -w "\nHTTP_STATUS:%{http_code}")
    fi
    
    http_status=$(echo "$response" | grep "HTTP_STATUS" | cut -d: -f2)
    response_body=$(echo "$response" | grep -v "HTTP_STATUS")
    
    echo "Status: $http_status"
    echo "Response: $response_body"
    echo "---"
    echo ""
}

# Check if server is running
echo "Checking if API server is running..."
if ! curl -s "$API_BASE/api/health" > /dev/null; then
    echo "❌ API server is not running on $API_BASE"
    echo "Please start the server with: go run main.go"
    exit 1
fi

echo "✅ API server is running!"
echo ""

# Test Health Check
test_endpoint "GET" "/api/health" "" "Health Check"

# Test API Info
test_endpoint "GET" "/api" "" "API Information"

# Test Get All Houses
test_endpoint "GET" "/api/houses" "" "Get All Houses"

# Test Get Top Houses
test_endpoint "GET" "/api/houses/top?limit=3" "" "Get Top 3 Houses"

# Test Get House by ID
test_endpoint "GET" "/api/houses/1" "" "Get House by ID (1)"

# Test Get All Agents
test_endpoint "GET" "/api/agents" "" "Get All Agents"

# Test Get House Types
test_endpoint "GET" "/api/house-types" "" "Get House Types"

# Test Create House
create_data='{
  "name": "Test Property",
  "description": "A test property created via API",
  "house_type_id": 1,
  "price": 500000.00,
  "tags": ["test", "api", "modern"],
  "agent_id": 1
}'
test_endpoint "POST" "/api/houses" "$create_data" "Create New House"

# Test Update House (assuming ID 8 was created)
update_data='{
  "name": "Updated Test Property",
  "description": "An updated test property",
  "house_type_id": 2,
  "price": 550000.00,
  "tags": ["updated", "test"],
  "agent_id": 2
}'
test_endpoint "PUT" "/api/houses/8" "$update_data" "Update House (ID 8)"

# Test Delete House
test_endpoint "DELETE" "/api/houses/8" "" "Delete House (ID 8)"

# Test Invalid Endpoints
test_endpoint "GET" "/api/invalid" "" "Invalid Endpoint (Should return 404)"
test_endpoint "POST" "/api/houses/1" "" "Invalid Method (Should return 405)"

echo "=== Test Suite Complete ==="
echo ""
echo "API Endpoints Summary:"
echo "- GET    /api/health        - Health check"
echo "- GET    /api              - API information"
echo "- GET    /api/houses       - All houses"
echo "- GET    /api/houses/top   - Top houses"
echo "- GET    /api/houses/{id}  - Specific house"
echo "- POST   /api/houses       - Create house"
echo "- PUT    /api/houses/{id}  - Update house"
echo "- DELETE /api/houses/{id}  - Delete house"
echo "- GET    /api/agents       - All agents"
echo "- GET    /api/house-types  - All house types"
