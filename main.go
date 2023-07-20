package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/jclumbiarres/gofibertmpl/config"
	"github.com/jclumbiarres/gofibertmpl/controllers"
	"github.com/jclumbiarres/gofibertmpl/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	iniciaBD()
	jwt := middleware.NewAuthMiddleware(config.Secret)
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/assets", "./assets")

	app.Get("/protected", jwt, controllers.Protected)

	app.Post("/login", controllers.Login)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	log.Fatal(app.Listen(":3000"))
}

func iniciaBD() {
	var err error
	controllers.DB, err = gorm.Open(postgres.Open(config.PostgreCONNTXT))
	if err != nil {
		panic("Error al conectar a la base de datos.")
	}
	log.Println("Conectado a la base de datos.")
}
