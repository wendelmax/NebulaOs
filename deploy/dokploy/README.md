# Dokploy Deployment for NebulaOS

[Dokploy](https://dokploy.com) is an open-source, self-hosted Platform as a Service (PaaS)
that simplifies deploying and managing applications using Docker and Traefik.
It serves as a free alternative to Heroku, Vercel, and Netlify.

This directory contains a Docker Compose file optimized for Dokploy deployment.

## Quick Start

### Prerequisites
- A VPS with at least 2GB RAM and 30GB disk
- Ubuntu 22.04+, Debian 12, or Fedora 40
- A domain name pointing to your server IP
- Ports 80, 443, and 3000 available

### 1. Install Dokploy

```bash
curl -sSL https://dokploy.com/install.sh | sh
```

After installation, access the Dokploy web UI at `http://<your-server-ip>:3000`.

### 2. Deploy NebulaOS

1. **Create a new project** in Dokploy
2. **Create a new service** → select **Compose** → **Docker Compose**
3. **Configure source:**
   - Provider: `Git` or `GitHub`
   - Repository: `https://github.com/wendelmax/NebulaOs`
   - Branch: `main`
   - Compose Path: `deploy/dokploy/docker-compose.dokploy.yml`
4. **Set environment variables:**
   - `JWT_SECRET` — a strong random secret for JWT signing
   - `POSTGRES_PASSWORD` — a strong database password
   - `VAULT_TOKEN` — Vault root token
   - `API_DOMAIN` — domain for the API (e.g., `api.nebula.example.com`)
   - `DASHBOARD_DOMAIN` — domain for the dashboard (e.g., `nebula.example.com`)
5. **Add domains** in the Dokploy Domains tab for each service
6. **Click Deploy**

### 3. Post-Deployment

After all services are running:

1. Access the dashboard at your configured domain
2. Log in with the default credentials:
   - Username: `admin`
   - Password: `admin`
3. Change the default password immediately

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `JWT_SECRET` | — | JWT signing secret (required) |
| `POSTGRES_DB` | `nebula` | PostgreSQL database name |
| `POSTGRES_USER` | `nebula` | PostgreSQL user |
| `POSTGRES_PASSWORD` | `password123` | PostgreSQL password |
| `VAULT_TOKEN` | `root` | Vault authentication token |
| `AWS_ENDPOINT` | `http://moto:4566` | AWS mock endpoint (moto) |
| `AWS_REGION` | `us-east-1` | AWS region |
| `API_DOMAIN` | `api.nebula.local` | API domain for Traefik routing |
| `DASHBOARD_DOMAIN` | `nebula.local` | Dashboard domain for Traefik routing |
| `VITE_API_URL` | — | Dashboard API URL (auto-detected) |

## Architecture

```
                      Users
                        |
                    [Traefik] (provided by Dokploy)
                   /         \
            [nebula-api]  [nebula-dashboard]
                  |
            [postgres]  [moto]  [vault]  [keycloak]
```

## Monitoring

Dokploy provides built-in monitoring for:
- CPU, memory, and network usage per container
- Real-time and historical logs
- Deployment history (last 10 deployments)
- Automatic SSL certificate renewal via Let's Encrypt

## Updating

To update NebulaOS to the latest version:

1. In Dokploy UI, go to your service
2. Click "Update" or trigger a new deployment
3. Dokploy will pull the latest code and rebuild the images

Or trigger a deployment via webhook (configured in Dokploy settings).
