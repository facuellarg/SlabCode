package models

import (
	"time"

	"gorm.io/gorm"
)

//User Model
type User struct {
	gorm.Model `faker:"-"`
	ID         uint      `gorm:"primaryKey;autoIncrement" faker:"-"`
	UserName   string    `faker:"username"`
	Email      string    `gorm:"unique" faker:"email"`
	Password   string    `faker:"password"`
	Birthday   time.Time `faker:"-"`
	RolID      uint
	WhoBaned   []Baned   `gorm:"foreignKey:WhoBanedId" json:"-"`
	Projects   []Project `gorm:"foreignKey:UserId" json:"-"`
	BanedId    Baned     `gorm:"foreignKey:UserBanedId" json:"-"`
	// Rol        Rol `gorm:"foreignKey:RolID" faker:"-"`
}
