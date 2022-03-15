package handlers

import (
	"github.com/ashah360/fibertools"
	"github.com/gofiber/fiber/v2"
)

func (h *HandlerGroup) GetAllRestaurants(c *fiber.Ctx) error {
	r, err := h.restaurantService.GetAllRestaurants(c.Context())
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(r)
}

func (h *HandlerGroup) GetRestaurantByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		fibertools.Message(c, fiber.StatusBadRequest, "ID required")
	}

	r, err := h.restaurantService.GetRestaurantPageData(c.Context(), id)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(r)
}

func (h *HandlerGroup) GetRestaurantMenuData(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		fibertools.Message(c, fiber.StatusBadRequest, "ID required")
	}

	r, err := h.restaurantService.GetRestaurantMenuData(c.Context(), id)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(r)
}
