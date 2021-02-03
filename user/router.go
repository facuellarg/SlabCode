package user

import (
	"slabcode/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//SetRoutes initialize routes for users service
func SetRoutes(echoServer *echo.Echo) {
	echoServer.POST("/singIn", SingIn)

	usersPetitions := echoServer.Group("/user")
	usersPetitions.Use(middleware.JWTWithConfig(
		config.GetJWTConfig()))
	usersPetitions.POST("/singUp", SingUp)
	usersPetitions.GET("", GetUsers)
	usersPetitions.GET("/welcome", Welcome)
	usersPetitions.GET("/ban", BanUser)
	usersPetitions.PATCH("/updatePassword", UpdatePassword)
}
