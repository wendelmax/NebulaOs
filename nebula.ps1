# NebulaOS Platform Orchestrator (Windows)
# This script automates the startup of all NebulaOS requirements.

Write-Host "== NebulaOS Enterprise Orchestrator (Windows) ==" -ForegroundColor Blue

# Check for Docker
if (!(Get-Command docker -ErrorAction SilentlyContinue)) {
    Write-Host "Error: docker is not installed or not in PATH." -ForegroundColor Red
    exit 1
}

# Check if Docker is running
docker ps > $null 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "Error: Docker is not running. Please start Docker Desktop." -ForegroundColor Red
    exit 1
}

Write-Host "[1/3] Starting Core Infrastructure (Postgres, NATS, Traefik)..." -ForegroundColor Green
docker-compose -f deploy/local/docker-compose.yml up -d postgres nats traefik localstack vault

Write-Host "[2/3] Building and Starting NebulaOS API & Dashboard..." -ForegroundColor Green
docker-compose -f deploy/local/docker-compose.yml up -d --build nebula-api nebula-dashboard

Write-Host "[3/3] Verifying Health..." -ForegroundColor Green
$maxRetries = 30
$count = 0
$success = $false

while ($count -lt $maxRetries) {
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:8000/health" -Method Get -UseBasicParsing -ErrorAction SilentlyContinue
        if ($response.StatusCode -eq 200) {
            $success = $true
            break
        }
    } catch {
        # API not ready yet or port not open
    }
    Write-Host "." -NoNewline
    Start-Sleep -Seconds 5
    $count++
}

if ($success) {
    Write-Host "`nNebulaOS is UP and Running!" -ForegroundColor Green
    Write-Host "Dashboard: http://nebula.local" -ForegroundColor Blue
    Write-Host "API Health: http://api.nebula.local/health" -ForegroundColor Blue
    Write-Host "Initial Credentials: admin / admin" -ForegroundColor Blue
    Write-Host "NOTE: Ensure nebula.local and api.nebula.local point to 127.0.0.1 in C:\Windows\System32\drivers\etc\hosts" -ForegroundColor Yellow
} else {
    Write-Host "`nError: Nebula API failed to start in time." -ForegroundColor Red
    exit 1
}
