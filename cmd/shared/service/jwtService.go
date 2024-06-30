package service

import (
	"fmt"
	"time"

	"github.com/altsaqif/go-rest/cmd/config"
	"github.com/altsaqif/go-rest/cmd/entity/dto"
	"github.com/altsaqif/go-rest/cmd/shared/model"
	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	CreateToken(author dto.UserWithProducts) (dto.AuthResponseDto, error)
	ParseToken(tokenHeader string) (jwt.MapClaims, error)
}

type jwtService struct {
	cfg config.TokenConfig
}

func (j *jwtService) CreateToken(user dto.UserWithProducts) (dto.AuthResponseDto, error) {
	claims := model.MyCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.cfg.IssuerName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.JwtExpiresTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserId: user.ID,
		Role:   user.Role,
	}

	token := jwt.NewWithClaims(j.cfg.JwtSigningMethod, claims)
	ss, err := token.SignedString(j.cfg.JwtSignatureKey)
	if err != nil {
		return dto.AuthResponseDto{}, fmt.Errorf("oops, failed to create token: %v", err)
	}
	return dto.AuthResponseDto{Token: ss}, nil
}

func (j *jwtService) ParseToken(tokenHeader string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		// if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 	return nil, fmt.Errorf("oops, unexpected signing method: %v", token.Header["alg"])
		// }
		return j.cfg.JwtSignatureKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("oops, failed to verify token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("oops, failed to parse token claims")
	}
	return claims, nil
}

func NewJwtService(cfg config.TokenConfig) JwtService {
	return &jwtService{cfg: cfg}
}
