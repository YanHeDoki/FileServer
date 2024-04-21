package main

import (
	"FileServerFiber/logic"

	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()
	//app.Use(recover.New())

	app.Get("/", logic.GetFileList)
	app.Get("/download/+", logic.DownloadFile)

	app.Listen(":9512")

}
