package utils

import (
	"slabcode/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

//GetClaims return a claims for request
func GetClaims(c echo.Context) *config.PayloadJWT {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*config.PayloadJWT)
	return claims
}
