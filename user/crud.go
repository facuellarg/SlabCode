package user

import (
	"net/http"
	"slabcode/database"
	"slabcode/models"
	"slabcode/response"
	"slabcode/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

//GetUsers get all users
func GetUsers(c echo.Context) error {
	var users []models.User
	connection := database.GetDefaultConnection()
	connection.Connection.Find(&users)
	return c.JSON(http.StatusOK, users)
}

//SingUp insert user in database
func SingUp(c echo.Context) error {
	user := &SingUpDto{}
	claims := utils.GetClaims(c)
	if err := c.Bind(user); err != nil {
		return c.JSON(response.SomthingWrongResponse(err))
	}
	status, err := SingUpService(c, claims, user)
	if err != nil {
		return c.JSON(status, echo.Map{"message": err.Error()})
	}
	return c.JSON(http.StatusCreated, user)
}

//SingIn login user
func SingIn(c echo.Context) error {
	singInData := &SingInDto{}
	if err := c.Bind(singInData); err != nil {
		return c.JSON(response.SomthingWrongResponse(err))
	}
	var token string
	status, err := SingInServive(c, singInData, &token)
	if err != "" {
		return c.JSON(status, echo.Map{"message": err})
	}
	return c.JSON(status, echo.Map{
		"token": token,
	})
}

//Welcome say welcome
func Welcome(c echo.Context) error {
	claims := utils.GetClaims(c)
	return c.String(http.StatusOK, "Welcome "+claims.Name+"!")
}

//BanUser ban user only if user that make petition is admin
func BanUser(c echo.Context) error {
	claims := utils.GetClaims(c)
	userBanedId, err := strconv.Atoi(c.QueryParam("userBanedId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Id must be a valid number",
		})
	}
	status, err := BanUserService(claims, userBanedId)
	if err != nil {
		return c.JSON(status, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(status, echo.Map{
		"message": "User Baned successfully.",
	})
}

func UpdatePassword(c echo.Context) error {
	claims := utils.GetClaims(c)
	updatePasswordDto := UpdatePasswordDto{}
	err := c.Bind(&updatePasswordDto)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{
			"message": err.Error(),
		})
	}
	status, err := UpdatePasswordService(claims, updatePasswordDto)
	if err != nil {
		return c.JSON(status, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(status, echo.Map{
		"message": "Password updated successfully.",
	})
}
