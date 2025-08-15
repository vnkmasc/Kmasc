package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"os"
)

func SignJSONWithPrivateKey(data map[string]interface{}, privateKeyPath string) (string, error) {
	// Canonicalize JSON
	canonical, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// Read private key from file
	keyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return "", errors.New("failed to decode PEM private key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// Hash the data
	hashed := sha256.Sum256(canonical)

	// Sign it
	signature, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}

	return EncodeBase64(signature), nil
}

func EncodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
