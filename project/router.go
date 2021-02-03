package project

import (
	"slabcode/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//SetRoutes initialize routes for project service
func SetRoutes(echoServer *echo.Echo) {

	projectPetitions := echoServer.Group("/project")
	projectPetitions.Use(middleware.JWTWithConfig(
		config.GetJWTConfig()))
	projectPetitions.POST("/", CreateProject)
	projectPetitions.PUT("/", UpdateProject)
	projectPetitions.DELETE("/:projectId", DeleteProject)
	projectPetitions.PATCH("/:projectId", FinalizeProject)

}
