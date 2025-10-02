# Playability Deployment Checklist

This document outlines all tasks required to deploy the Playability application to production.

---

## 1. Environment Configuration

### 1.1 Environment Variables Setup
- [ ] Create `.env.example` for frontend
  - [ ] `IGDB_CLIENT_SECRET`
  - [ ] `NUXT_PUBLIC_API_URL` (backend API URL for production)

- [ ] Create `.env.example` for backend
  - [ ] `DB_HOST`
  - [ ] `DB_PORT`
  - [ ] `DB_USER`
  - [ ] `DB_PASSWORD`
  - [ ] `DB_NAME`
  - [ ] `DB_SSLMODE` (require for production)
  - [ ] `JWT_SECRET` (must be cryptographically secure, 32+ chars)
  - [ ] `IGDB_ACCESS_TOKEN`
  - [ ] `CLAUDE_API_KEY`
  - [ ] `PORT` (default 8080)
  - [ ] `ENVIRONMENT` (development/staging/production)
  - [ ] `ALLOWED_ORIGINS` (comma-separated list for CORS)
  - [ ] `FRONTEND_URL` (for CORS configuration)

### 1.2 Production Environment Variables
- [ ] Generate secure JWT_SECRET (use `openssl rand -base64 32`)
- [ ] Obtain production IGDB credentials
- [ ] Obtain Claude API key with appropriate rate limits
- [ ] Configure production database credentials
- [ ] Set up separate staging environment variables

---

## 2. Security Hardening

### 2.1 Backend Security
- [ ] **CORS Configuration** (critical issue found)
  - [ ] Replace `AllowedOrigins: []string{"*"}` with specific production domain
  - [ ] Remove wildcard CORS in production (main.go:29)
  - [ ] Use environment variable for CORS origins
  - [ ] Validate origin headers properly

- [ ] **JWT Security**
  - [ ] Ensure JWT_SECRET is strong (32+ random bytes)
  - [ ] Consider rotating JWT secrets strategy
  - [ ] Implement token refresh mechanism
  - [ ] Add token revocation/blacklist for logout
  - [ ] Review JWT expiration time (currently 24h)
  - [ ] Add rate limiting to login/register endpoints

- [ ] **Password Security**
  - [ ] Verify bcrypt cost factor is appropriate (currently using DefaultCost=10)
  - [ ] Add password strength requirements
  - [ ] Implement account lockout after failed attempts

- [ ] **Input Validation**
  - [ ] Add request size limits
  - [ ] Implement rate limiting middleware
  - [ ] Validate all user inputs server-side
  - [ ] Add SQL injection protection verification
  - [ ] Sanitize all database queries

- [ ] **Cookie Security**
  - [ ] Ensure cookies are Secure in production only
  - [ ] Verify HttpOnly flag is set (already done)
  - [ ] Add SameSite attribute (Strict or Lax)
  - [ ] Consider using __Host- prefix for cookies

- [ ] **API Security**
  - [ ] Add request timeout middleware
  - [ ] Implement API rate limiting (per IP/user)
  - [ ] Add request ID tracking for debugging
  - [ ] Implement proper error handling (no stack traces in production)
  - [ ] Add security headers middleware (helmet equivalent)

- [ ] **HTTPS/TLS**
  - [ ] Configure TLS 1.3 minimum
  - [ ] Set up SSL/TLS certificates
  - [ ] Implement HSTS headers
  - [ ] Configure secure cipher suites

### 2.2 Frontend Security
- [ ] Enable CSP (Content Security Policy) headers
- [ ] Add XSS protection headers
- [ ] Configure X-Frame-Options
- [ ] Ensure all API calls use environment-based URLs (not hardcoded localhost)
- [ ] Implement proper error handling (no sensitive data leaks)
- [ ] Add Subresource Integrity (SRI) for CDN assets
- [ ] Configure secure cookie settings based on environment

