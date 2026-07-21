
# Panovista Proxy: Enterprise DevOps Compliance & Tuning Guide (V3.0)

Panovista operates as a high-performance, stateless Layer 7 security proxy. Version 3.0 introduces a frictionless **Product-Led Growth (PLG) Evaluation Tier** that shifts automatically into an **Offline Cryptographic Sovereign Tier** upon injection of an enterprise license token. 

When first deployed without a token, the proxy functions in **Evaluation Mode**—traffic is seamlessly routed while full L7 threat detection and telemetry are logged for discovery. This document outlines the compliance configurations, hardcoded mitigation limits, and red team verification strategies required to move your deployment into production.

---

## The Onboarding Motion: 14-Day Evaluation Tier

By default, the Panovista container boots in Evaluation Mode without requiring an upfront license token. 

To prevent trial-clock evasion across container restarts, the proxy enforces a persistent state baseline via a mapped directory. The container tracks its first-boot timestamp locally on the host machine. The proxy seamlessly routes traffic for 14 calendar days. On day 15, the proxy initiates a hard **Fail-Closed** state, drops all traffic, and logs an expiration fault until a valid 12-month cryptographic token is applied.

---

## Compliance Deployment Playbooks

### Playbook A: GDPR & Privacy (Light Compliance)
**Goal:** Prevent Customer PII from reaching external Large Language Models (LLMs) without breaking automated agent workflows.
**Action Required:** Minimal configuration overhead.
* Leave the `rbac_rules` array empty or set to `"allowed_roles": ["*"]`.
* Add custom JSON `rules` targeting specific semantic intents or data patterns (e.g., `pii_email`, `home_address`).
* The proxy utilizes Format-Preserving Tokenization (FPT) to swap out sensitive PII with secure tokens in volatile memory. Upstream workflows remain completely unblocked while data privacy boundaries are maintained transparently.

### Playbook B: SOC 2 Type II (Standard Enterprise)
**Goal:** Enforce Role-Based Access Control (RBAC) and maintain tamper-evident audit trails for security auditing.
**Action Required:** Gradual operational tightening.
* **Day 1-7 (Discovery):** Run the proxy using `"action": "evaluate"` in your RBAC rules. Monitor structured JSON `stdout` metrics in your centralized SIEM engine (Splunk/Datadog) to discover exactly which corporate roles are calling specific internal backend utilities or tools.
* **Day 8 (Enforcement):** Update your local system configuration to define explicit `"allowed_roles"` based on your collected telemetry. 
* **Cryptographic Identity Verification:** Set the environment variable `JWT_PUBLIC_KEY` with your corporate identity provider's (IdP) PEM certificate to strictly cryptographically verify all inbound user identity tokens.
* **Cobalt Strike Mitigation:** Once `JWT_PUBLIC_KEY` is active, your IdP *must* be configured to include an explicit `exp` (expiration) claim in the JWT payload. To mitigate advanced threat persistence or credential reuse, Panovista automatically drops any payload with an expired or missing `exp` claim, downgrading the actor to an `unauthenticated_agent`.
* Change the rule `"action"` from `"evaluate"` to `"block"` and restart the container to instantly lock down unauthorized lateral movement.

### Playbook C: FedRAMP & DoD IL5/IL6 (Strict Compliance)
**Goal:** Prevent all unauthorized external boundary transmission of Controlled Unclassified Information (CUI) and guarantee FIPS-validated cryptography.
**Action Required:** Absolute Zero-Trust Topology.
* The container is natively compiled with a modern Go toolchain enforcing FIPS-compliant symbols and executes entirely as an unprivileged container user (UID 10001), satisfying NIST SP 800-53 (SC-13) system protections.
* Define explicit `"denied_roles": ["*"]` for all highly destructive tools (e.g., database drops, key generation) unless explicitly whitelisted via targeted rules.
* Ensure all sensitive configuration parameters (`JWT_PUBLIC_KEY`, `TARGET_MCP_URL`) are injected securely via AWS Secrets Manager, host environment configurations, or HashiCorp Vault. 
* **Anti-Automation Guardrail:** For containerized cloud orchestration (Kubernetes / ECS / Docker), you must explicitly pass the environment variable `PANOVISTA_ALLOW_HEADLESS=TRUE`. If this flag is missing on a headless system device, the firewall assumes it is running in an unauthorized automated exploitation environment and will intentionally crash on boot (Exit Code 1) to protect its internal configuration.

---

## Active L7 Threat Mitigations & Operational Thresholds

