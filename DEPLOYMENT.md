# Deployment Checklist

Use this checklist when deploying the Kenya Info API to production.

## Pre-Deployment

### Code Quality
- [ ] All tests passing (`make test`)
- [ ] Code linted successfully (`make lint`)
- [ ] Build successful (`make build`)
- [ ] No security vulnerabilities
- [ ] Environment variables documented

### Configuration
- [ ] Production MongoDB URI configured
- [ ] Environment set to `production`
- [ ] Log level set appropriately (`info` or `warn`)
- [ ] Port configuration verified
- [ ] CORS settings configured for production domains

### Database
- [ ] MongoDB Atlas cluster created
- [ ] Database user created with appropriate permissions
- [ ] IP whitelist configured (or 0.0.0.0/0 for cloud deployments)
- [ ] Database indexes created
- [ ] Sample data seeded (optional)

## Deployment Options

### Option 1: Deploy to Render

1. **Setup**
   - [ ] Create Render account
   - [ ] Connect GitHub repository
   - [ ] Create new Web Service

2. **Configuration**
   - [ ] Set service name
   - [ ] Select Docker as environment
   - [ ] Set region (closest to users)
   - [ ] Choose instance type

3. **Environment Variables**
   ```
   MONGODB_URI=mongodb+srv://...
   MONGODB_DATABASE=kenya_info
   ENV=production
   LOG_LEVEL=info
   PORT=8080
   ```

4. **Deploy**
   - [ ] Click "Create Web Service"
   - [ ] Wait for build and deployment
   - [ ] Verify health endpoint: `https://your-app.onrender.com/health`
   - [ ] Test API endpoints
   - [ ] Check logs for errors

### Option 2: Deploy to Fly.io

1. **Setup**
   - [ ] Install Fly CLI: `curl -L https://fly.io/install.sh | sh`
   - [ ] Login: `fly auth login`
   - [ ] Initialize app: `fly launch`

2. **Configuration**
   - [ ] Review and edit `fly.toml` if needed
   - [ ] Set secrets:
   ```bash
   fly secrets set MONGODB_URI="mongodb+srv://..."
   fly secrets set ENV=production
   ```

3. **Deploy**
   - [ ] Deploy: `fly deploy`
   - [ ] Check status: `fly status`
   - [ ] View logs: `fly logs`
   - [ ] Verify health: `https://your-app.fly.dev/health`

### Option 3: Deploy with Docker

1. **Build Image**
   ```bash
   docker build -t kenya-info-api:latest .
   ```

2. **Push to Registry**
   ```bash
   docker tag kenya-info-api:latest your-registry/kenya-info-api:latest
   docker push your-registry/kenya-info-api:latest
   ```

3. **Deploy**
   - [ ] Configure container service (AWS ECS, GCP Cloud Run, etc.)
   - [ ] Set environment variables
   - [ ] Configure port mapping (8080)
   - [ ] Set up health checks
   - [ ] Deploy container

### Option 4: Deploy to VPS (Ubuntu)

1. **Server Setup**
   ```bash
   # Update system
   sudo apt update && sudo apt upgrade -y
   
   # Install Go
   wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
   
   # Or install Docker
   sudo apt install docker.io docker-compose -y
   ```

2. **Deploy Application**
   ```bash
   # Clone repository
   git clone https://github.com/Ismael-Njihia/Kenya-info-api.git
   cd Kenya-info-api
   
   # Set environment variables
   cp .env.example .env
   # Edit .env with production values
   
   # Build
   go build -o kenya-api cmd/api/main.go
   
   # Run (or use systemd service)
   ./kenya-api
   ```

3. **Setup systemd Service**
   - [ ] Create service file: `/etc/systemd/system/kenya-api.service`
   - [ ] Enable and start service
   - [ ] Configure nginx reverse proxy (optional)

## Post-Deployment

### Verification
- [ ] Health endpoint responding: `/health`
- [ ] Swagger documentation accessible: `/swagger/index.html`
- [ ] All API endpoints working
- [ ] Pagination working correctly
- [ ] Search functionality working
- [ ] CORS headers present
- [ ] Response times acceptable (<200ms for simple queries)

### Monitoring
- [ ] Set up uptime monitoring (UptimeRobot, Pingdom, etc.)
- [ ] Configure error alerting
- [ ] Set up log aggregation (if using cloud)
- [ ] Monitor database connections
- [ ] Track API response times

### Security
- [ ] HTTPS enabled
- [ ] Environment variables not exposed
- [ ] MongoDB credentials secure
- [ ] CORS configured appropriately
- [ ] Rate limiting configured (if needed)
- [ ] No sensitive data in logs

### Documentation
- [ ] Update README with production URL
- [ ] Document API base URL
- [ ] Share Swagger documentation URL
- [ ] Create API usage guide for consumers
- [ ] Document rate limits (if any)

## Rollback Plan

If deployment fails:

1. **Render/Fly.io**
   - Rollback to previous deployment in dashboard
   - Or redeploy previous Git commit

2. **Docker**
   - Deploy previous image version
   - Verify rollback successful

3. **VPS**
   - Restore previous binary
   - Restart service

## Maintenance

### Regular Tasks
- [ ] Monitor error logs daily
- [ ] Review performance metrics weekly
- [ ] Update dependencies monthly
- [ ] Backup database regularly
- [ ] Test disaster recovery plan quarterly

### Updates
- [ ] Test updates in staging environment first
- [ ] Create database backup before updates
- [ ] Deploy during low-traffic periods
- [ ] Monitor closely after deployment
- [ ] Have rollback plan ready

## Troubleshooting

### Common Issues

**Connection Errors**
- Verify MongoDB URI is correct
- Check IP whitelist in MongoDB Atlas
- Verify network connectivity

**High Response Times**
- Check database indexes
- Monitor MongoDB performance
- Check server resources
- Consider caching

**Build Failures**
- Verify all dependencies in go.mod
- Check Go version compatibility
- Review build logs

## Success Criteria

Deployment is successful when:
- ✅ Health endpoint returns 200
- ✅ All tests passing in production
- ✅ API responses within acceptable time
- ✅ No errors in logs
- ✅ Monitoring alerts configured
- ✅ Documentation updated
- ✅ Team notified of deployment

## Support Contacts

- **Technical Issues**: [Your contact]
- **Database Issues**: [DBA contact]
- **Infrastructure**: [DevOps contact]

---

**Last Updated**: 2025-10-27
