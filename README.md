# Panovista L7 Agent Firewall (Evaluation Edition) - Core V3.0

**Zero-Trust, Stateless L7 Security Proxy for Enterprise AI Deployments.**

Panovista provides an offline, cryptographically locked architectural boundary for AI Agents and Model Context Protocol (MCP) servers. By dropping this stateless sidecar proxy into your Virtual Private Cloud (VPC), you instantly enforce strict L7 routing constraints, payload validation, and Article 12/PCI-DSS compliance logging without a single external ping to a vendor server.

## The 14-Day Cryptographic Trial

This public repository provides access to the **14-Day Evaluation Tier** of the V3.0 Core Engine. 

**SECURITY NOTICE:** Panovista operates on mathematically enforced Zero-Trust principles. This evaluation container requires a cryptographically signed Ed25519 token to boot. Exactly 14 days after the token's issue date, the internal failsafe will trigger, and the proxy will permanently lock itself to prevent unauthorized production use. 

## Quick Start Deployment

Deploying the Panovista proxy requires zero external dependencies. It runs completely offline in your local environment with a strictly bounded memory footprint (<20MB).

### 1. Create your `docker-compose.yml`

Copy the following configuration file into your target environment. 

```yaml
version: '3.8'

services:
  panovista-security-proxy:
    image: ghcr.io/panovista/agent-firewall:eval-v3
    container_name: panovista-eval
    ports:
      - "4321:4321"
    environment:
      - PANOVISTA_PORT=4321
      - TARGET_MCP_URL=http://your-internal-mcp-server:8080
      - PANOVISTA_MODE=sidecar
      # The 14-Day Cryptographic Evaluation Token
      - PANOVISTA_LICENSE=<INSERT_14_DAY_EVAL_TOKEN_HERE>
      - PANOVISTA_ALLOW_HEADLESS=TRUE
      - PANOVISTA_LOG_LEVEL=info
```

### 2. Boot the Firewall

Run the following command to pull the signed Panovista image and boot the engine:

```bash
docker compose up -d
```

### 3. Verify Telemetry & Orchestration Probes

Check your container logs to ensure the offline license verified and the strict L7 traffic parsing is active:

```bash
docker logs -f panovista-eval
```

You should instantly see the Phase 1 Ingress Stamp from our **Passive Telemetry Odometer** outputting to `stdout`:

```json
{"level":"info","tag":"PANOVISTA_METRIC","status":"boot","node_id":"panovista-eval","license_tier_claimed":"standard_vpc","uptime_seconds":0,"peak_concurrent_streams":0}
```

You can verify the container readiness for your orchestrator via our native health probes:

```bash
curl http://localhost:4321/health/ready
```

Your downstream MCP database is now shielded.

---

## Configuration Lexicon (V3.0)

Panovista's runtime behavior is controlled entirely via environment variables and local declarative JSON schemas. 

| Variable Name | Required | Default | Functional Description |
| :--- | :--- | :--- | :--- |
| `PANOVISTA_MODE` | **Yes** | `sidecar` | Sets layout architecture: `sidecar` (protecting one tool) or `gateway` (routing multiple tools). |
| `TARGET_MCP_URL` | **Yes** | *None* | The internal, isolated network URL of the raw upstream Model Context Protocol tool server. |
| `UPSTREAM_PROVIDER` | No | *None* | External LLM routing destination if injecting API keys (e.g., `openai`, `anthropic`). |
| `PROVIDER_API_KEY` | No | *None* | Securely injected upstream platform credential. |
| `PANOVISTA_PORT` | No | `4321` | Local port the proxy core listens on for incoming AI agent or IDE requests. |
| `SCHEMA_MOUNT_PATH` | No | `/etc/panovista` | Local directory containing declarative JSON schemas for Data Loss Prevention (DLP) rules. |
| `PANOVISTA_LICENSE` | **Yes** | *None* | The offline, cryptographically signed token string dictating contract compliance. |

---

## Upgrading to Sovereign Enterprise

The Panovista Sovereign Enterprise tier transitions your infrastructure to our strictly-versioned private registry, removes the cryptographic time-bomb, and unlocks high-margin compliance modules. Sovereign builds include:
* Native compilation using the FIPS 140-3 validated cryptographic toolchain (`GOEXPERIMENT=boringcrypto`).
* Sequential tamper-evident log-chaining algorithms.
* Native hardware security module (HSM) integrations via PKCS#11.

To upgrade your evaluation environment to a permanent enterprise license, contact our sales engineering team to receive your scoped Enterprise Token and GHCR Registry Access Keys.