To maintain high availability while fulfilling strict compliance requirements (FIPS 140-3, Zero-Trust Architecture), Panovista actively enforces the following network mitigations.

### Active Threat Mitigations
* **Recursive De-obfuscation:** Intercepts payloads actively attempting to bypass standard Regex checks via nested encoding chains (URL -> Base64 -> Hexadecimal).
* **Greedy Decoder Patch:** Uses explicit `utf8.Valid()` entropy checks to prevent Base64 false-positive unravelling (Zero-Day defense).
* **Homoglyph Normalization:** Forces Cyrillic, Zalgo, and non-standard Unicode variations into flattened ASCII before inspection via `unicode/norm.NFKC`.
* **L4 SSRF Filtering:** Silently drops any payload attempting to route to internal cloud infrastructure or metadata endpoints (e.g., `169.254.169.254`, `127.0.0.1`, `10.x.x.x`, and Hex/Decimal IP variants).

### Production Hardcoded Guardrails
The following thresholds are strictly hardcoded in the Go proxy runtime to maintain a defensive Zero-Trust posture. If automated AI agents or client orchestrators exceed these limits, the firewall will intentionally sever the transaction. **These are architectural safety features, not application bugs.**

| Defense Mechanism | Hardcoded Threshold | Behavioral Trigger | System Resolution |
| :--- | :--- | :--- | :--- |
| **Token Bucket Rate Limiting** | `300 requests / min` | Tracks incoming client IP. Protects the enterprise against "Denial of Wallet" API billing exhaustion. | Returns `HTTP 429 Too Many Requests`. |
| **Egress Volume Choking** | `5,242,880 Bytes (5MB)` | Triggers if the downstream backend attempts to return a massive payload. Prevents catastrophic Intellectual Property (IP) data dumps. | Returns `HTTP 413 Payload Too Large`. |
| **Cryptographic Replay Window** | `30 Seconds` | Caches the SHA-256 hash of inbound JSON-RPC IDs. Blocks intercepted and replayed credential attacks. | Returns `HTTP 409 Conflict`. |
| **JSON Recursion Limiter** | `20 Layers Deep` | Triggers if an agent passes a deeply nested JSON object. Defends the proxy's volatile memory against payload crash exploits (DoS). | Returns `HTTP 422 Unprocessable Entity`. |
| **Trial Clock Expiration** | `14 Days from Boot` | Evaluates file creation metadata on the mounted storage volume. Restricts unauthorized extended evaluation. | Fail-Closed. Container terminates (`Exit Code 1`). |
| **JWT Expiration Lock** | `Mandatory 'exp' Claim` | Enforces a fail-closed Identity Access configuration. If the client passes an authorization token missing a definitive expiration timestamp, the proxy refuses to route it. | Returns `HTTP 403 Forbidden` / Downgrades role. |

> **Frontend Developer Note (Replay Attack Prevention):**
> Due to the 30-second Cryptographic Replay Window, all frontend clients **must** generate a unique, randomized `id` field in their JSON-RPC payload for every single request. Static payloads sent sequentially (e.g., hardcoding `"id": 1`) will match existing SHA-256 cache signatures and will be aggressively dropped as a replay attack.
> 
> **The Quick Fix (JavaScript / TypeScript):**
> Do not use static IDs. Instead, use the native runtime API to generate a UUID on the fly:
> ```javascript
> const payload = { 
>     jsonrpc: "2.0", 
>     id: crypto.randomUUID(), // Generates a unique ID instantly
>     method: "chat", 
>     params: { /* your payload */ } 
> };
> ```

---

## Red Team & Penetration Testing Notes

If your internal security team or an external auditor is pen-testing the Panovista V3.0 architecture using LLM vulnerability frameworks (such as Promptfoo), take note of the following tactical operational guardrails.

### 1. The Obfuscation Gauntlet (De-obfuscation)
* **The Vector:** Attackers often attempt to sneak malicious payloads (like SSRF metadata addresses) past firewalls by nesting them inside combinations of URL encoding, Hexadecimal, or Base64 strings.
* **The Defense:** Panovista uses **Heuristic Transcoding**. It will automatically unravel up to 3 layers of nested encoding in volatile memory, normalize Unicode variations into standard ASCII via `unicode/norm.NFKC`, scan the flattened payload, and block the packet if a violation is found. 
* **Testing Note:** Standard encoding bypass techniques will fail. Testers do not need to write manual regex rules for encoded strings; the proxy handles the decoding transparently in memory.

