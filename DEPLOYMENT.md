# Deployment Guide

This guide covers deploying the Kenya Info API to production.

## Table of Contents
- [MongoDB Atlas Setup](#mongodb-atlas-setup)
- [Local Development](#local-development)
- [Docker Deployment](#docker-deployment)
- [Production Deployment](#production-deployment)
- [Environment Variables](#environment-variables)

## MongoDB Atlas Setup

1. **Create a MongoDB Atlas Account**
   - Go to https://www.mongodb.com/cloud/atlas
   - Sign up for a free account

2. **Create a Cluster**
   - Click "Build a Database"
   - Choose the free tier (M0)
   - Select your preferred cloud provider and region
   - Click "Create Cluster"

3. **Configure Database Access**
   - Go to "Database Access" in the left sidebar
   - Click "Add New Database User"
   - Create a username and password
   - Grant "Read and write to any database" role

4. **Configure Network Access**
   - Go to "Network Access" in the left sidebar
   - Click "Add IP Address"
   - For development: Add your current IP
   - For production: Add your server's IP or 0.0.0.0/0 (allow from anywhere)

5. **Get Connection String**
   - Go to "Database" and click "Connect"
   - Choose "Connect your application"
   - Copy the connection string
   - Replace `<password>` with your database user password
   - Update your `.env` file with this connection string

## Local Development

1. **Install Dependencies**
```bash
make install
```

2. **Set Environment Variables**
```bash
cp .env.example .env
# Edit .env with your MongoDB Atlas connection string
```

3. **Generate Swagger Docs**
```bash
go install github.com/swaggo/swag/cmd/swag@latest
make swagger
```

4. **Run the Application**
```bash
make run
```

5. **Seed the Database**
```bash
make seed
```

6. **Access the API**
- API: http://localhost:8080
- Swagger Docs: http://localhost:8080/swagger/index.html
- Health Check: http://localhost:8080/health

## Docker Deployment

### Using Docker Compose (Recommended for Development)

1. **Update docker-compose.yml**
   - If using MongoDB Atlas, update the `MONGODB_URI` in docker-compose.yml
   - Or use the included MongoDB service for local testing

2. **Build and Run**
```bash
make docker-build
make docker-up
```

3. **Check Logs**
```bash
make docker-logs
```

4. **Stop Containers**
```bash
make docker-down
```

### Using Docker Only

1. **Build the Image**
```bash
docker build -t kenya-info-api .
```

2. **Run the Container**
```bash
docker run -d \
  --name kenya-api \
  -p 8080:8080 \
  -e MONGODB_URI="your-mongodb-atlas-uri" \
  -e MONGODB_DATABASE="kenya_info" \
  -e SERVER_PORT="8080" \
  -e LOG_LEVEL="info" \
  kenya-info-api
```

## Production Deployment

### Option 1: Deploy to Cloud Platform (Recommended)

#### Deploy to Heroku

1. **Install Heroku CLI**
```bash
curl https://cli-assets.heroku.com/install.sh | sh
```

2. **Login to Heroku**
```bash
heroku login
```

3. **Create a New App**
```bash
heroku create your-app-name
```

4. **Set Environment Variables**
```bash
heroku config:set MONGODB_URI="your-mongodb-atlas-uri"
heroku config:set MONGODB_DATABASE="kenya_info"
heroku config:set SERVER_PORT="8080"
heroku config:set LOG_LEVEL="info"
```

5. **Deploy**
```bash
git push heroku main
```

#### Deploy to DigitalOcean App Platform

1. **Connect Your GitHub Repository**
   - Go to DigitalOcean App Platform
   - Click "Create App"
   - Connect your GitHub repository

2. **Configure the App**
   - Select the branch to deploy
   - DigitalOcean will auto-detect the Dockerfile
   - Set environment variables in the App settings

3. **Deploy**
   - Click "Deploy"
   - Your app will be built and deployed automatically

#### Deploy to AWS ECS/Fargate

1. **Build and Push Docker Image**
```bash
aws ecr create-repository --repository-name kenya-info-api
docker build -t kenya-info-api .
docker tag kenya-info-api:latest <aws-account-id>.dkr.ecr.<region>.amazonaws.com/kenya-info-api:latest
docker push <aws-account-id>.dkr.ecr.<region>.amazonaws.com/kenya-info-api:latest
```

2. **Create ECS Task Definition**
   - Define your container with environment variables
   - Set MongoDB Atlas URI as an environment variable

3. **Create ECS Service**
   - Choose Fargate launch type
   - Configure load balancer
   - Deploy the service

### Option 2: Deploy to VPS (Ubuntu)

1. **SSH into Your Server**
```bash
ssh user@your-server-ip
```

2. **Install Docker**
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
```

3. **Clone Repository**
```bash
git clone https://github.com/Ismael-Njihia/Kenya-info-api.git
cd Kenya-info-api
```

4. **Create .env File**
```bash
cp .env.example .env
nano .env  # Edit with your MongoDB Atlas credentials
```

5. **Build and Run**
```bash
docker build -t kenya-info-api .
docker run -d \
  --name kenya-api \
  --restart unless-stopped \
  -p 80:8080 \
  --env-file .env \
  kenya-info-api
```

6. **Set Up Nginx (Optional)**
```bash
sudo apt install nginx
# Configure nginx as reverse proxy
```

7. **Set Up SSL with Let's Encrypt (Optional)**
```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d yourdomain.com
```

## Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `SERVER_PORT` | Port for the API server | No | `8080` |
| `SERVER_HOST` | Host address | No | `0.0.0.0` |
| `SERVER_MODE` | Gin mode (release/debug) | No | `release` |
| `MONGODB_URI` | MongoDB connection string | Yes | `mongodb://localhost:27017` |
| `MONGODB_DATABASE` | Database name | No | `kenya_info` |
| `MONGODB_TIMEOUT` | Connection timeout in seconds | No | `10` |
| `LOG_LEVEL` | Logging level (debug/info/warn/error) | No | `info` |
| `LOG_OUTPUT_PATH` | Log output destination | No | `stdout` |

## Monitoring and Maintenance

### Health Checks
```bash
curl http://your-api-url/health
```

### View Logs (Docker)
```bash
docker logs -f kenya-api
```

### Update Deployment
```bash
git pull origin main
docker build -t kenya-info-api .
docker stop kenya-api
docker rm kenya-api
docker run -d --name kenya-api --restart unless-stopped -p 80:8080 --env-file .env kenya-info-api
```

### Backup MongoDB Atlas
- MongoDB Atlas provides automatic backups
- Configure backup schedule in Atlas dashboard
- Set up point-in-time recovery if needed

## Performance Optimization

1. **Enable MongoDB Atlas Indexes**
   - Create indexes on frequently queried fields
   - Monitor slow queries in Atlas dashboard

2. **Add Caching** (Future Enhancement)
   - Consider adding Redis for caching
   - Cache frequently accessed data

3. **Load Balancing** (For High Traffic)
   - Use cloud load balancers
   - Run multiple instances of the API

4. **Database Optimization**
   - Monitor query performance
   - Optimize data models based on usage patterns

## Troubleshooting

### Cannot Connect to MongoDB
- Check MongoDB Atlas network access settings
- Verify connection string is correct
- Check if database user has proper permissions

### API Returns 502/503 Errors
- Check if container is running: `docker ps`
- Check container logs: `docker logs kenya-api`
- Verify MongoDB connection is active

### High Response Times
- Check MongoDB Atlas metrics
- Monitor API logs for slow queries
- Consider adding database indexes
- Review server resources (CPU, memory)

## Security Best Practices

1. **Use Environment Variables**
   - Never commit credentials to git
   - Use secrets management in production

2. **Enable HTTPS**
   - Use SSL/TLS certificates
   - Redirect HTTP to HTTPS

3. **MongoDB Security**
   - Use strong passwords
   - Enable MongoDB Atlas encryption
   - Restrict IP access

4. **API Security** (Future Enhancements)
   - Add rate limiting
   - Implement authentication
   - Add request validation

## Support

For issues and questions:
- GitHub Issues: https://github.com/Ismael-Njihia/Kenya-info-api/issues
- Documentation: See README.md
