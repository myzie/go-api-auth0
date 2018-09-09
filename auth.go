package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	certBegin = "-----BEGIN CERTIFICATE-----"
	certEnd   = "-----END CERTIFICATE-----"
)

// JSONWebKey is a single JSON Web Key definition
type JSONWebKey struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// GetCertificate returns the certificate string for this JSONWebKey
func (key JSONWebKey) GetCertificate() string {
	return fmt.Sprintf("%s\n%s\n%s", certBegin, key.X5c[0], certEnd)
}

// JSONWebKeys is a response from a JWK endpoint
type JSONWebKeys struct {
	Keys []JSONWebKey `json:"keys"`
}

// Find looks for a key with the given key ID
// token.Header["kid"]
func (keys JSONWebKeys) Find(keyID string) (JSONWebKey, bool) {
	for k := range keys.Keys {
		if keyID == keys.Keys[k].Kid {
			return keys.Keys[k], true
		}
	}
	return JSONWebKey{}, false
}

// GetKeys retrieves public keys from a remote JWKS endpoint
// "https://" + os.Getenv("AUTH0_DOMAIN") + "/.well-known/jwks.json"
func GetKeys(url string) (JSONWebKeys, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return JSONWebKeys{}, err
	}

	ctx, cancel := context.WithTimeout(req.Context(), 500*time.Millisecond)
	defer cancel()

	req = req.WithContext(ctx)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return JSONWebKeys{}, err
	}
	defer resp.Body.Close()

	var jwks JSONWebKeys
	if json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return JSONWebKeys{}, err
	}
	return jwks, nil
}

// CustomClaims include a scope field with the JWT claims
type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

// CheckScope confirms the specified scope is declared in the JWT
func CheckScope(scope string, tokenString string) (bool, error) {

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, nil)
	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return false, errors.New("CustomClaims type assertion failed")
	}

	result := strings.Split(claims.Scope, " ")
	for i := range result {
		if result[i] == scope {
			return true, nil
		}
	}
	return false, nil
}
