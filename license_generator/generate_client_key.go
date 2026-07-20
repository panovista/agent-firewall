package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// LicenseClaims maps precisely to the validated V3.0 schema checked at container initialization.
type LicenseClaims struct {
	CustomerID          string `json:"customer_id"`
	LicenseTier         string `json:"license_tier"`
	IssuedAt            int64  `json:"issued_at"`
	ExpiresAt           int64  `json:"expires_at"`
	TargetEnvironmentID string `json:"target_environment_id"`
}

func main() {
	// -------------------------------------------------------------------------
	// YOUR HARDCODED MASTER SYSTEM CREDENTIALS
	// -------------------------------------------------------------------------
	// REPLACE THE STRING BELOW WITH YOUR 128-CHARACTER PRIVATE KEY
	const rawMasterPrivateKeyHex = "66b9688892374126bfbf6c7e9ebb87527334d65df4738a669253ac2b6f0d01e3ebb30ad59645d4bcbc1c7fcdc976db676fef024db5c0713b552cfac72f1d4d55"
	
	// -------------------------------------------------------------------------
	// CUSTOMER CONTRACT PARAMETERS (Edit these when you close a new deal)
	// -------------------------------------------------------------------------
	targetCustomerID := "enterprise-client-id-xyz"
	targetVpcNamespace := "aws-vpc-id-or-k8s-namespace" // Enforces structural scoping
	durationMonths := 12                                 // Standard 1-year procurement block
	
	// Parse the master root key directly from the local vault configuration
	privateKeyBytes, err := hex.DecodeString(rawMasterPrivateKeyHex)
	if err != nil {
		log.Fatalf("CRYPTOGRAPHIC SECURITY FAILURE: Hex decoding failed. Check your key formatting.")
	}
	
	// Failsafe to ensure you pasted the 128-character (64-byte) key, not the 64-character (32-byte) one.
	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		log.Fatalf("CRYPTOGRAPHIC SECURITY FAILURE: Master Private Key is malformed. Expected 64 bytes, got %d bytes.", len(privateKeyBytes))
	}
	privKey := ed25519.PrivateKey(privateKeyBytes)

	// Calculate precise epoch timestamps for boundary compliance loops
	issueTime := time.Now().Unix()
	expirationTime := time.Now().AddDate(0, durationMonths, 0).Unix()

	// Build the exact V3.0 verified payload token structure
	claims := LicenseClaims{
		CustomerID:          targetCustomerID,
		LicenseTier:         "sovereign_enterprise", // Upgrades them past Path A trial locks
		IssuedAt:            issueTime,
		ExpiresAt:           expirationTime,
		TargetEnvironmentID: targetVpcNamespace,
	}

	// Marshal to strict JSON data bytes
	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		log.Fatalf("SYSTEM ERROR: Failed to marshal client license claims.")
	}

	// Encode payload strictly utilizing Base64 Raw URL Encoding without trailing padding
	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadBytes)

	// Generate asymmetric cryptographic signature over the base64 payload bytes
	signatureBytes := ed25519.Sign(privKey, []byte(encodedPayload))
	hexSignature := hex.EncodeToString(signatureBytes)

	// Form the definitive shippable production string split configuration
	finalTokenString := fmt.Sprintf("pv_lic_%s.%s", encodedPayload, hexSignature)

	// Print clean verification results
	fmt.Println("\n=======================================================")
	fmt.Println("         PANOVISTA LICENSING VAULT MINT COMPLETED      ")
	fmt.Println("=======================================================")
	fmt.Printf("CLIENT IDENTITY : %s\n", targetCustomerID)
	fmt.Printf("NETWORK ENCLAVE : %s\n", targetVpcNamespace)
	fmt.Printf("TIER ESCALATION : sovereign_enterprise\n")
	fmt.Printf("EXPIRATION DATE : %s\n", time.Unix(expirationTime, 0).UTC().Format(time.RFC1123))
	fmt.Println("-------------------------------------------------------")
	fmt.Println("PASTE THIS BLOCK IN CLIENT ENVIRONMENT VARIABLE:")
	fmt.Printf("PANOVISTA_LICENSE=%s\n", finalTokenString)
	fmt.Println("=======================================================\n")
}