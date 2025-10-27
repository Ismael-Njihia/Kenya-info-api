# Kenya Info API - Example Requests

This file contains example API requests you can use with curl, Postman, or any HTTP client.

## Base URL
```
http://localhost:8080
```

## Health Check

### Check API Health
```bash
curl -X GET http://localhost:8080/health
```

## Counties

### Get All Counties (Paginated)
```bash
curl -X GET "http://localhost:8080/api/v1/counties?page=1&page_size=10"
```

### Get All Counties (All 47)
```bash
curl -X GET "http://localhost:8080/api/v1/counties?page_size=47"
```

### Get County by ID
```bash
# Replace {id} with actual county ID from the GET all counties response
curl -X GET http://localhost:8080/api/v1/counties/{id}
```

### Get County by Code (1-47)
```bash
# Get Nairobi County (code 47)
curl -X GET http://localhost:8080/api/v1/counties/code/47

# Get Mombasa County (code 1)
curl -X GET http://localhost:8080/api/v1/counties/code/1
```

### Search County by Name
```bash
# Search for Nairobi
curl -X GET "http://localhost:8080/api/v1/counties/search?name=Nairobi"

# Search for Mombasa
curl -X GET "http://localhost:8080/api/v1/counties/search?name=Mombasa"
```

### Create New County
```bash
curl -X POST http://localhost:8080/api/v1/counties \
  -H "Content-Type: application/json" \
  -d '{
    "code": 48,
    "name": "Example County",
    "capital": "Example Capital",
    "population": 500000,
    "area": 2500.5,
    "governor": "Example Governor"
  }'
```

### Update County
```bash
# Replace {id} with actual county ID
curl -X PUT http://localhost:8080/api/v1/counties/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "code": 48,
    "name": "Updated County",
    "capital": "Updated Capital",
    "population": 550000,
    "area": 2600.5,
    "governor": "Updated Governor"
  }'
```

### Delete County
```bash
# Replace {id} with actual county ID
curl -X DELETE http://localhost:8080/api/v1/counties/{id}
```

## Wards

### Get All Wards (Paginated)
```bash
curl -X GET "http://localhost:8080/api/v1/wards?page=1&page_size=20"
```

### Get Ward by ID
```bash
# Replace {id} with actual ward ID
curl -X GET http://localhost:8080/api/v1/wards/{id}
```

### Get Wards by County ID
```bash
# Replace {county_id} with actual county ID
curl -X GET "http://localhost:8080/api/v1/counties/{county_id}/wards?page=1&page_size=50"
```

### Create New Ward
```bash
# Replace {county_id} with actual county ID
curl -X POST http://localhost:8080/api/v1/wards \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Example Ward",
    "county_id": "{county_id}",
    "county_name": "Nairobi",
    "mca": "Example MCA",
    "population": 80000
  }'
```

### Update Ward
```bash
# Replace {id} with actual ward ID
curl -X PUT http://localhost:8080/api/v1/wards/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Ward",
    "county_id": "{county_id}",
    "county_name": "Nairobi",
    "mca": "Updated MCA",
    "population": 85000
  }'
```

### Delete Ward
```bash
# Replace {id} with actual ward ID
curl -X DELETE http://localhost:8080/api/v1/wards/{id}
```

## Leaders

### Get All Leaders (Paginated)
```bash
curl -X GET "http://localhost:8080/api/v1/leaders?page=1&page_size=20"
```

### Get Leader by ID
```bash
# Replace {id} with actual leader ID
curl -X GET http://localhost:8080/api/v1/leaders/{id}
```

### Get Leaders by Position
```bash
# Get all Governors
curl -X GET "http://localhost:8080/api/v1/leaders/position?position=Governor"

# Get all Senators
curl -X GET "http://localhost:8080/api/v1/leaders/position?position=Senator"

# Get all MCAs
curl -X GET "http://localhost:8080/api/v1/leaders/position?position=MCA"

# Get President
curl -X GET "http://localhost:8080/api/v1/leaders/position?position=President"
```