### 2.3 Database Security
- [ ] Enable SSL/TLS for database connections
- [ ] Use least-privilege database user for application
- [ ] Create separate read-only user for analytics/reporting
- [ ] Implement database connection pooling with limits
- [ ] Enable database audit logging
- [ ] Set up automated database backups
- [ ] Encrypt sensitive data at rest (if applicable)
- [ ] Review and restrict database network access

### 2.4 Secrets Management
- [ ] Never commit `.env` files to repository (verify .gitignore)
- [ ] Use secrets manager (AWS Secrets Manager, HashiCorp Vault, etc.)
- [ ] Rotate API keys and secrets regularly
- [ ] Implement secret scanning in CI/CD
- [ ] Document secret rotation procedures

### 2.5 Dependency Security
- [ ] Run `pnpm audit` and fix vulnerabilities
- [ ] Run `go list -m all | nancy sleuth` for Go vulnerabilities
- [ ] Set up automated dependency updates (Dependabot/Renovate)
- [ ] Pin dependency versions in production
- [ ] Review all third-party package licenses

---

## 3. Database Setup

### 3.1 Migration System
- [ ] Create database migration tool/framework
  - [ ] Consider golang-migrate, goose, or Flyway
- [ ] Create initial migration from existing SQL files
  - [ ] `001_create_feature_support_enum.sql`
  - [ ] `002_create_users_table.sql`
  - [ ] `003_create_games_table.sql`
  - [ ] `004_create_reports_table.sql`
- [ ] Add foreign key constraints
  - [ ] reports.game_id → games.id
  - [ ] reports.user_id → users.id
- [ ] Add database indexes
  - [ ] games(name) for search performance
  - [ ] reports(game_id, created_at)
  - [ ] reports(user_id)
  - [ ] users(email) unique index
  - [ ] users(username) unique index
- [ ] Create database initialization script
- [ ] Add seed data for testing

### 3.2 Database Optimization
- [ ] Add connection pooling configuration
- [ ] Set appropriate timeout values
- [ ] Configure max connections
- [ ] Add prepared statement caching
- [ ] Review query performance
- [ ] Add query logging for slow queries

### 3.3 Database Backup & Recovery
- [ ] Set up automated daily backups
- [ ] Test backup restoration procedure
- [ ] Document recovery process
- [ ] Set up point-in-time recovery
- [ ] Configure backup retention policy

---

## 4. Docker & Containerization

### 4.1 Frontend Dockerfile
- [ ] Create multi-stage Dockerfile for Nuxt app
  - [ ] Stage 1: Dependencies installation
  - [ ] Stage 2: Build production bundle
  - [ ] Stage 3: Production runtime (minimal image)
- [ ] Optimize image size (use alpine base)
- [ ] Add health check endpoint
- [ ] Configure non-root user
- [ ] Add .dockerignore file

### 4.2 Backend Dockerfile
- [ ] Create multi-stage Dockerfile for Go API
  - [ ] Stage 1: Build binary
  - [ ] Stage 2: Runtime (scratch or alpine)
- [ ] Use CGO_ENABLED=0 for static binary
- [ ] Add health check endpoint
- [ ] Configure non-root user
- [ ] Add .dockerignore file

### 4.3 Database Container
- [ ] Create PostgreSQL Docker configuration
- [ ] Add initialization scripts
- [ ] Configure persistent volumes
- [ ] Set up proper user permissions

### 4.4 Docker Compose
- [ ] Create docker-compose.yml for local development
- [ ] Create docker-compose.production.yml
- [ ] Configure service dependencies
- [ ] Set up networking between services
- [ ] Add health checks
- [ ] Configure restart policies
- [ ] Set resource limits (CPU/memory)

### 4.5 Container Security
- [ ] Scan images for vulnerabilities
- [ ] Use official base images only
- [ ] Keep base images updated
- [ ] Don't run as root user
- [ ] Minimize attack surface (minimal images)

---

## 5. Application Configuration

