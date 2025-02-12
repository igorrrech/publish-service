package models

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type Token string
type AccesToken Token
type RefreshToken Token
type TokenPair struct {
	Access  AccesToken
	Refresh RefreshToken
}
type AccessClaims struct {
	UUID uint   `json:"uuid"`
	Role string `json:"rol"`
	jwt.RegisteredClaims
}
type RefreshClaims struct {
	UUID uint `json:"uuid"`
	jwt.RegisteredClaims
}

func NewTokenPair(uuid uint, ttl time.Duration, rol string, secret string) (p TokenPair, e error) {
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessClaims{
		uuid,
		rol,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	})
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshClaims{
		uuid,
		jwt.RegisteredClaims{},
	})
	a_signed, err := access.SignedString([]byte(secret))
	if err != nil {
		return p, err
	}
	r_signed, err := refresh.SignedString([]byte(secret))
	if err != nil {
		return p, err
	}
	p.Access = AccesToken(a_signed)
	p.Refresh = RefreshToken(r_signed)
	return p, nil
}
func (t AccesToken) VerifyToken(claims AccessClaims, secret string) (*jwt.Token, error) {
	// Parse the token with the secret key
	token, err := jwt.ParseWithClaims(string(t), &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	// Check if the token is valid
	if !token.Valid {
		return token, fmt.Errorf("invalid token")
	}

	// Check for verification errors
	if err != nil {
		return token, err
	}

	// Return the verified token
	return token, nil
}
func (t RefreshToken) VerifyToken(claims RefreshClaims, secret string) (*jwt.Token, error) {
	// Parse the token with the secret key
	token, err := jwt.ParseWithClaims(string(t), &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	// Check for verification errors
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Return the verified token
	return token, nil
}