### Get Leaders by County ID
```bash
# Replace {county_id} with actual county ID
curl -X GET "http://localhost:8080/api/v1/counties/{county_id}/leaders?page=1&page_size=50"
```

### Create New Leader
```bash
curl -X POST http://localhost:8080/api/v1/leaders \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Example Leader",
    "position": "Senator",
    "county_id": "{county_id}",
    "county_name": "Nairobi",
    "party": "Example Party",
    "email": "leader@example.com",
    "phone": "+254712345678"
  }'
```

### Update Leader
```bash
# Replace {id} with actual leader ID
curl -X PUT http://localhost:8080/api/v1/leaders/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Leader",
    "position": "Senator",
    "county_id": "{county_id}",
    "county_name": "Nairobi",
    "party": "Updated Party",
    "email": "updated@example.com",
    "phone": "+254787654321"
  }'
```

### Delete Leader
```bash
# Replace {id} with actual leader ID
curl -X DELETE http://localhost:8080/api/v1/leaders/{id}
```

## Example Workflows

### 1. Get all counties and their details
```bash
# Get all counties
curl -X GET "http://localhost:8080/api/v1/counties?page_size=47" > counties.json

# Get specific county by code
curl -X GET http://localhost:8080/api/v1/counties/code/47
```

### 2. Find a county and its wards
```bash
# Step 1: Search for Nairobi county
curl -X GET "http://localhost:8080/api/v1/counties/search?name=Nairobi"

# Step 2: Copy the county ID from response and get its wards
# Example: curl -X GET "http://localhost:8080/api/v1/counties/507f1f77bcf86cd799439011/wards"
```

### 3. Find leaders in a specific county
```bash
# Step 1: Get county by code
curl -X GET http://localhost:8080/api/v1/counties/code/47

# Step 2: Use the county ID to get leaders
# Example: curl -X GET "http://localhost:8080/api/v1/counties/507f1f77bcf86cd799439011/leaders"
```

### 4. Filter leaders by position
```bash
# Get all Governors
curl -X GET "http://localhost:8080/api/v1/leaders/position?position=Governor&page_size=47"

# Get all national leaders
curl -X GET "http://localhost:8080/api/v1/leaders/position?position=President"
curl -X GET "http://localhost:8080/api/v1/leaders/position?position=Deputy"
```

## Response Format

### Success Response
```json
{
  "success": true,
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "code": 47,
    "name": "Nairobi",
    "capital": "Nairobi City",
    "population": 4397073,
    "area": 696.0,
    "governor": "Johnson Sakaja",
    "created_at": "2025-10-27T12:00:00Z",
    "updated_at": "2025-10-27T12:00:00Z"
  }
}
```

### Paginated Response
```json
{
  "success": true,
  "data": [...],
  "page": 1,
  "page_size": 10,
  "total_count": 47,
  "total_pages": 5
}
```

### Error Response
```json
{
  "success": false,
  "error": "County not found"
}
```

## Testing with Postman

1. Import this file into Postman
2. Create a new environment with variable `base_url` = `http://localhost:8080`
3. Replace all `http://localhost:8080` with `{{base_url}}`
4. Run the requests

## Testing with HTTPie

```bash
# Install HTTPie: pip install httpie

# Get all counties
http GET http://localhost:8080/api/v1/counties

# Search county
http GET http://localhost:8080/api/v1/counties/search name==Nairobi

# Create county
http POST http://localhost:8080/api/v1/counties \
  code:=48 \
  name="Test County" \
  capital="Test Capital" \
  population:=100000 \
  area:=1000.0
```

## Rate Limiting

Currently, there is no rate limiting implemented. In production, you may want to add rate limiting middleware.

## Authentication

Currently, the API is open. In production, you may want to add authentication/authorization.

## CORS

CORS is enabled for all origins. In production, configure specific allowed origins in the middleware.
