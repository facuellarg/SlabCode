package task

import (
	"slabcode/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//SetRoutes initialize routes for router service
func SetRoutes(echoServer *echo.Echo) {

	projectPetitions := echoServer.Group("/task")
	projectPetitions.Use(middleware.JWTWithConfig(
		config.GetJWTConfig()))
	projectPetitions.POST("/", CreateTask)
	projectPetitions.PATCH("/", FinalizeTask)
	projectPetitions.PUT("/", UpdateTask)

}
