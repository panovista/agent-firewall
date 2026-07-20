package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// --- RED TEAM TEST 1 & 6: The Obfuscation Gauntlet & Homoglyphs ---
func TestAnalyzeObfuscation(t *testing.T) {
	// Test A: Nested Encoding (URL -> Hex -> Base64 -> Raw) targeting standard SSRF
	// "169.254.169.254" Base64'd, then Hex'd, then URL encoded
	nestedPayload := "%33%31%33%36%33%39%32%65%33%32%33%35%33%34%32%65%33%31%33%36%33%39%32%65%33%32%33%35%33%34" 
	
	_, err := analyzeObfuscation(nestedPayload)
	if err == nil || !strings.Contains(err.Error(), "nested evasion detected") {
		t.Fatalf("FAILED: Obfuscation Gauntlet did not catch nested SSRF payload.")
	}

	// Test B: Hex SSRF bypass attempt (Decimal/Hex routing)
	hexIP := "0xa9.0xfe.0xa9.0xfe"
	_, err2 := analyzeObfuscation(hexIP)
	if err2 == nil || !strings.Contains(err2.Error(), "ssrf metadata attack detected") {
		t.Fatalf("FAILED: Engine failed to catch Hex-encoded L4 routing subversion.")
	}
}

// --- RED TEAM TEST 2: Algorithmic Complexity (JSON DoS) ---
func TestJSONRecursionLimiter(t *testing.T) {
	// Build a heavily nested payload that exceeds MaxJSONDepth (20)
	var nestedMap interface{} = "malicious_payload"
	for i := 0; i < 25; i++ {
		nestedMap = map[string]interface{}{
			"layer": nestedMap,
		}
	}

	_, err := scrubJSON(nestedMap, 0)
	if err == nil || err.Error() != "recursion depth limit exceeded" {
		t.Fatalf("FAILED: JSON DoS trap failed. Payload penetrated deeper than 20 layers.")
	}
}

// --- RED TEAM TEST 7: IP-Based Rate Limiting ---
func TestRateLimiter(t *testing.T) {
	testIP := "192.168.1.100"
	
	// Reset the cache for a clean test
	rateLimitCache.Range(func(key, value interface{}) bool {
		rateLimitCache.Delete(key)
		return true
	})

	// Simulate 300 valid requests
	for i := 0; i < RateLimitMax; i++ {
		if isRateLimited(testIP) {
			t.Fatalf("FAILED: Rate limiter falsely blocked traffic before the threshold at request %d", i)
		}
	}

	// The 301st request MUST be blocked
	if !isRateLimited(testIP) {
		t.Fatalf("FAILED: Rate limiter failed to block request 301. Denial of Wallet vulnerability exposed.")
	}
}

// --- RED TEAM TEST 3: Cryptographic Replay Cache ---
func TestReplayAttack(t *testing.T) {
	reqID := "audit-test-uuid-001"
	
	// Use the exact same cryptographic logic as production
	hashBytes := sha256.Sum256([]byte(reqID))
	hash := hex.EncodeToString(hashBytes[:])

	// Manually inject the valid hash into the cache
	replayCache.Store(hash, "timestamp")

	// Attempt to push a payload with the same ID
	reqBody := `{"jsonrpc": "2.0", "id": "audit-test-uuid-001", "method": "chat", "params": {}}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
	req.RemoteAddr = "10.0.0.5:12345"
	w := httptest.NewRecorder()

	proxyHandler(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("FAILED: Replay attack penetrated. Expected HTTP 409 Conflict, got HTTP %d", w.Code)
	}
}