### 5.1 Backend Configuration
- [ ] Make port configurable via environment variable
- [ ] Add graceful shutdown handling
- [ ] Implement structured logging (JSON format)
- [ ] Add log levels (DEBUG, INFO, WARN, ERROR)
- [ ] Configure CORS based on environment
- [ ] Add health check endpoint (`/health`, `/ready`)
- [ ] Add metrics endpoint for monitoring
- [ ] Implement request logging middleware
- [ ] Add panic recovery middleware
- [ ] Configure timeouts (read, write, idle)

### 5.2 Frontend Configuration
- [ ] Configure backend API URL from environment
- [ ] Set up SSR vs SSG decision
- [ ] Configure build optimizations
- [ ] Add robots.txt and sitemap.xml
- [ ] Set up meta tags for SEO
- [ ] Configure image optimization
- [ ] Add analytics (if required)
- [ ] Set up error tracking (Sentry, etc.)

### 5.3 Build Scripts
- [ ] Add `lint` script to package.json
- [ ] Add `typecheck` script to package.json
- [ ] Add `test` script to package.json
- [ ] Add `build:prod` script with optimizations
- [ ] Create backend build script
- [ ] Add pre-commit hooks (husky)

---

## 6. Testing

### 6.1 Backend Testing
- [ ] Add unit tests for handlers
- [ ] Add unit tests for database layer
- [ ] Add integration tests for API endpoints
- [ ] Test JWT authentication flow
- [ ] Test AI moderation integration
- [ ] Test external API integrations (IGDB, PCGamingWiki)
- [ ] Add database transaction tests
- [ ] Test error handling and edge cases
- [ ] Add load testing

### 6.2 Frontend Testing
- [ ] Add component unit tests
- [ ] Add E2E tests (Playwright/Cypress)
- [ ] Test authentication flows
- [ ] Test form validation
- [ ] Test API integration
- [ ] Add accessibility testing
- [ ] Test responsive design
- [ ] Cross-browser testing

### 6.3 Security Testing
- [ ] Run OWASP ZAP security scan
- [ ] Test SQL injection protection
- [ ] Test XSS protection
- [ ] Test CSRF protection
- [ ] Verify authentication/authorization
- [ ] Test rate limiting
- [ ] Penetration testing

---

## 7. CI/CD Pipeline

### 7.1 GitHub Actions Setup
- [ ] Create CI workflow
  - [ ] Lint frontend code
  - [ ] Typecheck frontend
  - [ ] Run frontend tests
  - [ ] Build frontend
  - [ ] Lint backend code (golangci-lint)
  - [ ] Run backend tests
  - [ ] Build backend binary
  - [ ] Security scanning
  - [ ] Dependency vulnerability scanning

- [ ] Create CD workflow
  - [ ] Build Docker images
  - [ ] Push to container registry
  - [ ] Deploy to staging
  - [ ] Run smoke tests
  - [ ] Deploy to production (with approval)

- [ ] Set up branch protection rules
- [ ] Configure deployment gates
- [ ] Add rollback mechanism

### 7.2 Container Registry
- [ ] Choose registry (Docker Hub, GitHub Container Registry, AWS ECR)
- [ ] Set up authentication
- [ ] Configure image tagging strategy (semantic versioning)
- [ ] Set up image scanning
- [ ] Configure retention policies

---

## 8. Infrastructure & Deployment

### 8.1 Choose Deployment Platform
- [ ] **Option A: Cloud Platform (Recommended)**
  - [ ] AWS (ECS/EKS, RDS, CloudFront)
  - [ ] Google Cloud (Cloud Run, Cloud SQL)
  - [ ] DigitalOcean (App Platform, Managed DB)
  - [ ] Railway.app (simple, all-in-one)
  - [ ] Render.com (simple, all-in-one)

- [ ] **Option B: Traditional VPS**
  - [ ] DigitalOcean Droplet
  - [ ] Linode
  - [ ] Hetzner

- [ ] **Option C: Serverless**
  - [ ] Vercel (frontend) + Railway (backend + DB)
  - [ ] Netlify (frontend) + Fly.io (backend + DB)

