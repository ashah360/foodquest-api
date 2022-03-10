package main

import (
	"database/sql"
	"fmt"
	"github.com/ashah360/fibertools"
	"github.com/ashah360/foodquest-api/internal/api/cerror"
	"github.com/ashah360/foodquest-api/internal/api/db"
	"github.com/ashah360/foodquest-api/internal/api/model"
	"github.com/ashah360/foodquest-api/internal/api/service"
	"github.com/ashah360/foodquest-api/internal/auth"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/mail"
	"os"
	"regexp"
	"strings"
	"time"
)

var port = "3000"

type RegisterUserPayload struct {
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

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

	userService := service.NewUserService(conn)

	app := fiber.New(fiber.Config{
		ErrorHandler: fibertools.ErrorHandler,
	})

	app.Use(fibertools.Recover())
	app.Use(cors.New())

	app.Get("/users", func(c *fiber.Ctx) error {
		u, err := userService.GetAllUsers(c.Context())
		if err != nil {
			panic(err)
		}
		return c.Status(fiber.StatusOK).JSON(u)
	})

	app.Post("/register", func(c *fiber.Ctx) error {
		var p RegisterUserPayload

		if err := c.BodyParser(&p); err != nil {
			return fibertools.Message(c, fiber.StatusBadRequest, "Invalid registration format")
		}

		if p.FirstName == "" {
			return fibertools.Message(c, fiber.StatusBadRequest, "First name cannot be blank")
		}

		if p.LastName == "" {
			return fibertools.Message(c, fiber.StatusBadRequest, "Last name cannot be blank")
		}

		if p.PhoneNumber == "" {
			return fibertools.Message(c, fiber.StatusBadRequest, "Phone number cannot be blank")
		}

		re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)

		if !re.MatchString(p.PhoneNumber) {
			return fibertools.Message(c, fiber.StatusBadRequest, "Invalid phone number")
		}

		if _, err := mail.ParseAddress(p.Email); err != nil {
			return fibertools.Message(c, fiber.StatusBadRequest, "Invalid email address")
		}

		var userCount int

		q := fmt.Sprintf("select count(*) from users where lower(users.email)='%s'", strings.ToLower(p.Email))

		if err := conn.QueryRow(q).Scan(&userCount); err != nil {
			panic(err)
		}

		if userCount > 0 {
			return fibertools.Message(c, fiber.StatusBadRequest, "Account already registered with email address")
		}

		pw, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		var u model.User

		q = `insert into users (email, password, first_name, last_name, phone_number) values ($1, $2, $3, $4, $5) returning *`

		if err := conn.GetContext(c.Context(), &u, q, p.Email, string(pw), p.FirstName, p.LastName, p.PhoneNumber); err != nil {
			panic(err)
		}

		return c.Status(fiber.StatusCreated).JSON(u)
	})

	type LoginPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	app.Post("/login", func(c *fiber.Ctx) error {
		p := new(LoginPayload)

		if err := c.BodyParser(p); err != nil {
			return fibertools.Message(c, fiber.StatusBadRequest, "Invalid login format")
		}

		var u model.User

		if err := conn.GetContext(c.Context(), &u, "select * from users where lower(email)=lower($1)", p.Email); err != nil {
			if err == sql.ErrNoRows {
				panic(cerror.ErrUserDoesNotExist)
			}

			panic(err)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p.Password)); err != nil {
			return fibertools.Message(c, fiber.StatusUnauthorized, "Invalid credentials")
		}

		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = u.ID
		claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
		claims["iat"] = time.Now().Unix()

		t, err := token.SignedString([]byte(os.Getenv("FOODQUEST_JWT_SECRET")))
		if err != nil {
			panic(err)
		}

		if _, err := conn.ExecContext(c.Context(), `update users set last_login=$1 where id=$2`, time.Now(), u.ID); err != nil {
			panic(err)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"token": t,
		})
	})

	app.Get("/users/me", func(c *fiber.Ctx) error {
		t := auth.ExtractJWT(c)
		if t == "" {
			return fibertools.Message(c, fiber.StatusUnauthorized, "User is not authenticated")
		}

		uid, err := auth.ValidateJWT(t)
		if err != nil {
			panic(err)
		}

		var u model.User

		if err := conn.GetContext(c.Context(), &u, "select * from users where id=$1", uid); err != nil {
			if err == sql.ErrNoRows {
				panic(cerror.ErrUserDoesNotExist)
			}

			panic(err)
		}

		return c.Status(200).JSON(u.ToSimple())

	})

	type RestaurantQuery struct {
		OwnerEmail     string `json:"owner_email" db:"owner_email"`
		RestaurantName string `json:"restaurant_name" db:"restaurant_name"`
	}

	app.Get("/restaurants", func(c *fiber.Ctx) error {
		var r []*RestaurantQuery

		if err := conn.SelectContext(c.Context(), &r, `select users.email as owner_email, restaurant.name as restaurant_name from users right join restaurant on users.id = restaurant.owner_id`); err != nil && err != sql.ErrNoRows {
			panic(err)
		}

		return c.Status(fiber.StatusOK).JSON(r)
	})

	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Println("An error occured, shutting down gracefully. ", err)
		_ = app.Shutdown()
	}
}
