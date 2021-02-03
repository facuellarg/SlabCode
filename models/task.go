package models

import "time"

const (
	TASK_IN_PROGRESS = "en proceso"
	TASK_FINALIZE    = "finalizada"
)

type Task struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	Name        string
	Description string
	EndDate     time.Time `gorm:"not null" json:"end_date"`
	ProjectId   uint
	State       string
}
