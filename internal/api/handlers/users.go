package handlers

import (
	"database/sql"
	"github.com/ashah360/fibertools"
	"github.com/ashah360/foodquest-api/internal/api/cerror"
	"github.com/ashah360/foodquest-api/internal/api/model"
	"github.com/ashah360/foodquest-api/internal/auth"
	"github.com/gofiber/fiber/v2"
)

func (h *HandlerGroup) GetCurrentUser(c *fiber.Ctx) error {
	t := auth.ExtractJWT(c)

	if t == "" {
		return fibertools.Message(c, fiber.StatusUnauthorized, "User is not authenticated")
	}

	uid, err := auth.ValidateJWT(t)
	if err != nil {
		panic(err)
	}

	var u model.User

	if err := h.db.GetContext(c.Context(), &u, "select * from users where id=$1", uid); err != nil {
		if err == sql.ErrNoRows {
			panic(cerror.ErrUserDoesNotExist)
		}

		panic(err)
	}

	return c.Status(200).JSON(u.ToPublicUser())

}
