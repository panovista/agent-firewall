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

// LicenseClaims matches the Panovista V2.0 Core strict initialization schema
type LicenseClaims struct {
	CustomerID          string `json:"customer_id"`
	LicenseTier         string `json:"license_tier"`
	IssuedAt            int64  `json:"issued_at"`
	ExpiresAt           int64  `json:"expires_at"`
	TargetEnvironmentID string `json:"target_environment_id"`
}

func main() {
	// 1. Generate the Master Key Pair 
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatalf("Fatal: Cryptographic generation failed: %v", err)
	}

	// 2. Calculate the exact 14-Day evaluation window based on your system clock
	issueTime := time.Now().Unix()
	expirationTime := issueTime + (14 * 24 * 60 * 60) // +14 Days

	// 3. Construct the Claims
	claims := LicenseClaims{
		CustomerID:          "public-eval-001",
		LicenseTier:         "standard_vpc",
		IssuedAt:            issueTime,
		ExpiresAt:           expirationTime,
		TargetEnvironmentID: "vpc-eval",
	}

	// 4. Serialize and Encode Payload (Base64URL without padding)
	payloadBytes, _ := json.Marshal(claims)
	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadBytes)

	// 5. Sign Payload with Ed25519 Private Key and Hex Encode
	signature := ed25519.Sign(privKey, payloadBytes)
	encodedSignature := hex.EncodeToString(signature)

	// 6. Assemble the Final Token String
	finalToken := fmt.Sprintf("pv_lic_%s.%s", encodedPayload, encodedSignature)

	// Output the immutable artifacts
	fmt.Println("--- PANOVISTA CA VAULT: ARTIFACT EXPORT ---")
	fmt.Printf("PUBLIC_VERIFICATION_KEY (Hardcode in main.go) : %s\n", hex.EncodeToString(pubKey))
	fmt.Printf("14_DAY_EVAL_TOKEN (Paste in docker-compose)   : %s\n", finalToken)
	fmt.Println("-------------------------------------------")
}