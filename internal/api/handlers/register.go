package handlers

import (
	"fmt"
	"github.com/ashah360/fibertools"
	"github.com/ashah360/foodquest-api/internal/api/model"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"regexp"
	"strings"
)

type RegisterUserPayload struct {
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phoneNumber"`
}

func (h *HandlerGroup) Register(c *fiber.Ctx) error {
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

	if err := h.db.QueryRow(q).Scan(&userCount); err != nil {
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

	if err := h.db.GetContext(c.Context(), &u, q, p.Email, string(pw), p.FirstName, p.LastName, p.PhoneNumber); err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusCreated).JSON(u.ToPublicUser())

}
