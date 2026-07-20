package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
)

func main() {
	pubKey, privKey, _ := ed25519.GenerateKey(nil)
	
	fmt.Println("=== YOUR PERMANENT MASTER KEYS ===")
	fmt.Printf("PUBLIC KEY (32 Bytes - Put in main.go):\n%s\n\n", hex.EncodeToString(pubKey))
	fmt.Printf("PRIVATE KEY (64 Bytes - Put in generate_client_key.go):\n%s\n", hex.EncodeToString(privKey))
	fmt.Println("==================================")
}