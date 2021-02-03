package models

import "time"

const (
	PROJECT_IN_PROGRESS = "en proceso"
	PROJECT_FINALIZE    = "finalizado"
)

type Project struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"not null"`
	Description string    `gorm:"not null"`
	StartDate   time.Time `gorm:"not null" json:"start_date"`
	EndDate     time.Time `gorm:"not null" json:"end_date"`
	UserId      uint
	State       string
	Tasks       []Task `gorm:"foreignKey:ProjectId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-" `
}
