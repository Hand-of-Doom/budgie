package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/configor"
)

func runApplication() error {
	app := fiber.New()
	app.Static("/", "./public")

	config := struct{ Port string }{}
	err := configor.Load(&config, "./config.yaml")
	if err != nil {
		return err
	}

	return app.Listen(":" + config.Port)
}

func main() {
	if err := runApplication(); err != nil {
		panic(err)
	}
}
