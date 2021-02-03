package project

import (
	"fmt"
	"slabcode/models"
	"slabcode/response"
	"slabcode/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

//CreateProject comment
func CreateProject(c echo.Context) error {
	var (
		project models.Project
		claims  = utils.GetClaims(c)
		message = echo.Map{
			"message": "Project created successfully.",
		}
	)
	err := c.Bind(&project)
	if err != nil {
		return c.JSON(response.SomthingWrongResponse(err))
	}
	status, err := CreateProjectService(project, claims)
	if err != nil {
		message["message"] = err.Error()
	}
	return c.JSON(status, message)
}

func UpdateProject(c echo.Context) error {
	var (
		project models.Project
		claims  = utils.GetClaims(c)
		message = echo.Map{
			"message": "Project was updated successfully.",
		}
	)
	err := c.Bind(&project)
	if err != nil {
		return c.JSON(response.SomthingWrongResponse(err))
	}
	status, err := UpdateProjectService(project, claims)
	if err != nil {
		message["message"] = err.Error()
	}
	return c.JSON(status, message)
}

func DeleteProject(c echo.Context) error {
	var (
		projectId int
		claims    = utils.GetClaims(c)
		message   = echo.Map{
			"message": "Project was deleted successfully.",
		}
	)
	projectId, err := strconv.Atoi(c.Param("projectId"))

	if err != nil {
		return c.JSON(response.SomthingWrongResponse(err))
	}
	fmt.Print(projectId)
	status, err := DeleteProjectService(projectId, claims)
	if err != nil {
		message["message"] = err.Error()
	}
	return c.JSON(status, message)
}
func FinalizeProject(c echo.Context) error {
	var (
		projectId int
		claims    = utils.GetClaims(c)
		message   = echo.Map{
			"message": "Project was finalized successfully.",
		}
	)
	projectId, err := strconv.Atoi(c.Param("projectId"))
	if err != nil {
		return c.JSON(response.SomthingWrongResponse(err))
	}
	fmt.Print(projectId)
	status, err := FinalizeProjectService(projectId, claims)
	if err != nil {
		message["message"] = err.Error()
	}
	return c.JSON(status, message)
}
