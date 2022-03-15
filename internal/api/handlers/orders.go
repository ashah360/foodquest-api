package handlers

import (
	"github.com/ashah360/fibertools"
	"github.com/ashah360/foodquest-api/internal/auth"
	"github.com/gofiber/fiber/v2"
)

func (h *HandlerGroup) GetUserOrders(c *fiber.Ctx) error {
	t := auth.ExtractJWT(c)

	if t == "" {
		return fibertools.Message(c, fiber.StatusUnauthorized, "User is not authenticated")
	}

	uid, err := auth.ValidateJWT(t)
	if err != nil {
		panic(err)
	}

	o, err := h.userService.GetUserOrders(c.Context(), uid)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(o)
}
