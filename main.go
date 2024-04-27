package main

import (
	"FileServerFiber/logic"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	app := fiber.New()
	app.Use(recover.New())

	app.Get("/", logic.GetFileList)
	app.Get("/download/+", logic.MyDownloadFile)

	app.Listen(":9512")

}