### 8.2 Database Hosting
- [ ] Managed PostgreSQL (AWS RDS, DigitalOcean, Supabase)
- [ ] Configure automatic backups
- [ ] Set up read replicas (if needed)
- [ ] Enable connection pooling (PgBouncer)
- [ ] Configure monitoring and alerts

### 8.3 Domain & DNS
- [ ] Purchase/configure domain name
- [ ] Set up DNS records
- [ ] Configure CDN (CloudFlare, AWS CloudFront)
- [ ] Set up SSL/TLS certificates (Let's Encrypt)
- [ ] Configure DNSSEC (optional)

### 8.4 Load Balancing & Scaling
- [ ] Set up load balancer (if using multiple instances)
- [ ] Configure auto-scaling rules
- [ ] Set up health checks
- [ ] Configure SSL termination

---

## 9. Monitoring & Observability

### 9.1 Logging
- [ ] Set up centralized logging (ELK, Grafana Loki, CloudWatch)
- [ ] Configure log aggregation
- [ ] Set up log retention policies
- [ ] Add structured logging to application
- [ ] Configure log levels per environment

### 9.2 Monitoring & Metrics
- [ ] Set up application monitoring (Prometheus, DataDog, New Relic)
- [ ] Monitor database performance
- [ ] Track API response times
- [ ] Monitor error rates
- [ ] Track resource usage (CPU, memory, disk)
- [ ] Set up uptime monitoring (UptimeRobot, Pingdom)

### 9.3 Alerting
- [ ] Configure alerts for critical errors
- [ ] Set up alerts for high response times
- [ ] Configure alerts for database issues
- [ ] Set up alerts for security events
- [ ] Configure on-call rotation (PagerDuty, etc.)

### 9.4 Error Tracking
- [ ] Set up error tracking (Sentry, Rollbar)
- [ ] Configure source maps for frontend
- [ ] Set up error notifications
- [ ] Configure error grouping and deduplication

### 9.5 Analytics & Performance
- [ ] Set up web analytics (if needed)
- [ ] Add performance monitoring (Web Vitals)
- [ ] Track API usage patterns
- [ ] Monitor AI API costs and usage

---

## 10. Documentation

### 10.1 Technical Documentation
- [ ] API documentation (OpenAPI/Swagger)
- [ ] Database schema documentation
- [ ] Architecture diagrams
- [ ] Deployment architecture diagram
- [ ] Authentication flow documentation
- [ ] Rate limiting documentation

### 10.2 Operations Documentation
- [ ] Deployment runbook
- [ ] Rollback procedures
- [ ] Incident response plan
- [ ] Database backup/restore procedures
- [ ] Secret rotation procedures
- [ ] Scaling procedures
- [ ] Monitoring and alerting guide

### 10.3 Developer Documentation
- [ ] Update README.md with production setup
- [ ] Local development setup guide
- [ ] Contributing guidelines
- [ ] Code style guide
- [ ] Testing guide
- [ ] Troubleshooting guide

---

## 11. Performance Optimization

### 11.1 Frontend Performance
- [ ] Enable Nuxt image optimization
- [ ] Configure CDN for static assets
- [ ] Implement lazy loading for images
- [ ] Add caching headers
- [ ] Minify CSS/JS
- [ ] Enable compression (gzip/brotli)
- [ ] Optimize bundle size
- [ ] Implement code splitting
- [ ] Add service worker for PWA (optional)

### 11.2 Backend Performance
- [ ] Implement response caching
- [ ] Add Redis for session/cache storage
- [ ] Optimize database queries
- [ ] Add database indexes
- [ ] Implement connection pooling
- [ ] Add HTTP response compression
- [ ] Optimize AI API calls (caching, batching)

### 11.3 Database Performance
- [ ] Analyze slow queries
- [ ] Add appropriate indexes
- [ ] Optimize table structure
- [ ] Configure query caching
- [ ] Set up read replicas if needed

---

## 12. Compliance & Legal

### 12.1 Data Privacy
- [ ] Review GDPR compliance (if applicable)
- [ ] Add privacy policy
- [ ] Add terms of service
- [ ] Implement data export feature
- [ ] Implement data deletion feature
- [ ] Add cookie consent banner (if needed)
- [ ] Document data retention policies

### 12.2 Accessibility
- [ ] Run accessibility audit (WCAG 2.1 AA)
- [ ] Test with screen readers
- [ ] Add ARIA labels where needed
- [ ] Ensure keyboard navigation works
- [ ] Test color contrast ratios
- [ ] Add skip navigation links

### 12.3 Content Moderation
- [ ] Review AI moderation effectiveness
- [ ] Add manual moderation tools
- [ ] Implement user reporting system
- [ ] Create content policy
- [ ] Set up moderation workflow

---

## 13. Pre-Launch Checklist

### 13.1 Final Security Review
- [ ] Complete security audit
- [ ] Verify all secrets are properly managed
- [ ] Test authentication/authorization flows
- [ ] Verify HTTPS is enforced
- [ ] Test CORS configuration
- [ ] Review all API endpoints for security

### 13.2 Performance Testing
- [ ] Run load tests
- [ ] Test under expected traffic
- [ ] Verify caching works correctly
- [ ] Test CDN configuration
- [ ] Verify database performance

### 13.3 Functionality Testing
- [ ] Test all critical user flows
- [ ] Verify external API integrations work
- [ ] Test AI moderation
- [ ] Verify email notifications (if any)
- [ ] Test error handling and recovery

### 13.4 Operational Readiness
- [ ] Verify monitoring and alerts are working
- [ ] Test backup and restore procedures
- [ ] Verify logs are being collected
- [ ] Test deployment process
- [ ] Verify rollback procedures work
- [ ] Create runbooks for common issues

---

## 14. Post-Launch

### 14.1 Immediate Post-Launch
- [ ] Monitor application logs and errors
- [ ] Watch performance metrics
- [ ] Monitor user feedback
- [ ] Be ready for hotfixes
- [ ] Monitor API rate limits and costs

### 14.2 Ongoing Maintenance
- [ ] Regular security updates
- [ ] Dependency updates
- [ ] Performance monitoring and optimization
- [ ] Regular backups verification
- [ ] Quarterly security audits
- [ ] Review and update documentation

---

## Priority Breakdown

### **Critical (Must Complete Before Launch)**
1. Environment configuration with secure secrets
2. Fix CORS wildcard vulnerability (main.go:29)
3. Database migrations and setup
4. Docker containerization
5. HTTPS/TLS configuration
6. Production deployment infrastructure
7. Security hardening (JWT, cookies, input validation)
8. Monitoring and error tracking
9. Database backups

### **High Priority (Complete Before Public Launch)**
1. Rate limiting on all endpoints
2. Comprehensive testing (unit, integration, E2E)
3. CI/CD pipeline
4. Error tracking and logging
5. API documentation
6. Performance optimization
7. Health check endpoints

### **Medium Priority (Complete Within First Month)**
1. Enhanced security features (token refresh, account lockout)
2. Load testing and optimization
3. Complete documentation
4. Accessibility compliance
5. Advanced monitoring and alerting

### **Low Priority (Nice to Have)**
1. Analytics integration
2. Advanced caching strategies
3. PWA features
4. Read replicas for database

---

## Estimated Timeline

- **Week 1-2**: Environment setup, security fixes, Docker configuration
- **Week 2-3**: Database setup, testing, CI/CD pipeline
- **Week 3-4**: Infrastructure setup, deployment, monitoring
- **Week 4-5**: Security audit, performance testing, documentation
- **Week 5-6**: Staging deployment, final testing, launch preparation
- **Week 6+**: Production launch, monitoring, iterations

---

## Notes

- This is a living document - update as tasks are completed
- Prioritize security and stability over features
- Test everything in staging before production
- Keep documentation updated throughout the process
- Have a rollback plan for every deployment
