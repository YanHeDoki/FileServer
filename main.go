package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()
	app.Use()

	app.Get("/", func(c *fiber.Ctx) error {

		return c.SendString("hello world")

	})

	app.Listen(":9512")

}
