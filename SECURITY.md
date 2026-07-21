# Security Policy

At Panovista, the security and integrity of our Layer 7 proxy and your enterprise data are our highest priorities. We maintain a strict Zero-Trust architectural standard and take all security vulnerabilities seriously.

## Supported Versions

We currently provide active security updates, cryptographic patches, and vulnerability mitigations for the following product versions:

| Version | Supported | Notes |
| :--- | :--- | :--- |
| **V3.0.x** | Yes | Active main branch (Includes FPT and offline cryptographic licensing). |
| **V2.x** | No | Deprecated. Users must upgrade to V3.0 for active mitigation support. |
| **V1.x** | No | End of Life (EOL). |

## Reporting a Vulnerability

If you have discovered a potential security vulnerability in the Panovista Proxy, **please do not disclose it publicly in GitHub Issues.** 

Instead, we ask that you practice responsible disclosure by contacting our security team directly. This allows us to protect our enterprise customers by verifying and patching the vulnerability before it is made public.

**How to Report:**
*   Email your findings directly to our security engineering team at: **ian.ayliffe@panovistamarketing.com**
*   Use the subject line: `[VULNERABILITY] Panovista Proxy - <Brief Description>`

**What to Include:**
*   The exact version of the Panovista container you are testing (e.g., `v3.0.1`).
*   A description of the vulnerability and its potential impact.
*   Detailed steps to reproduce the issue (including any specific JSON-RPC payloads, curl commands, or proxy configurations).
*   Any relevant logs or HTTP response codes.

## Our Commitment

*   **Acknowledge:** We will acknowledge receipt of your vulnerability report within 48 hours.
*   **Investigate:** We will provide an estimated timeline for validation and mitigation within 5 business days.
*   **Resolve:** Once patched, we will issue a formal security release to the GitHub Container Registry (GHCR) and publish the mitigations in our release notes.