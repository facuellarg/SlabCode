package config

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4/middleware"
)

//PayloadJWT payload for jwt token
type PayloadJWT struct {
	ID    uint
	Name  string `json:"name"`
	Email string `json:"email"`
	RolID uint   `json:"rol"`
	jwt.StandardClaims
}

//GetJWTSecret return secret for jwt
func GetJWTSecret() string {
	return secret
}

//GetJWTConfig get configuration for jwt
func GetJWTConfig() middleware.JWTConfig {
	return middleware.JWTConfig{
		Claims:     &PayloadJWT{},
		SigningKey: []byte(secret),
	}
}
