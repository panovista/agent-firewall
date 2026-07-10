# =======================================================
# PANOVISTA AUTOMATED DEPLOYMENT SCRIPT
# =======================================================

$IMAGE_NAME = "ghcr.io/panovista/agent-firewall:latest"
$VERSION_TAG = "ghcr.io/panovista/agent-firewall:v3.0.0"

Write-Host "[*] Starting Panovista Core Build Process..." -ForegroundColor Cyan

# 1. Build the Docker Image locally
Write-Host "[*] Building Docker container..." -ForegroundColor Yellow
docker build -t $IMAGE_NAME -t $VERSION_TAG .

# 2. Push the Images to GitHub Container Registry
Write-Host "[*] Pushing to GitHub Registry..." -ForegroundColor Yellow
docker push $IMAGE_NAME
docker push $VERSION_TAG

Write-Host "[+] Deployment Complete! The public can now pull the latest version." -ForegroundColor Green