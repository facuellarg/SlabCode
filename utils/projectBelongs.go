package utils

import (
	"errors"
	"slabcode/database"
	"slabcode/models"

	"gorm.io/gorm"
)

//ProjectBelogns return true if a project belogns to user
func ProjectBelongs(userId, projectId uint) (bool, models.Project) {
	var project models.Project
	db := database.GetDefaultConnection()
	err := db.Connection.Where("id = ? and user_id = ?", projectId, userId).First(&project).Error
	return !errors.Is(err, gorm.ErrRecordNotFound), project
}