### 2. Automated Fuzzing & The Replay Trap
* **The Vector:** Automation tools typically send a rapid sequence of test payloads while keeping the JSON envelope structure identical.
* **The Defense:** The proxy caches the SHA-256 hash of every inbound JSON-RPC request identifier for 30 seconds.
* **Testing Note:** If an automated scanner runs with a static or sequential ID structure (e.g., hardcoding `"id": 1` across an entire fuzzing suite), the proxy will drop every single concurrent request after the first one as a **Cryptographic Replay Attack (HTTP 409 Conflict)**. Testers *must* configure their test fixtures to generate a dynamic UUID on every individual request loop.

### 3. Exfiltration Choking (Data Dumps)
* **The Vector:** Prompt injection or downstream LLM hijacking designed to dump massive internal corporate databases or knowledge bases out through the network perimeter.
* **The Defense:** Strict egress volume choking is hardcoded at exactly **5,242,880 Bytes (5MB)**.
* **Testing Note:** If a red team injection attempt successfully tricks the backend into returning a massive data dump, the proxy will instantly sever the socket transaction and return an `HTTP 413 Payload Too Large`. The raw data never clears the boundary.

### 4. Algorithmic Complexity Traps (JSON DoS)
* **The Vector:** Sending an extensively nested, recursive JSON object designed to exhaust the proxy's runtime memory allocation and crash the container infrastructure (Denial of Service).
* **The Defense:** The `scrubJSON` parser tracks exact structural recursion depth.
* **Testing Note:** Any payload that penetrates deeper than **20 layers** triggers an immediate `HTTP 422 Unprocessable Entity` fail-closed truncation block.

### V3.0 Red Team Execution Matrix

| Attack Vector | Target Function | Expected Firewall Action / Fail-Closed State |
| :--- | :--- | :--- |
| **Chained Encoding** | `analyzeObfuscation` | Hard Block (`[MALICIOUS_OBFUSCATION_BLOCKED]`) |
| **JSON Depth > 20** | `scrubJSON` | Connection Severed (`HTTP 422 Unprocessable Entity`) |
| **Static Payload Hash** | `proxyHandler` | Request Dropped (`HTTP 409 Conflict` / Replay Window) |
| **5MB+ Upstream Response**| `proxyHandler` | Connection Severed (`HTTP 413 Payload Too Large`) |
| **Cyrillic/Zalgo Spoofing** | `unicode/norm` | Flattens to ASCII -> Evaluates against SSRF rules |
| **> 300 Requests / Min** | `isRateLimited` | Volatile IP cache lock (`HTTP 429 Too Many Requests`) |
| **Container Tampering** | `enforceTrialOrLicense` | Clock manipulation or volume unmounting results in instant boot crash (`Exit Code 1`) |
| **SIGTERM under Load** | `main` (Signal Hook) | Final `"status": "shutdown"` metric printed to standard out |

---

## Integration & Troubleshooting

### Mandatory Deployment Parameters (V3.0)

To successfully spin up the Panovista container without triggering a fail-closed boot sequence, your deployment manifest must map the storage volume and define the following environment blocks:

| Configuration Parameter | Type | Purpose | Example / Parity |
| :--- | :--- | :--- | :--- |
| `volumes:` | Docker Mount | Mounts host directory to track the trial window persistently. | `./state:/var/lib/panovista` |
| `PANOVISTA_ALLOW_HEADLESS` | Env Var | Bypasses the anti-automation lock for cloud orchestration. | `TRUE` |
| `PANOVISTA_ENV_ID` | Env Var | Ties binary execution to a specific VPC or environment ID. | `vpc-nexus-prod-99` |
| `TARGET_MCP_URL` | Env Var | The downstream destination where scrubbed traffic is routed. | `http://mock-mcp-backend:80` |
| `PANOVISTA_LICENSE_TOKEN` | Env Var | **(Optional)** The cryptographically signed Ed25519 12-month license key. Unlocks full Sovereign Enterprise Tier. | `pv_lic_eyJjdXN0...` |
| `JWT_PUBLIC_KEY` | Env Var | **(Optional)** Injectable PEM-encoded public key certificate for enterprise identity verification. | `-----BEGIN PUBLIC KEY-----...` |

### 1. DevOps Health & Readiness Probes
Panovista provides native endpoints for container orchestration (Kubernetes, AWS ECS, Docker Compose) to monitor container readiness without triggering the DLP engine:

