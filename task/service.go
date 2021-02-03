package task

import (
	"errors"
	"net/http"
	"slabcode/config"
	"slabcode/database"
	"slabcode/models"
	"slabcode/response"
	"slabcode/utils"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreateTaskService(c echo.Context, task models.Task) (int, error) {
	claims := utils.GetClaims(c)
	db := database.GetDefaultConnection()
	if !utils.IsRol(db, claims, models.Operator) {
		return http.StatusForbidden, errors.New(response.FORBIDDEN_RESOURCE)
	}
	isValid, err := isValidTask(task, claims)
	if !isValid {
		return http.StatusBadRequest, err
	}
	task.State = models.TASK_IN_PROGRESS
	err = db.Connection.Create(&task).Error
	if err != nil {
		return http.StatusBadGateway, err
	}
	return http.StatusCreated, nil

}

func FinalizeTaskService(c echo.Context, finalizeTaskDto FinalizeTaskDto) (int, error) {
	claims := utils.GetClaims(c)
	db := database.GetDefaultConnection()
	if !utils.IsRol(db, claims, models.Operator) {
		return http.StatusForbidden, errors.New(response.FORBIDDEN_RESOURCE)
	}
	belongs, _ := utils.ProjectBelongs(claims.ID, finalizeTaskDto.ProjectId)
	if !belongs {
		return http.StatusBadRequest, errors.New("The project don't belongs to the user")
	}
	err := db.Connection.
		Model(&models.Task{}).
		Where("id = ?", finalizeTaskDto.TaskId).
		Update("state", models.TASK_FINALIZE).Error
	if err != nil {
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil

}

func UpdateTaskService(c echo.Context, newTask models.Task) (int, error) {
	claims := utils.GetClaims(c)
	db := database.GetDefaultConnection()
	var currentTask models.Task

	if !utils.IsRol(db, claims, models.Operator) {
		return http.StatusForbidden, errors.New(response.FORBIDDEN_RESOURCE)
	}

	err := db.Connection.Where("id = ? and project_id = ?", newTask.ID, newTask.ProjectId).First(&currentTask).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, errors.New("Cant found the task.")
	}

	newTask.State = currentTask.State
	if newTask.Description == "" {
		newTask.Description = currentTask.Description
	}
	if newTask.Name == "" {
		newTask.Name = currentTask.Name
	}
	if (newTask.EndDate == time.Time{}) {
		newTask.EndDate = currentTask.EndDate
	}
	isValid, err := isValidTask(newTask, claims)
	if !isValid {
		return http.StatusBadRequest, err
	}
	err = db.Connection.Save(newTask).Error
	if err != nil {
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil

}

func isValidTask(task models.Task, claims *config.PayloadJWT) (bool, error) {
	if task.Name == "" {
		return false, errors.New("The name should no be empty")
	}
	belongs, project := utils.ProjectBelongs(claims.ID, task.ProjectId)
	if !belongs {
		return false, errors.New("The project don't belongs to the user")
	}
	if task.EndDate.Before(project.StartDate) || task.EndDate.After(project.EndDate) {
		return false, errors.New("The end date for task should be in the time range of project selected.")
	}
	return true, nil
}
