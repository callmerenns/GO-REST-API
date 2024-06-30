package model

import "github.com/golang-jwt/jwt/v5"

type MyCustomClaims struct {
	jwt.RegisteredClaims
	UserId uint   `json:"userId"`
	Role   string `json:"role"`
}
