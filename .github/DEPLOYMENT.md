# Deployment Setup

This document describes how to set up automatic deployment to Digital Ocean using GitHub Actions.

## Overview

The deployment workflow automatically:
1. Builds Docker images for frontend and backend
2. Pushes images to DockerHub
3. Deploys to Digital Ocean droplet by pulling the latest images

## Required GitHub Secrets

Configure these secrets in your GitHub repository settings (Settings → Secrets and variables → Actions):

### DockerHub Credentials

1. **`DOCKERHUB_USERNAME`**
   - Your DockerHub username
   - Example: `myusername`

2. **`DOCKERHUB_TOKEN`**
   - DockerHub access token (NOT your password)
   - Create at: https://hub.docker.com/settings/security
   - Click "New Access Token"
   - Give it a description (e.g., "GitHub Actions")
   - Copy the token and save it as this secret

### Digital Ocean Access

3. **`DO_DROPLET_IP`**
   - Your Digital Ocean droplet's IP address
   - Example: `159.89.123.456`

4. **`DO_DROPLET_USER`**
   - SSH username for your droplet
   - Default: `root`
   - If you created a non-root user, use that username

5. **`DO_SSH_PRIVATE_KEY`**
   - SSH private key for accessing your droplet
   - Generate a new SSH key pair:
     ```bash
     ssh-keygen -t ed25519 -C "github-actions" -f ~/.ssh/do_github_actions
     ```
   - Copy the private key:
     ```bash
     cat ~/.ssh/do_github_actions
     ```
   - Paste the entire output (including `-----BEGIN OPENSSH PRIVATE KEY-----` and `-----END OPENSSH PRIVATE KEY-----`) as this secret
   - Add the public key to your droplet:
     ```bash
     ssh-copy-id -i ~/.ssh/do_github_actions.pub user@your-droplet-ip
     ```

## Digital Ocean Droplet Setup

Your droplet needs to have the following installed:

1. **Docker**
   ```bash
   curl -fsSL https://get.docker.com -o get-docker.sh
   sudo sh get-docker.sh
   sudo usermod -aG docker $USER
   ```

2. **Docker Compose**
   ```bash
   sudo apt-get update
   sudo apt-get install docker-compose-plugin
   ```

3. **Project Directory**
   ```bash
   mkdir -p ~/playability
   ```

4. **Environment File**
   - Create `~/playability/.env` with all required environment variables:
     ```bash
     # Database
     DB_HOST=postgres
     DB_PORT=5432
     DB_USER=playability_user
     DB_PASSWORD=your_secure_password
     DB_NAME=playability
     POSTGRES_USER=playability_user
     POSTGRES_PASSWORD=your_secure_password
     POSTGRES_DB=playability

     # Backend
     JWT_SECRET=your_jwt_secret
     IGDB_ACCESS_TOKEN=your_igdb_token
     CLAUDE_API_KEY=your_claude_api_key

     # Frontend
     NUXT_API_URL=http://backend:8080
     NUXT_PUBLIC_API_URL=http://your-droplet-ip:8080

     # DockerHub
     DOCKERHUB_USERNAME=your_dockerhub_username
     ```

5. **Database Initialization Script**
   ```bash
   # Copy your database initialization script to the droplet
   mkdir -p ~/playability/database/scripts
   # Upload your init script to ~/playability/database/scripts/docker-entrypoint-initdb.sh
   ```

## Deployment Process

### Automatic Deployment

Deployment happens automatically on every push to the `main` branch:

```bash
git push origin main
```

### Manual Deployment

You can also trigger deployment manually from GitHub:
1. Go to Actions tab in your repository
2. Select "Build, Push, and Deploy to Digital Ocean"
3. Click "Run workflow"
4. Select the branch and click "Run workflow"

## Monitoring Deployment

1. Watch the workflow progress in the Actions tab
2. Check logs on your droplet:
   ```bash
   ssh user@your-droplet-ip
   cd ~/playability
   docker-compose -f docker-compose.prod.yml logs -f
   ```

## Troubleshooting

### Build Fails
- Check GitHub Actions logs for build errors
- Verify Dockerfile syntax
- Ensure all dependencies are properly specified

### Push to DockerHub Fails
- Verify `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN` are correct
- Ensure the DockerHub token has write permissions

### SSH Connection Fails
- Verify `DO_DROPLET_IP` is correct
- Ensure SSH private key is properly formatted (includes header/footer)
- Check that the public key is in `~/.ssh/authorized_keys` on the droplet
- Verify `DO_DROPLET_USER` exists on the droplet

### Deployment Fails
- SSH into your droplet and check Docker logs
- Verify `.env` file exists and contains all required variables
- Check disk space: `df -h`
- Check Docker status: `docker ps -a`

### Application Not Accessible
- Check if containers are running: `docker-compose -f docker-compose.prod.yml ps`
- Verify firewall settings allow ports 80 and 8080
- Check container logs for errors

## Rolling Back

If a deployment fails, you can roll back to a previous version:

```bash
ssh user@your-droplet-ip
cd ~/playability

# Pull a specific version (replace SHA with git commit hash)
docker pull your-dockerhub-username/playability-backend:SHA
docker pull your-dockerhub-username/playability-frontend:SHA

# Update docker-compose to use that version, then restart
docker-compose -f docker-compose.prod.yml up -d --force-recreate
```

## Security Best Practices

1. Never commit secrets to the repository
2. Use strong passwords for database credentials
3. Regularly rotate SSH keys and access tokens
4. Keep Docker and system packages updated
5. Use a firewall (UFW) to restrict access:
   ```bash
   sudo ufw allow 22/tcp
   sudo ufw allow 80/tcp
   sudo ufw allow 443/tcp
   sudo ufw enable
   ```
