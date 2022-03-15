package handlers

import (
	"database/sql"
	"github.com/ashah360/fibertools"
	"github.com/ashah360/foodquest-api/internal/api/model"
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

func (h *HandlerGroup) GetRecentlyPurchasedItems(c *fiber.Ctx) error {
	var m []*model.MenuItem

	q := `SELECT *
FROM menu_item M
WHERE M.id IN (
	SELECT L.menu_item_id 
	FROM line_item L 
	WHERE L.menu_item_id = M.id
	AND (SELECT status FROM order_details WHERE order_details.id = L.order_id AND order_details.created_at > NOW() - INTERVAL '72 HOUR') != 'cancelled' 	
)`

	if err := h.db.SelectContext(c.Context(), &m, q); err != nil {
		if err == sql.ErrNoRows {
			m = []*model.MenuItem{}
		} else {
			panic(err)
		}
	}

	return c.Status(200).JSON(m)
}
