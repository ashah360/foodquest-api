package handlers

import (
	"database/sql"
	"github.com/ashah360/fibertools"
	"github.com/ashah360/foodquest-api/internal/api/model"
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

func (h *HandlerGroup) GetBestSellingItems(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		fibertools.Message(c, fiber.StatusBadRequest, "ID required")
	}

	r, err := h.restaurantService.GetBestSellingItems(c.Context(), id, 7)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(r)
}

func (s *HandlerGroup) GetOrderCountForUsers(c *fiber.Ctx) error {

	var o []*model.UserOrders

	q := `SELECT users.id, users.email, users.first_name, users.last_name, users.phone_number, COUNT(order_details.id) as "orders"
FROM users
FULL OUTER JOIN order_details
ON order_details.user_id = users.id
GROUP BY users.id, users.email, users.first_name, users.last_name, users.phone_number
ORDER BY COUNT(order_details.id) DESC
`

	if err := s.db.SelectContext(c.Context(), &o, q); err != nil {
		if err == sql.ErrNoRows {
			o = []*model.UserOrders{}
		} else {
			panic(err)
		}
	}

	return c.Status(fiber.StatusOK).JSON(o)
}

func (s *HandlerGroup) GetFeaturedRestaurants(c *fiber.Ctx) error {
	var ids []string

	q := `SELECT DISTINCT(restaurant_id) 
FROM menu_item 
INNER JOIN menu ON menu.id = menu_item.menu_id
WHERE price < (SELECT AVG(price) FROM menu_item)
UNION
SELECT restaurant_id FROM ratings
GROUP BY restaurant_id
HAVING AVG(stars) > (select AVG(stars) from ratings)`

	if err := s.db.SelectContext(c.Context(), &ids, q); err != nil {
		if err == sql.ErrNoRows {
			ids = []string{}
		} else {
			panic(err)
		}
	}
	return c.Status(fiber.StatusOK).JSON(ids)

}
