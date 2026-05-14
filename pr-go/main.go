package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"

	config "project-go/comfig"
	database "project-go/db"
	"project-go/handlers"
	"project-go/middleware"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Подключаемся к SQLite (файл создастся автоматически)
	if err := database.Connect(cfg.DBPath); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	app := fiber.New()

	app.Use(logger.New())
	app.Use(recover.New())

	authHandler := handlers.NewAuthHandler(cfg.JWTSecret)
	linkHandler := handlers.NewLinkHandler()
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)

	api := app.Group("/api")
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	links := api.Group("/links")
	links.Use(authMiddleware)
	links.Post("/", linkHandler.CreateLink)
	links.Get("/:id", linkHandler.GetLink)
	links.Delete("/:id", linkHandler.DeleteLink)
	links.Get("/:id/stats", linkHandler.GetStats)

	app.Get("/:code", linkHandler.Redirect)

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