* **Liveness Probe (`GET /health/live`):** Returns `HTTP 200 OK` instantly once the Go runtime web server is initialized and accepting socket connections.
* **Readiness Probe (`GET /health/ready`):** Returns `HTTP 200 OK` only after the 14-day evaluation validation completes or the Ed25519 license token is successfully verified in memory. Returns `HTTP 503 Service Unavailable` if evaluation states or required internal configurations are missing or faulted.

### 2. Standardized Error Handling (Frontend & Clients)
When Panovista intercepts a malicious payload, rate limit violation, or unauthorized access attempt, it severs the connection to the downstream backend and returns a standardized JSON-RPC 2.0 error format directly to the client. 

Frontend applications should monitor the structural **HTTP Status Code** and the custom `-32001` JSON-RPC error block to handle security mitigations gracefully.

**Example Block Response:**
```json
{
  "jsonrpc": "2.0",
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "error": {
    "code": -32001,
    "message": "Security Policy Violation" 
  }
}
```

*Note: The internal message string will dynamically update to reflect the explicit Layer 7 block condition encountered (e.g., "Security Policy Violation: Cryptographic Replay Attack", "Denial of Service Prevention: Rate Limit Exceeded").*

---

## Enterprise Quality Assurance & Verification

To maintain strict Zero-Trust compliance, the Panovista Proxy is subjected to an automated Continuous Integration (CI/CD) testing pipeline prior to any container artifact creation. Our internal testing architecture guarantees:
* **Cryptographic Verification:** Every build is mathematically verified via automated integration frameworks to ensure FIPS-compliant cryptographic linkages are active.
* **Obfuscation Gauntlets:** The DLP de-obfuscation engine is rigorously tested against nested encoding attacks (Base64, Hex, URL) and Cyrillic homoglyph spoofing to ensure malicious payloads are met with a hard Zero-Trust block state.
* **No Testing Backdoors:** Panovista does not contain bypass flags or development backdoors. The production container artifact is the exact, uncompromised element that passes verification.

### Client-Side Verification (Black-Box Testing)
Security teams and compliance auditors can verify the firewall's active defenses without requiring direct access to the underlying source code. To validate the intercept engine, send a standard JSON-RPC request containing a known malicious signature (such as the AWS Metadata SSRF IP address) to the exposed proxy endpoint.

**Example Verification Command (Linux/macOS):**
```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{
        "jsonrpc": "2.0", 
        "id": "audit-test-001", 
        "method": "chat", 
        "params": {
            "name": "agent", 
            "arguments": {
                "target": "http://169.254.169.254/latest/meta-data/"
            }
        }
      }'
```

**Expected System Response:**
The firewall will instantly identify the Layer 4 SSRF target signature within the parameters structure, sever the transaction, and return a sanitized JSON-RPC error payload (`-32001: Security Policy Violation`). The downstream backend service will never receive the transaction.

---

### Guidance for Client Automated Testing (E2E / API Suites)

If your organization utilizes automated API testing frameworks (e.g., Postman, Playwright, Cypress, k6) to validate upstream routing through Panovista, tests executing static payloads sequentially will fail. 

The proxy's **30-second Cryptographic Replay Window** is absolute; it does not contain development backdoors, and the tracking cache cannot be bypassed via environment parameters. To successfully execute automated test suites without triggering false-positive Replay Attack blocks (`HTTP 409 Conflict`), your QA engineering teams must configure their test runners to dynamically inject a unique UUID into the JSON-RPC `id` field for every discrete request fixture.

#### Example 1: k6 Volume & Load Testing Implementation
```javascript
import http from 'k6/http';
import { check } from 'k6';
import { uuidv4 } from '[https://jslib.k6.io/k6-utils/1.4.0/index.js](https://jslib.k6.io/k6-utils/1.4.0/index.js)';

export default function () {
    const payload = JSON.stringify({
        jsonrpc: "2.0",
        id: uuidv4(), // Dynamic injection for strict replay cache parity compliance
        method: "chat",
        params: { /* target evaluation payload */ }
    });

    const res = http.post('http://localhost:8080/', payload, {
        headers: { 'Content-Type': 'application/json' },
    });
}
```

#### Example 2: Playwright E2E Functional Testing Implementation
```javascript
const response = await request.post('http://localhost:8080', {
  data: {
    jsonrpc: "2.0",
    id: crypto.randomUUID(), // Native cryptographic runtime module for dynamic IDs
    method: "database_query",
    params: { /* target evaluation payload */ }
  }
});

```