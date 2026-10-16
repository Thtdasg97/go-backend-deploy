package dto

import "github.com/golang-jwt/jwt/v5"

// define struct
type CustomClaim struct {
	UserId int
	jwt.RegisteredClaims
}
