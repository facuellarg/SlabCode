package models

import "time"

//Rol Model
type Baned struct {
	ID          uint `gorm:"primaryKey;autoIncrement" faker:"-"`
	UserBanedId uint `gorm:"unique"`
	WhoBanedId  uint
	Date        time.Time `faker:"-"`
}
