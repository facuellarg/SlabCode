package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	secret     string
	debug      bool
	dbUserName string
	dbPassword string
	dbDatabase string
	dbHost     string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	secret = os.Getenv("JWT_SECRET")
	dbUserName = os.Getenv("DB_USERNAME")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbDatabase = os.Getenv("DB_DATABASE")
	dbHost = os.Getenv("DB_HOST")
	debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
}
