package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slabcode/config"
	"slabcode/database"
	"slabcode/mailer"
	"slabcode/models"
	"slabcode/response"
	"slabcode/utils"
	"time"

	"gorm.io/gorm"
)

func CreateProjectService(project models.Project, claims *config.PayloadJWT) (int, error) {

	db := database.GetDefaultConnection()
	if !utils.IsRol(db, claims, models.Operator) {
		return http.StatusForbidden, errors.New(response.FORBIDDEN_RESOURCE)
	}
	isValid, err := isValidProject(project)
	if !isValid {
		return http.StatusBadRequest, err
	}
	project.UserId = claims.ID
	project.State = models.PROJECT_IN_PROGRESS
	result := db.Connection.Create(&project)
	if result.Error != nil {
		return http.StatusBadGateway, result.Error
	}
	return http.StatusCreated, nil
}
func FinalizeProjectService(projectId int, claims *config.PayloadJWT) (int, error) {
	db := database.GetDefaultConnection()
	var (
		currentProject   models.Project
		projectTasksLen  int64
		rolAdministrator models.Rol
	)
	if !utils.IsRol(db, claims, models.Operator) {
		return http.StatusForbidden, errors.New(response.FORBIDDEN_RESOURCE)
	}

	err := db.Connection.Where("id = ? and user_id = ?", projectId, claims.ID).First(&currentProject).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, errors.New("Project not found")
	}
	db.Connection.Model(&models.Task{}).
		Where("project_id = ? and state = ?", projectId, models.TASK_IN_PROGRESS).
		Count(&projectTasksLen)
	if projectTasksLen > 0 {
		return http.StatusOK, errors.New("The project have some task whitout finalize")
	}
	currentProject.State = models.PROJECT_FINALIZE
	body, err := json.MarshalIndent(&currentProject, "", "    ")
	if err != nil {
		return http.StatusBadGateway, err
	}

	db.Connection.Where("rol = ?", models.Administrator.Rol).First(&rolAdministrator)
	administrators := make([]models.User, 0)
	db.Connection.Where("rol_id = ?", rolAdministrator.ID).Find(&administrators)
	emailsAdministrators := make([]string, len(administrators))
	for i, admin := range administrators {
		emailsAdministrators[i] = admin.Email
	}
	fmt.Printf("%v\n", emailsAdministrators)
	message := mailer.MailMessage{
		To:      emailsAdministrators,
		Subject: "finalized project",
		Body:    string(body),
	}
	err = mailer.SendMail(message)
	if err != nil {
		return http.StatusBadGateway, err
	}
	result := db.Connection.Save(&currentProject)
	if result.Error != nil {
		return http.StatusBadGateway, result.Error
	}
	return http.StatusOK, nil
}

func DeleteProjectService(projectId int, claims *config.PayloadJWT) (int, error) {
	db := database.GetDefaultConnection()
	var (
		currentProject models.Project
	)
	if !utils.IsRol(db, claims, models.Operator) {
		return http.StatusForbidden, errors.New(response.FORBIDDEN_RESOURCE)
	}
	err := db.Connection.Where("id = ? and user_id = ?", projectId, claims.ID).First(&currentProject).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, errors.New("Project not found")
	}
	result := db.Connection.Delete(&currentProject)
	if result.Error != nil {
		return http.StatusBadGateway, result.Error
	}
	return http.StatusCreated, nil
}

func UpdateProjectService(project models.Project, claims *config.PayloadJWT) (int, error) {
	db := database.GetDefaultConnection()
	currentProject := models.Project{}
	projectTasks := make([]models.Task, 0)
	if !utils.IsRol(db, claims, models.Operator) {
		return http.StatusForbidden, errors.New(response.FORBIDDEN_RESOURCE)
	}
	err := db.Connection.Where("id = ? and user_id = ?", project.ID, claims.ID).First(&currentProject).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, errors.New("Project not found")
	}
	isValid, err := isValidProject(project)
	if !isValid {
		return http.StatusBadRequest, err
	}
	if project.Name != "" {
		currentProject.Name = project.Name
	}
	if project.Description != "" {
		currentProject.Description = project.Description
	}
	db.Connection.
		Where("project_id = ? and state = ?", project.ID, models.TASK_IN_PROGRESS).
		Find(&projectTasks)
	if (project.EndDate != time.Time{}) {
		for _, projectTask := range projectTasks {
			if projectTask.EndDate.After(project.EndDate) {
				return http.StatusBadRequest, errors.New("Cant update date, some task have a lastest end date.")
			}
		}
		currentProject.EndDate = project.EndDate
	}
	result := db.Connection.Save(&currentProject)
	if result.Error != nil {
		return http.StatusBadGateway, result.Error
	}
	return http.StatusOK, nil
}

func isValidProject(project models.Project) (bool, error) {
	if project.Name == "" || project.Description == "" {
		return false, errors.New("The name and description should no be empty")
	}
	if project.StartDate.Before(time.Now()) || project.EndDate.Before(project.StartDate) {
		return false, errors.New("The end date should be greather than start and start should be greather than now")
	}
	return true, nil
}
