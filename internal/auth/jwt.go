package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type Auth interface {
}

type auth struct {
	pv *ecdsa.PrivateKey
	pb *ecdsa.PublicKey
}

func NewJWTAuth() (Auth, error) {
	// Read and parse the private key
	pv, err := parsePrivateKey("./configs/ssl/private.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Read and parse the public key
	pb, err := parsePublicKey("./configs/ssl/public.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return &auth{
		pv: pv,
		pb: pb,
	}, nil
}

func parsePrivateKey(filePath string) (*ecdsa.PrivateKey, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block from %s", filePath)
	}

	return x509.ParseECPrivateKey(block.Bytes)
}

func parsePublicKey(filePath string) (*ecdsa.PublicKey, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block from %s", filePath)
	}

	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pubKey, ok := parsedKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("parsed key is not of type *ecdsa.PublicKey")
	}

	return pubKey, nil
}
