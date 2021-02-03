package utils

import (
	"slabcode/config"
	"slabcode/database"
	"slabcode/models"
)

//IsRol return true is rol is the corresponding
func IsRol(db *database.DatabaseConnection, claims *config.PayloadJWT, contrast *models.Rol) bool {
	var rol models.Rol
	db.Connection.Where("id", claims.RolID).First(&rol)
	return rol.Rol == contrast.Rol

}
