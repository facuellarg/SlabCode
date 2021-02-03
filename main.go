package main

import (
	"log"
	"slabcode/project"
	"slabcode/task"
	"slabcode/user"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	echoServer := echo.New()
	echoServer.Use(
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
		}),
	)

	user.SetRoutes(echoServer)
	project.SetRoutes(echoServer)
	task.SetRoutes(echoServer)
	if err := echoServer.Start(":8000"); err != nil {
		log.Fatal(err.Error())
	}

}
