package main

import (
	"fmt"
	"github.com/ashah360/fibertools"
	"github.com/ashah360/foodquest-api/internal/api/db"
	"github.com/ashah360/foodquest-api/internal/api/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var port = "4000"

func main() {
	creds := db.DBInfo{
		Name:     os.Getenv("SQL_DBNAME"),
		Host:     os.Getenv("SQL_HOSTNAME"),
		Port:     os.Getenv("SQL_PORT"),
		User:     os.Getenv("SQL_USER"),
		Password: os.Getenv("SQL_PASSWORD"),
	}

	conn, err := creds.Connect("postgres")
	if err != nil {
		log.Fatal(err)
	}

	// userService := service.NewUserService(conn)

	server := fiber.New(fiber.Config{
		ErrorHandler: fibertools.ErrorHandler,
	})

	h := handlers.NewHandlerGroup(conn)

	server.Use(fibertools.Recover())
	server.Use(cors.New())

	server.Post("/register", h.Register)
	server.Post("/login", h.Login)
	server.Get("/users/me", h.GetCurrentUser)

	if err := server.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Println("An error occured, shutting down gracefully. ", err)
		_ = server.Shutdown()
	}
}
