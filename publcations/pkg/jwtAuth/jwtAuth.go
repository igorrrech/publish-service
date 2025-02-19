package jwtauth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	User
	jwt.RegisteredClaims
}
type User struct {
	UUID uint   `json:"uuid"`
	Role string `json:"rol"`
}

var (
	UserRole    = "user"
	ManagerRole = "manager"
	AdminRole   = "admin"
)

func AccessValidate(token string, secret string) (*User, error) {
	t, err := jwt.ParseWithClaims(
		token,
		&AccessClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, ErrTokenIsNotValid
	}
	claims := t.Claims.(*AccessClaims)
	return &User{
		UUID: claims.UUID,
		Role: claims.Role,
	}, nil
}

var (
	ErrTokenIsNotValid = errors.New("token is not valid")
)
