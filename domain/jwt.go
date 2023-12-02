package domain

import "github.com/golang-jwt/jwt/v5"

type JwtCustomClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}
