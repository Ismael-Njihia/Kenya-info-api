# Kenya Info API

A fast and reliable backend system providing structured data on Kenya's counties, wards, and political leaders. Built with Go, MongoDB Atlas, and designed for speed, scalability, and easy data updates.

## Features

- 🚀 **Fast & Reliable**: Built with Go and optimized for high performance
- 📊 **Structured Data**: Clean JSON responses for counties, wards, and leaders
- 🔍 **RESTful API**: Intuitive endpoints for easy integration
- 📝 **Swagger Documentation**: Interactive API documentation
- 🐳 **Docker Support**: Easy deployment with Docker and docker-compose
- 📋 **Logging**: Structured logging with Zap
- ✅ **Testing**: Comprehensive unit tests
- 🔄 **CI/CD**: GitHub Actions pipeline for automated testing and building
- 🌐 **CORS Enabled**: Ready for web applications

## Tech Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: MongoDB Atlas
- **Logging**: Zap
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker & Docker Compose
- **Testing**: Go testing package with testify
- **CI/CD**: GitHub Actions

## API Endpoints

### Health Check
- `GET /health` - Check API and database health

### Counties
- `GET /api/v1/counties` - Get all counties
- `GET /api/v1/counties/:id` - Get county by ID
- `GET /api/v1/counties/code/:code` - Get county by code
- `POST /api/v1/counties` - Create a new county
- `PUT /api/v1/counties/:id` - Update a county
- `DELETE /api/v1/counties/:id` - Delete a county

### Wards
- `GET /api/v1/wards` - Get all wards
- `GET /api/v1/wards/:id` - Get ward by ID
- `GET /api/v1/wards/county/:county_id` - Get wards by county
- `POST /api/v1/wards` - Create a new ward
- `PUT /api/v1/wards/:id` - Update a ward
- `DELETE /api/v1/wards/:id` - Delete a ward

### Leaders
- `GET /api/v1/leaders` - Get all leaders
- `GET /api/v1/leaders/:id` - Get leader by ID
- `GET /api/v1/leaders/position/:position` - Get leaders by position
- `GET /api/v1/leaders/county/:county` - Get leaders by county
- `POST /api/v1/leaders` - Create a new leader
- `PUT /api/v1/leaders/:id` - Update a leader
- `DELETE /api/v1/leaders/:id` - Delete a leader

### Documentation
- `GET /swagger/index.html` - Interactive Swagger documentation

## Quick Start

### Prerequisites
- Go 1.21 or higher
- MongoDB (local or MongoDB Atlas)
- Docker (optional)

### Local Development

1. **Clone the repository**
```bash
git clone https://github.com/Ismael-Njihia/Kenya-info-api.git
cd Kenya-info-api
```

2. **Install dependencies**
```bash
make install
```

3. **Set up environment variables**
```bash
cp .env.example .env
# Edit .env with your MongoDB connection string
```

4. **Generate Swagger documentation**
```bash
go install github.com/swaggo/swag/cmd/swag@latest
make swagger
```

5. **Run the application**
```bash
make run
```

The API will be available at `http://localhost:8080`

6. **Seed the database (optional)**
```bash
make seed
```

### Using Docker

1. **Build and start containers**
```bash
make docker-build
make docker-up
```

2. **View logs**
```bash
make docker-logs
```

3. **Stop containers**
```bash
make docker-down
```

## Configuration

Configuration is done via environment variables. See `.env.example` for all available options:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | API server port | `8080` |
| `SERVER_HOST` | API server host | `0.0.0.0` |
| `SERVER_MODE` | Gin mode (release/debug) | `release` |
| `MONGODB_URI` | MongoDB connection string | `mongodb://localhost:27017` |
| `MONGODB_DATABASE` | Database name | `kenya_info` |
| `MONGODB_TIMEOUT` | Connection timeout (seconds) | `10` |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | `info` |
| `LOG_OUTPUT_PATH` | Log output (stdout or file path) | `stdout` |

## Development

### Running Tests
```bash
make test
```

### Test Coverage
```bash
make test-coverage
```

### Linting
```bash
make lint
```

### Building
```bash
make build
```

## Project Structure

```
.
├── cmd/
│   └── api/              # Application entrypoint
│       └── main.go
├── config/               # Configuration management
│   └── config.go
├── internal/
│   ├── handlers/         # HTTP handlers
│   ├── middleware/       # HTTP middleware
│   ├── models/          # Data models
│   ├── repository/      # Database operations
│   └── logger/          # Logging utilities
├── scripts/             # Utility scripts
│   └── seed.go          # Database seeding
├── docs/                # Swagger documentation (generated)
├── .github/
│   └── workflows/       # GitHub Actions CI/CD
├── docker-compose.yml   # Docker Compose configuration
├── Dockerfile           # Docker image definition
├── Makefile            # Development commands
└── README.md           # This file
```

## MongoDB Atlas Setup

1. Create a free MongoDB Atlas account at https://www.mongodb.com/cloud/atlas
2. Create a new cluster
3. Add your IP address to the IP whitelist
4. Create a database user
5. Get your connection string and update `MONGODB_URI` in `.env`

Example connection string:
```
mongodb+srv://username:password@cluster.mongodb.net/kenya_info?retryWrites=true&w=majority
```

## Data Updates

To update data in the API:

1. **Via API**: Use POST, PUT, DELETE endpoints
2. **Via Seed Script**: Modify `scripts/seed.go` and run `make seed`
3. **Direct MongoDB**: Connect to your MongoDB instance and update collections directly

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support, email support@kenyainfo.api or open an issue in the GitHub repository.

## Roadmap

- [ ] Add authentication and authorization
- [ ] Add rate limiting
- [ ] Add caching layer (Redis)
- [ ] Add more comprehensive data for all 47 counties
- [ ] Add constituency data
- [ ] Add electoral data
- [ ] Add GraphQL support
- [ ] Add WebSocket support for real-time updates

## Acknowledgments

- Data sourced from official Kenyan government sources
- Built with ❤️ for civic tech and open data
