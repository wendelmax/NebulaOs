#!/bin/bash

# NebulaOS Platform Orchestrator
# This script automates the startup of all NebulaOS requirements.

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

# Load .env if present
if [ -f .env ]; then
  set -a
  source .env
  set +a
fi

NEBULA_DOMAIN="${NEBULA_DOMAIN:-nebula.local}"
NEBULA_API_DOMAIN="${NEBULA_API_DOMAIN:-api.nebula.local}"
NEBULA_API_PORT="${NEBULA_API_PORT:-8000}"
SEED_ADMIN_USERNAME="${SEED_ADMIN_USERNAME:-admin}"
SEED_ADMIN_PASSWORD="${SEED_ADMIN_PASSWORD:-admin}"

echo -e "${BLUE}== NebulaOS Enterprise Orchestrator ==${NC}"

# Check for Docker
if ! [ -x "$(command -v docker)" ]; then
  echo -e "${RED}Error: docker is not installed.${NC}" >&2
  exit 1
fi

# Check for Docker Compose
if ! [ -x "$(command -v docker-compose)" ]; then
  echo -e "${RED}Error: docker-compose is not installed.${NC}" >&2
  exit 1
fi

echo -e "${GREEN}[1/3] Starting Core Infrastructure (Postgres, NATS, Traefik)...${NC}"
docker-compose -f deploy/local/docker-compose.yml up -d postgres nats traefik localstack vault

echo -e "${GREEN}[2/3] Building and Starting NebulaOS API & Dashboard...${NC}"
docker-compose -f deploy/local/docker-compose.yml up -d --build nebula-api nebula-dashboard

echo -e "${GREEN}[3/3] Verifying Health...${NC}"
max_retries=30
count=0
until $(curl --output /dev/null --silent --head --fail "http://localhost:${NEBULA_API_PORT}/health"); do
    printf '.'
    sleep 2
    count=$((count+1))
    if [ $count -eq $max_retries ]; then
      echo -e "${RED}\nError: Nebula API failed to start in time.${NC}"
      exit 1
    fi
done

echo -e "\n${GREEN}NebulaOS is UP and Running!${NC}"
echo -e "${BLUE}Dashboard:${NC} http://${NEBULA_DOMAIN} (ensure it points to 127.0.0.1 in /etc/hosts)"
echo -e "${BLUE}API Health:${NC} http://${NEBULA_API_DOMAIN}/health"
echo -e "${BLUE}Initial Credentials:${NC} ${SEED_ADMIN_USERNAME} / ${SEED_ADMIN_PASSWORD}"
