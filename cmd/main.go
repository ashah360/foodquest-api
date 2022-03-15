package main

import (
	"fmt"
	"github.com/ashah360/fibertools"
	"github.com/ashah360/foodquest-api/internal/api/db"
	"github.com/ashah360/foodquest-api/internal/api/handlers"
	"github.com/ashah360/foodquest-api/internal/api/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
)

var port = "3000"

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

	rs := service.NewRestaurantService(conn)
	us := service.NewUserService(conn)

	h := handlers.NewHandlerGroup(conn, rs, us)

	server.Use(fibertools.Recover())
	server.Use(cors.New(cors.Config{
		Next:         nil,
		AllowOrigins: "*",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
			fiber.MethodOptions,
		}, ","),
		AllowHeaders: strings.Join([]string{
			fiber.HeaderAccept,
			fiber.HeaderAuthorization,
			fiber.HeaderReferer,
			fiber.HeaderUserAgent,
			fiber.HeaderOrigin,
			fiber.HeaderAcceptLanguage,
			fiber.HeaderContentType,
		}, ","),
		AllowCredentials: true,
		ExposeHeaders:    "",
		MaxAge:           0,
	}))

	/*mocker := mock.NewMocker(conn)
	if err := mocker.InsertUsers(1000); err != nil {
		panic(err)
	}*/

	server.Post("/register", h.Register)
	server.Post("/login", h.Login)
	server.Get("/users/me/orders", h.GetUserOrders)
	server.Get("/users/me", h.GetCurrentUser)

	server.Get("/restaurants/:id", h.GetRestaurantByID)
	server.Get("/restaurants/:id/menuData", h.GetRestaurantMenuData)
	server.Get("/restaurants", h.GetAllRestaurants)

	if err := server.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Println("An error occured, shutting down gracefully. ", err)
		_ = server.Shutdown()
	}
}
