package main

import (
	"FileServerFiber/logic"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

func main() {

	engine := html.New("./tmpl", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(recover.New())

	//app.Get("/*", logic.GetFileList)
	app.Get("/download/+", logic.MyDownloadFile)

	app.Get("/*", logic.Tmpl)

	app.Listen(":9512")

}
