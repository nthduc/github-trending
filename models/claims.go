package models

import "github.com/dgrijalva/jwt-go"

type JwtCustomClaims struct {
	UserId string
	Role   string
	jwt.StandardClaims
}