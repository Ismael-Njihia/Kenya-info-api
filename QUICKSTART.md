# Quick Start Guide

This guide will help you get the Kenya Info API up and running quickly.

## Prerequisites

- Go 1.21+ installed
- MongoDB (local) or MongoDB Atlas account
- Docker (optional)

## Quick Setup (5 minutes)

### 1. Clone and Setup

```bash
# Clone the repository
git clone https://github.com/Ismael-Njihia/Kenya-info-api.git
cd Kenya-info-api

# Copy environment template
cp .env.example .env
```

### 2. Configure MongoDB

Edit `.env` file and set your MongoDB connection:

```env
MONGODB_URI=mongodb://localhost:27017
# or for MongoDB Atlas:
# MONGODB_URI=mongodb+srv://username:password@cluster.mongodb.net/
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Seed the Database

```bash
make seed
# or
go run scripts/seed.go
```

This will populate your database with:
- All 47 Kenyan counties
- Sample wards
- Sample leaders

### 5. Run the API

```bash
make run
# or
go run cmd/api/main.go
```

The API will start at `http://localhost:8080`

## Using Docker (Alternative)

```bash
# Start with Docker Compose (includes MongoDB)
docker-compose up -d

# Check logs
docker-compose logs -f

# Stop
docker-compose down
```

## Test the API

### Health Check

```bash
curl http://localhost:8080/health
```

### Get All Counties

```bash
curl http://localhost:8080/api/v1/counties
```

### Get Specific County by Code

```bash
curl http://localhost:8080/api/v1/counties/code/47
```

### Search County by Name

```bash
curl http://localhost:8080/api/v1/counties/search?name=Nairobi
```

### Get Wards in a County

First, get the county ID from the counties endpoint, then:

```bash
curl http://localhost:8080/api/v1/counties/{county_id}/wards
```

## Access Swagger Documentation

Open your browser and navigate to:

```
http://localhost:8080/swagger/index.html
```

## Common Commands

```bash
# Run tests
make test

# Build the application
make build

# Run linter (requires golangci-lint)
make lint

# View all available commands
make help
```

## Troubleshooting

### MongoDB Connection Issues

1. **Local MongoDB not running:**
   ```bash
   # Start MongoDB
   sudo systemctl start mongod  # Linux
   brew services start mongodb-community  # macOS
   ```

2. **Connection string incorrect:**
   - Verify your MONGODB_URI in `.env`
   - For Atlas, ensure your IP is whitelisted
   - Check username and password

### Port Already in Use

If port 8080 is already in use, change it in `.env`:

```env
PORT=3000
```

### Build Errors

```bash
# Clean and rebuild
make clean
go mod tidy
make build
```

## Next Steps

1. **Explore the API** using Swagger UI
2. **Add more data** by modifying `scripts/seed.go`
3. **Customize** the configuration in `.env`
4. **Deploy** to Render or Fly.io (see README.md)

## Sample API Calls

### Create a New County

```bash
curl -X POST http://localhost:8080/api/v1/counties \
  -H "Content-Type: application/json" \
  -d '{
    "code": 48,
    "name": "Test County",
    "capital": "Test Capital",
    "population": 100000,
    "area": 1000.0,
    "governor": "Test Governor"
  }'
```

### Get Leaders by Position

```bash
curl "http://localhost:8080/api/v1/leaders/position?position=Governor"
```

### Pagination Example

```bash
curl "http://localhost:8080/api/v1/counties?page=1&page_size=10"
```

## Need Help?

- Check the [README.md](README.md) for detailed documentation
- See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines
- Open an issue on GitHub for bugs or questions

Happy coding! 🇰🇪
