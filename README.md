# Kenya Info API 🇰🇪

[![CI/CD Pipeline](https://github.com/Ismael-Njihia/Kenya-info-api/actions/workflows/ci-cd.yml/badge.svg)](https://github.com/Ismael-Njihia/Kenya-info-api/actions/workflows/ci-cd.yml)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A fast, reliable, and production-ready backend API providing structured data on Kenya's 47 counties, wards, and leaders. Built with Go and MongoDB Atlas, it delivers clean JSON responses optimized for civic apps, dashboards, and research tools.

## 🚀 Features

- **Fast & Reliable**: Built with Go and Gin framework for high performance
- **RESTful API**: Clean, intuitive endpoints following REST principles
- **MongoDB Atlas**: Scalable NoSQL database for flexible data storage
- **Swagger Documentation**: Interactive API documentation at `/swagger/index.html`
- **Structured Logging**: Production-grade logging with Zap
- **Docker Support**: Fully containerized for easy deployment
- **CI/CD Pipeline**: Automated testing and deployment with GitHub Actions
- **Pagination**: Efficient data retrieval with built-in pagination
- **CORS Enabled**: Cross-origin resource sharing support
- **Health Checks**: Monitor API status and uptime

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Deployment](#deployment)
- [Project Structure](#project-structure)

## 🛠️ Prerequisites

- Go 1.21 or higher
- MongoDB 7.0 or MongoDB Atlas account
- Docker (optional, for containerized deployment)
- Git

## 📦 Installation

1. **Clone the repository**

```bash
git clone https://github.com/Ismael-Njihia/Kenya-info-api.git
cd Kenya-info-api
```

2. **Install dependencies**

```bash
go mod download
```

3. **Set up environment variables**

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
PORT=8080
ENV=development
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=kenya_info
MONGODB_TIMEOUT=10
LOG_LEVEL=info
```

## ⚙️ Configuration

The application uses environment variables for configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `ENV` | Environment (development/production) | `development` |
| `MONGODB_URI` | MongoDB connection string | `mongodb://localhost:27017` |
| `MONGODB_DATABASE` | Database name | `kenya_info` |
| `MONGODB_TIMEOUT` | Connection timeout (seconds) | `10` |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | `info` |

## 🏃 Running the Application

### Local Development

```bash
go run cmd/api/main.go
```

### Using Docker

```bash
docker-compose up -d
```

### Build and Run

```bash
go build -o main cmd/api/main.go
./main
```

The API will be available at `http://localhost:8080`

### Seed the Database

To populate the database with sample data for all 47 Kenyan counties, wards, and leaders:

```bash
# Make sure MONGODB_URI is set in your .env file
make seed

# Or run directly
go run scripts/seed.go
```

This will:
- Clear existing data (optional, can be commented out in the script)
- Insert all 47 Kenyan counties with accurate data
- Add sample wards for multiple counties
- Add sample leaders (national, senators, MPs, MCAs)

## 📚 API Documentation

### Interactive Documentation

Once the server is running, access the Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

### Key Endpoints

#### Counties

- `GET /api/v1/counties` - Get all counties (paginated)
- `GET /api/v1/counties/:id` - Get county by ID
- `GET /api/v1/counties/code/:code` - Get county by code (1-47)
- `GET /api/v1/counties/search?name=Nairobi` - Search county by name
- `GET /api/v1/counties/:id/wards` - Get wards in a county
- `GET /api/v1/counties/:id/leaders` - Get leaders in a county
- `POST /api/v1/counties` - Create new county
- `PUT /api/v1/counties/:id` - Update county
- `DELETE /api/v1/counties/:id` - Delete county

#### Wards

- `GET /api/v1/wards` - Get all wards (paginated)
- `GET /api/v1/wards/:id` - Get ward by ID
- `POST /api/v1/wards` - Create new ward
- `PUT /api/v1/wards/:id` - Update ward
- `DELETE /api/v1/wards/:id` - Delete ward

#### Leaders

- `GET /api/v1/leaders` - Get all leaders (paginated)
- `GET /api/v1/leaders/:id` - Get leader by ID
- `GET /api/v1/leaders/position?position=Governor` - Filter by position
- `POST /api/v1/leaders` - Create new leader
- `PUT /api/v1/leaders/:id` - Update leader
- `DELETE /api/v1/leaders/:id` - Delete leader

#### Health Check

- `GET /health` - API health status

### Example Requests

**Get all counties:**
```bash
curl http://localhost:8080/api/v1/counties
```

**Get specific county:**
```bash
curl http://localhost:8080/api/v1/counties/code/1
```

**Search county by name:**
```bash
curl http://localhost:8080/api/v1/counties/search?name=Nairobi
```

**Get wards in a county (paginated):**
```bash
curl http://localhost:8080/api/v1/counties/{county_id}/wards?page=1&page_size=50
```

### Response Format

**Successful Response:**
```json
{
  "success": true,
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "code": 1,
    "name": "Nairobi",
    "capital": "Nairobi City",
    "population": 4397073,
    "area": 696.0,
    "governor": "Johnson Sakaja"
  }
}
```

**Paginated Response:**
```json
{
  "success": true,
  "data": [...],
  "page": 1,
  "page_size": 47,
  "total_count": 47,
  "total_pages": 1
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "County not found"
}
```

## 🧪 Testing

Run all tests:

```bash
go test ./... -v
```

Run tests with coverage:

```bash
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Run tests with race detection:

```bash
go test ./... -race
```

## 🚢 Deployment

### Deploy to Render

1. Create a new Web Service on [Render](https://render.com)
2. Connect your GitHub repository
3. Set environment variables in Render dashboard
4. Deploy!

### Deploy to Fly.io

1. Install Fly CLI: `curl -L https://fly.io/install.sh | sh`
2. Login: `fly auth login`
3. Launch app: `fly launch`
4. Set secrets: `fly secrets set MONGODB_URI=your_connection_string`
5. Deploy: `fly deploy`

### Deploy with Docker

```bash
docker build -t kenya-info-api .
docker run -p 8080:8080 -e MONGODB_URI=your_connection_string kenya-info-api
```

## 📁 Project Structure

```
Kenya-info-api/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   │   ├── config.go
│   │   └── config_test.go
│   ├── database/                # Database connection
│   │   └── database.go
│   ├── handlers/                # HTTP request handlers
│   │   ├── county_handler.go
│   │   ├── ward_handler.go
│   │   ├── leader_handler.go
│   │   └── health_handler.go
│   ├── middleware/              # HTTP middleware
│   │   └── middleware.go
│   ├── models/                  # Data models
│   │   └── models.go
│   └── services/                # Business logic
│       ├── county_service.go
│       ├── ward_service.go
│       ├── leader_service.go
│       └── services_test.go
├── pkg/
│   └── logger/                  # Logging utilities
│       └── logger.go
├── docs/                        # Swagger documentation
│   └── docs.go
├── .github/
│   └── workflows/
│       └── ci-cd.yml           # CI/CD pipeline
├── Dockerfile                  # Docker configuration
├── docker-compose.yml          # Docker Compose setup
├── .env.example               # Environment variables template
├── .gitignore                 # Git ignore file
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
└── README.md                  # This file
```

## 🏗️ Architecture

The application follows a clean architecture pattern:

- **cmd/api**: Application entry point
- **internal**: Private application code
  - **config**: Configuration loading and validation
  - **database**: Database connection management
  - **handlers**: HTTP request handlers (controllers)
  - **middleware**: HTTP middleware (logging, CORS, recovery)
  - **models**: Data structures and DTOs
  - **services**: Business logic layer
- **pkg**: Public, reusable packages (logger)
- **docs**: Auto-generated Swagger documentation

## 🔧 Technologies Used

- **Go 1.21+** - Programming language
- **Gin** - Web framework
- **MongoDB** - Database
- **Zap** - Structured logging
- **Swagger** - API documentation
- **Docker** - Containerization
- **GitHub Actions** - CI/CD
- **Testify** - Testing framework

## 📝 License

This project is licensed under the MIT License.

## 👥 Authors

- **Ismael Njihia** - Initial work

## 🙏 Acknowledgments

- Kenya's Independent Electoral and Boundaries Commission (IEBC)
- Open Data Kenya
- All contributors and supporters

---

**Built with ❤️ for Kenya**
