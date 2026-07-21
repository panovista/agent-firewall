[![Docker Pulls](https://img.shields.io/badge/Docker-ghcr.io%2Fpanovista%2Fagent--firewall-blue?logo=docker)](https://github.com/panovista/agent-firewall/pkgs/container/agent-firewall)
[![Release](https://img.shields.io/github/v/release/panovista/agent-firewall?color=green)](https://github.com/panovista/agent-firewall/releases)
[![Security Policy](https://img.shields.io/badge/Security-SECURITY.md-red)](SECURITY.md)

# Panovista L7 Agent Firewall (Evaluation Edition) - Core V3.0

**[🌐 Visit panovista.io](https://panovista.io)** | **Zero-Trust, Stateless L7 Security Proxy for Enterprise AI Deployments.**

Panovista provides an offline, cryptographically locked architectural boundary for AI Agents and Model Context Protocol (MCP) servers. By dropping this stateless sidecar proxy into your network, you instantly enforce strict L7 routing constraints, payload validation, and compliance logging without a single external ping to an outside vendor server.

## The Frictionless 14-Day Evaluation

This public repository provides access to the 14-Day Evaluation Tier of the V3.0 Core Engine. 

Unlike older versions, this evaluation container **does not require an upfront token to boot**. Instead, it uses a Frictionless First-Boot mechanism. The moment the container launches on your network, it securely creates an internal cryptographic timestamp file within your mounted state volume. The proxy will seamlessly route traffic for 14 calendar days. On day 15, the evaluation window closes, and the proxy will permanently lock itself until a valid 12-month enterprise license token is provided.

State Persistence: Panovista utilizes a local volume at /var/lib/panovista to track the evaluation trial period. Ensure this directory is mapped to a persistent host volume in your docker-compose configuration. Deleting this directory or failing to map it will reset your evaluation window.

---

## Quick Start Deployment

Deploying the Panovista proxy requires zero external dependencies. It runs completely offline in your local environment with a strictly bounded memory footprint.

### 1. Prepare Your Host Directories
Create an empty local directory named `state` to house the persistent evaluation trial clock metadata:

```bash
mkdir -p state
```

### 2. Create Your `docker-compose.yml`

Copy the following streamlined configuration file into your target environment folder.

```yaml
version: '3.8'

services:
  panovista-security-proxy:
    image: ghcr.io/panovista/agent-firewall:latest
    container_name: panovista-core-v3
    ports:
      - "8080:8080"
    volumes:
      # Mandatory persistent directory for tracking the 14-day evaluation clock
      - ./state:/var/lib/panovista  # Ensures the trial.lock file persists across restarts
      - ./panovista-config:/etc/panovista/policies:ro
    environment:
      - PANOVISTA_PORT=8080
      - TARGET_MCP_URL=http://your-internal-mcp-server:80
      - PANOVISTA_ENV_ID=eval-vpc-local
      - PANOVISTA_ALLOW_HEADLESS=TRUE

      # ==========================================================
      # 14-DAY TRIAL ACTIVE BY DEFAULT
      # To upgrade to the 12-Month Paid Tier, uncomment below:
      # - PANOVISTA_LICENSE_TOKEN=pv_lic_YOUR_TOKEN_HERE
      # ==========================================================
```

### 3. Boot the Firewall

Run the following command to pull the signed Panovista image and boot the engine:

```bash
docker compose up -d
```

### 4. Verify Telemetry & Trial Status

Check your container logs to ensure the evaluation clock initialized, the persistent storage mounted, and the strict L7 traffic parsing is active:

```bash
docker compose logs panovista-security-proxy
```

You should see our structured Passive Telemetry metrics outputting directly to `stdout`, confirming a successful local first boot:

```text
panovista-core-v3 | 2026/07/10 12:19:25 [*] No license token provided. Initiating Evaluation Tier checks...
panovista-core-v3 | 2026/07/10 12:19:25 ⏱️ FIRST BOOT DETECTED. 14-Day Evaluation Clock Started.
panovista-core-v3 | 2026/07/10 12:19:25 ⚠️ EVALUATION MODE ACTIVE. 14 Days Remaining.
panovista-core-v3 | 2026/07/10 12:19:25 {"level":"info","tag":"PANOVISTA_METRIC","status":"boot","node_id":"8ea012c93647","license_tier_claimed":"core_v3","uptime_seconds":0,"peak_concurrent_streams":0}
panovista-core-v3 | 2026/07/10 12:19:25 🚀 Panovista Core V3 active on port 8080
```

Verify the container readiness for your orchestrator via our native health probes:

```bash
# Returns 200 OK instantly if the HTTP server socket loop is up
curl http://localhost:8080/health/live

# Returns 200 OK only after trial validation passes or license verifies
curl http://localhost:8080/health/ready
```

---

## Configuration Lexicon (V3.0 Core)

Panovista's runtime behavior is controlled cleanly via environment variables and local volume mappings.

| Variable Name | Required | Default | Functional Description |
| :--- | :--- | :--- | :--- |
| `TARGET_MCP_URL` | Yes | *None* | The internal network URL of the raw upstream tool server or LLM backend being shielded. |
| `PANOVISTA_ENV_ID` | Yes | *None* | Ties binary execution to a specific VPC or environment ID for isolation matching. |
| `PANOVISTA_ALLOW_HEADLESS` | Yes | *None* | Bypasses the anti-automation security lock for cloud orchestration engines (`TRUE`). |
| `PANOVISTA_PORT` | No | `8080` | Local port the proxy core listens on for incoming AI agent or client application requests. |
| `PANOVISTA_LICENSE_TOKEN` | No | *None* | The offline, cryptographically signed Ed25519 token string that unlocks the unlimited tier. |

---

## Instant Attack Deflection Test

Prove Panovista is actively shielding your backend parameters. Run a standard curl command attempting a Layer 4 Server-Side Request Forgery (SSRF) bypass through the proxy:

```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{
        "jsonrpc": "2.0", 
        "id": "verification-test-001", 
        "method": "chat", 
        "params": {
            "name": "agent", 
            "arguments": {
                "target": "http://169.254.169.254/latest/meta-data/"
            }
        }
      }'
```

### Expected Defensive Response
The proxy instantly intercepts the forbidden cloud infrastructure address, terminates the transaction safely before it breaches your perimeter, and returns a sanitized JSON-RPC error payload:

```json
{
  "jsonrpc": "2.0",
  "id": "verification-test-001",
  "error": {
    "code": -32001,
    "message": "Security Policy Violation" 
  }
}
```

---

## Upgrading to Sovereign Enterprise

The Panovista Sovereign Enterprise tier transitions your infrastructure to our strictly-versioned private registry, removes the 14-day trial clock file limitations, and unlocks unlimited enterprise scale.

Sovereign builds include:
* Native compilation using a FIPS 140-3 validated cryptographic toolchain.
* Hardcoded defensive thresholds: 5MB egress volume choking, 300 req/min rate limiting, and a 20-layer deep JSON recursion filter.
* Injectable PEM-encoded Public Key certificates (`JWT_PUBLIC_KEY`) for full IdP identity access verification.

To upgrade your evaluation environment to a permanent enterprise license, visit **[panovista.io](https://panovista.io)** or contact our sales engineering team at **[ian.ayliffe@panovistamarketing.com](mailto:ian.ayliffe@panovistamarketing.com)**  to receive your scoped Cryptographic Enterprise Token.







