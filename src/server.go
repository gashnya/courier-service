package main

import (
	"github.com/labstack/echo/v4"
	"log"
	"yandex-team.ru/bstask/controller"
	"yandex-team.ru/bstask/database"
	"yandex-team.ru/bstask/models"
	"yandex-team.ru/bstask/repository"
	"yandex-team.ru/bstask/service"
)

func main() {
	db, err := database.SetupDatabase()
	if err != nil {
		log.Fatal(err)
	}

	e := setupServer(&controller.Controller{Service: service.Service{Repository: repository.Repository{DB: db}}})

	e.Logger.Fatal(e.Start(":8080"))
}

func setupServer(c *controller.Controller) *echo.Echo {
	e := echo.New()

	if err := models.SetupValidators(e); err != nil {
		log.Fatal(err)
	}

	c.SetupRoutes(e)

	return e
}
