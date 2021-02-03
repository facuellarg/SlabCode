package user

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"sync"
	"text/template"
	"time"

	"slabcode/config"
	"slabcode/database"
	"slabcode/mailer"
	"slabcode/models"
	"slabcode/response"
	"slabcode/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//SingUpService singup user
func SingUpService(c echo.Context, claims *config.PayloadJWT, userDto *SingUpDto) (int, error) {

	var (
		buf           bytes.Buffer
		wg            sync.WaitGroup
		txStatus      int
		plainPassword string
	)
	db := database.GetDefaultConnection()
	if !utils.IsRol(db, claims, models.Administrator) {
		return http.StatusForbidden, errors.New("This route is only for administrator user.")
	}
	user := &userDto.User
	plainPassword = user.Password
	hash, err := bcrypt.
		GenerateFromPassword(
			[]byte((*user).Password),
			bcrypt.MinCost)

	if err == nil {
		user.Password = string(hash)
	}
	rol := models.Rol{}

	err = db.Connection.Where("id = ?", userDto.RolID).First(&rol).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusBadRequest, errors.New("Rol id  id invalid.")

	}
	user.RolID = rol.ID

	err = db.Connection.Transaction(func(tx *gorm.DB) error {
		result := tx.
			Clauses(clause.OnConflict{DoNothing: true}).
			Create(user)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			txStatus = http.StatusConflict
			return fmt.Errorf(response.ALREADY_EXISTS, "User")
		}
		tmpl := template.Must(template.ParseFiles(path.Join(exPath, "user/template.html")))
		err = tmpl.Execute(&buf, map[string]string{
			"username": user.Email,
			"password": plainPassword,
		})
		if err != nil {
			txStatus = http.StatusBadGateway
			return err
		}

		message := mailer.MailMessage{
			To:      []string{user.Email},
			Body:    buf.String(),
			Subject: "Bienvenido",
		}
		wg.Add(1)
		go func() {
			err := mailer.SendMail(message)
			if err != nil {
				tx.Rollback()
			}
			wg.Done()
		}()
		if err != nil {
			txStatus = http.StatusBadGateway
			return err
		}
		return nil
	})
	if err != nil {
		return txStatus, err
	}

	wg.Wait()
	return http.StatusCreated, nil
}

//SingInServive singin user
func SingInServive(c echo.Context, singinData *SingInDto, token *string) (int, string) {
	connection := database.GetDefaultConnection()

	user := &models.User{}
	err := connection.Connection.Where("email = ?", singinData.Email).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, "User not Found"
	}
	baned := &models.Baned{UserBanedId: user.ID + 1}
	err = connection.Connection.Where("user_baned_id = ?", user.ID).First(baned).Error
	if baned.UserBanedId == user.ID {
		return http.StatusForbidden, "User is baned"
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusBadGateway, err.Error()
	}
	if !validatePassword(singinData.Password, user) {
		return http.StatusUnauthorized, response.INVALID_CREDENTIALS
	}
	*token, err = generateJWTToken(user)
	if err != nil {
		return http.StatusBadGateway, err.Error()
	}

	return http.StatusOK, ""

}

func BanUserService(claims *config.PayloadJWT, userBanedId int) (int, error) {
	db := database.GetDefaultConnection()
	if !utils.IsRol(db, claims, models.Administrator) {
		return http.StatusForbidden, errors.New(response.FORBIDDEN_RESOURCE)
	}
	var (
		user  models.User
		baned models.Baned
	)

	err := db.Connection.Where("id = ?", userBanedId).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, errors.New("User not found")
	}
	baned.Date = time.Now()
	baned.UserBanedId = user.ID
	baned.WhoBanedId = claims.ID
	result := db.Connection.Clauses(clause.OnConflict{DoNothing: true}).Create(&baned)
	if result.Error != nil {
		return http.StatusBadGateway, result.Error
	}
	return http.StatusOK, nil
}

func UpdatePasswordService(claims *config.PayloadJWT, updatePasswordDto UpdatePasswordDto) (int, error) {
	hash, err := bcrypt.
		GenerateFromPassword(
			[]byte(updatePasswordDto.NewPassword),
			bcrypt.MinCost,
		)
	var user models.User
	db := database.GetDefaultConnection()
	err = db.Connection.Where("id = ?", claims.ID).First(&user).Error
	if err != nil {
		return http.StatusBadGateway, err
	}

	if !validatePassword(updatePasswordDto.OldPassword, &user) {
		return http.StatusForbidden, errors.New("Bad credentials")
	}
	user.Password = string(hash)
	err = db.Connection.Save(user).Error
	if err != nil {
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil

}

func generateJWTToken(user *models.User) (string, error) {
	claims := &config.PayloadJWT{
		ID:    user.ID,
		Name:  user.UserName,
		Email: user.Email,
		RolID: user.RolID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 5).Unix(),
		},
	}
	tokenGenerator := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenGenerator.SignedString([]byte(config.GetJWTSecret()))
	return token, err
}

func validatePassword(password string, user *models.User) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte((*user).Password),
		[]byte(password))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
