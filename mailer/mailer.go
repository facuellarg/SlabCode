package mailer

import (
	"crypto/tls"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

type mailerConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type MailMessage struct {
	To      []string
	Subject string
	Body    string
}

var (
	MailerConfig *mailerConfig
	once         sync.Once
	mailer       *gomail.Dialer
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No se pudieron cargar las variables de entorno")
	}
	MailerConfig = &mailerConfig{
		Host:     "smtp.gmail.com",
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
		Port:     587,
	}
}

//GetMailer return reference to mailer
//for use declare USERNAME and PASSWORD enviroment variables
func GetMailer() *gomail.Dialer {
	once.Do(func() {
		mailer = gomail.NewDialer(
			MailerConfig.Host,
			MailerConfig.Port,
			MailerConfig.Username,
			MailerConfig.Password,
		)
		mailer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	})
	return mailer

}
func SendMail(message MailMessage) error {
	mailer := GetMailer()
	gomailMessage := gomail.NewMessage()
	gomailMessage.SetHeader("From", MailerConfig.Username)
	gomailMessage.SetHeader("To", message.To...)
	gomailMessage.SetHeader("Subject", message.Subject)
	gomailMessage.SetBody("text/html", message.Body)

	if err := mailer.DialAndSend(gomailMessage); err != nil {
		return err
	}
	return nil
}
