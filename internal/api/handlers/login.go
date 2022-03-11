package handlers

import (
	"database/sql"
	"github.com/ashah360/fibertools"
	"github.com/ashah360/foodquest-api/internal/api/cerror"
	"github.com/ashah360/foodquest-api/internal/api/model"
	"github.com/ashah360/foodquest-api/internal/auth"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *HandlerGroup) Login(c *fiber.Ctx) error {
	p := new(LoginPayload)

	if err := c.BodyParser(p); err != nil {
		return fibertools.Message(c, fiber.StatusBadRequest, "Invalid login format")
	}

	var u model.User

	if err := h.db.GetContext(c.Context(), &u, "select * from users where lower(email)=lower($1)", p.Email); err != nil {
		if err == sql.ErrNoRows {
			panic(cerror.ErrUserDoesNotExist)
		}

		panic(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p.Password)); err != nil {
		return fibertools.Message(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	t, err := auth.IssueToken(u.ID)
	if err != nil {
		panic(err)
	}

	if _, err := h.db.ExecContext(c.Context(), `update users set last_login=$1 where id=$2`, time.Now(), u.ID); err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": t,
	})

}
