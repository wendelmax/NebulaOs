# NebulaOS Platform Orchestrator (Windows)
# This script automates the startup of all NebulaOS requirements.

param(
    [string]$envFile = ".env"
)

# Load .env if present
if (Test-Path $envFile) {
    Get-Content $envFile | ForEach-Object {
        if ($_ -match "^\s*([^#=]+)=(.*)$") {
            $key = $matches[1].Trim()
            $val = $matches[2].Trim()
            Set-Item -Path "env:$key" -Value $val
        }
    }
}

$NebulaDomain = $env:NEBULA_DOMAIN
if (-not $NebulaDomain) { $NebulaDomain = "nebula.local" }

$NebulaApiDomain = $env:NEBULA_API_DOMAIN
if (-not $NebulaApiDomain) { $NebulaApiDomain = "api.nebula.local" }

$NebulaApiPort = $env:NEBULA_API_PORT
if (-not $NebulaApiPort) { $NebulaApiPort = "8000" }

$SeedAdminUser = $env:SEED_ADMIN_USERNAME
if (-not $SeedAdminUser) { $SeedAdminUser = "admin" }

$SeedAdminPass = $env:SEED_ADMIN_PASSWORD
if (-not $SeedAdminPass) { $SeedAdminPass = "admin" }

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
docker-compose -f deploy/local/docker-compose.yml up -d postgres nats traefik moto vault

Write-Host "[2/3] Building and Starting NebulaOS API & Dashboard..." -ForegroundColor Green
docker-compose -f deploy/local/docker-compose.yml up -d --build nebula-api nebula-dashboard

Write-Host "[3/3] Verifying Health..." -ForegroundColor Green
$maxRetries = 30
$count = 0
$success = $false

while ($count -lt $maxRetries) {
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:$NebulaApiPort/health" -Method Get -UseBasicParsing -ErrorAction SilentlyContinue
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
    Write-Host "Dashboard: http://$NebulaDomain" -ForegroundColor Blue
    Write-Host "API Health: http://$NebulaApiDomain/health" -ForegroundColor Blue
    Write-Host "Initial Credentials: $SeedAdminUser / $SeedAdminPass" -ForegroundColor Blue
    Write-Host "NOTE: Ensure $NebulaDomain and $NebulaApiDomain point to 127.0.0.1 in C:\Windows\System32\drivers\etc\hosts" -ForegroundColor Yellow
} else {
    Write-Host "`nError: Nebula API failed to start in time." -ForegroundColor Red
    exit 1
}
