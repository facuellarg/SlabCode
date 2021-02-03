package task

import (
	"net/http"
	"slabcode/models"
	"slabcode/response"

	"github.com/labstack/echo/v4"
)

//CreateTask create task with post request
func CreateTask(c echo.Context) error {
	var (
		task    models.Task
		message = echo.Map{
			"message": "Task created successfully.",
		}
	)
	err := c.Bind(&task)
	if err != nil {
		return c.JSON(response.SomthingWrongResponse(err))
	}
	status, err := CreateTaskService(c, task)
	if err != nil {
		message["message"] = err.Error()
	}
	return c.JSON(status, message)
}

func FinalizeTask(c echo.Context) error {
	var (
		message = echo.Map{
			"message": "Task finalized successfully.",
		}
		finalizeTaskDto FinalizeTaskDto
	)

	err := c.Bind(&finalizeTaskDto)
	if err != nil {
		return c.JSON(http.StatusBadGateway, err)
	}
	status, err := FinalizeTaskService(c, finalizeTaskDto)
	if err != nil {
		message["message"] = err.Error()
	}
	return c.JSON(status, message)
}

func UpdateTask(c echo.Context) error {
	var (
		message = echo.Map{
			"message": "Task updated successfully.",
		}
		task models.Task
	)

	err := c.Bind(&task)
	if err != nil {
		return c.JSON(http.StatusBadGateway, err)
	}
	status, err := UpdateTaskService(c, task)
	if err != nil {
		message["message"] = err.Error()
	}
	return c.JSON(status, message)
}
