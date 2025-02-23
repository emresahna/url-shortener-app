package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/emresahna/url-shortener-app/configs"
	"github.com/emresahna/url-shortener-app/internal/models"
	"github.com/emresahna/url-shortener-app/internal/sqlc"
	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenExpireMin  = 3
	refreshTokenExpireMin = 15
)

type Auth interface {
	Create(u sqlc.User) (models.LoginUserResponse, error)
	Parse(token string) (jwt.MapClaims, error)
}

type auth struct {
	pv *ecdsa.PrivateKey
	pb *ecdsa.PublicKey
}

func (a *auth) Create(u sqlc.User) (models.LoginUserResponse, error) {
	p := &models.Payload{
		Id: u.ID.String(),
	}

	at, err := a.createAccessToken(p)
	if err != nil {
		return models.LoginUserResponse{}, err
	}

	rt, err := a.createRefreshToken(p)
	if err != nil {
		return models.LoginUserResponse{}, err
	}

	return models.LoginUserResponse{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

func (a *auth) createAccessToken(p *models.Payload) (string, error) {
	c := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"id":  p.Id,
		"exp": time.Now().Add(time.Minute * accessTokenExpireMin).Unix(),
	})

	at, err := c.SignedString(a.pv)
	if err != nil {
		return "", err
	}

	return at, nil
}

func (a *auth) createRefreshToken(p *models.Payload) (string, error) {
	c := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"id":  p.Id,
		"exp": time.Now().Add(time.Minute * refreshTokenExpireMin).Unix(),
	})

	at, err := c.SignedString(a.pv)
	if err != nil {
		return "", err
	}

	return at, nil
}

func (a *auth) Parse(token string) (jwt.MapClaims, error) {
	c, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, nil
		}
		return a.pb, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := c.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, nil
}

func NewJWT(cfg configs.Auth) (Auth, error) {
	// Read and parse the private key
	pv, err := parsePrivateKey(cfg.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Read and parse the public key
	pb, err := parsePublicKey(cfg.PublicKeyPath)
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
