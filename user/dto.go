package user

import "slabcode/models"

//SingInDto data for user sing in
type SingInDto struct {
	Email    string
	Password string
}

//SingUpDto data for user sing in
type SingUpDto struct {
	User  models.User
	RolID uint
}

type UpdatePasswordDto struct {
	OldPassword string
	NewPassword string
}
