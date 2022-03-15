package service

import (
	"context"
	"database/sql"
	"github.com/ashah360/foodquest-api/internal/api/cerror"
	"github.com/ashah360/foodquest-api/internal/api/db"
	"github.com/ashah360/foodquest-api/internal/api/model"
	"github.com/jmoiron/sqlx"
)

type RestaurantService interface {
	GetRestaurantsOwnedByUser(ctx context.Context, userID string) ([]*model.Restaurant, error)
	GetRestaurantByID(ctx context.Context, restaurantID string) (*model.Restaurant, error)
	GetRestaurantPageData(ctx context.Context, restaurantID string) (*model.RestaurantPageData, error)
	GetRestaurantMenuItems(ctx context.Context, restaurantID string) ([]*model.RestaurantMenuItem, error)
	GetRestaurantMenuData(ctx context.Context, restaurantID string) (*RestaurantMenuData, error)
	GetAllRestaurants(ctx context.Context) ([]*model.Restaurant, error)
}

type restaurantService struct {
	db *sqlx.DB
}

func (s *restaurantService) GetAllRestaurants(ctx context.Context) ([]*model.Restaurant, error) {
	var r []*model.Restaurant

	q := `SELECT * FROM restaurant`

	if err := s.db.SelectContext(ctx, &r, q); err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	return r, nil
}

func (s *restaurantService) GetRestaurantMenuItems(ctx context.Context, restaurantID string) ([]*model.RestaurantMenuItem, error) {
	var items []*model.RestaurantMenuItem

	q := `SELECT restaurant.id as "restaurant_id", menu.id as "menu_id", menu.menu_name as "menu_name", menu.description as "menu_description", menu_item.id as "menu_item_id", menu_item.section, menu_item.title, menu_item.item_description, menu_item.price, menu_item.image_url FROM restaurant 
INNER JOIN menu ON restaurant.id = menu.restaurant_id
INNER JOIN menu_item ON menu.id = menu_item.menu_id WHERE restaurant.id=$1`

	if err := s.db.SelectContext(ctx, &items, q, restaurantID); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return items, nil
}

type MenuData struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description db.NullString   `json:"description"`
	Sections    []db.NullString `json:"sections" db:"sections"`
	MenuItemIDs []string        `json:"menuItemIds"`
}

type RestaurantMenuData struct {
	RestaurantID string                     `json:"restaurantId"`
	MenuIDs      []string                   `json:"menuIds"`
	Menus        map[string]*MenuData       `json:"menus"`
	MenuItems    map[string]*model.MenuItem `json:"menuItems"`
}

func (s *restaurantService) GetRestaurantMenuData(ctx context.Context, restaurantID string) (*RestaurantMenuData, error) {
	i, err := s.GetRestaurantMenuItems(ctx, restaurantID)
	if err != nil {
		panic(err)
	}

	d := &RestaurantMenuData{
		MenuIDs:   []string{},
		Menus:     make(map[string]*MenuData),
		MenuItems: make(map[string]*model.MenuItem),
	}

	for _, m := range i {
		if d.RestaurantID == "" {
			d.RestaurantID = m.RestaurantID
		}

		if _, exists := d.Menus[m.MenuID]; !exists {
			d.Menus[m.MenuID] = &MenuData{
				ID:          m.MenuID,
				Title:       m.MenuName,
				Sections:    []db.NullString{},
				Description: m.MenuDescription,
				MenuItemIDs: []string{},
			}
			d.MenuIDs = append(d.MenuIDs, m.MenuID)
		}

		d.MenuItems[m.MenuItemID] = &model.MenuItem{
			ID:              m.MenuItemID,
			MenuID:          m.MenuID,
			Section:         m.Section,
			Title:           m.Title,
			ItemDescription: m.ItemDescription,
			Price:           m.Price,
			ImageURL:        m.ImageURL,
		}

		d.Menus[m.MenuID].MenuItemIDs = append(d.Menus[m.MenuID].MenuItemIDs, m.MenuItemID)
	}

	for _, mid := range d.MenuIDs {
		if err := s.db.SelectContext(ctx, &(d.Menus[mid].Sections), `SELECT DISTINCT(section) FROM menu_item WHERE menu_id=$1`, mid); err != nil && err != sql.ErrNoRows {
			panic(err)
		}
	}

	return d, nil
}

func (s *restaurantService) GetRestaurantPageData(ctx context.Context, restaurantID string) (*model.RestaurantPageData, error) {
	var r model.RestaurantPageData

	q := `SELECT id, owner_id, restaurant_name, category, open_time, close_time, address_line_1, address_line_2, address_line_3, city, state, postal_code, country, image_url, stars, num_ratings FROM restaurant LEFT JOIN (SELECT restaurant_id, ROUND(AVG(stars), 2) as "stars", COUNT(stars) as "num_ratings" FROM ratings GROUP BY restaurant_id) AS R on R.restaurant_id = restaurant.id WHERE id=$1`

	if err := s.db.GetContext(ctx, &r, q, restaurantID); err != nil {
		if err == sql.ErrNoRows {
			return nil, cerror.ErrResourceNotFound
		}
		return nil, err
	}

	return &r, nil
}

func (s *restaurantService) GetRestaurantByID(ctx context.Context, restaurantID string) (*model.Restaurant, error) {
	var r model.Restaurant

	q := `SELECT * FROM restaurant WHERE id=$1`
	if err := s.db.GetContext(ctx, &r, q, restaurantID); err != nil {
		if err == sql.ErrNoRows {
			return nil, cerror.ErrResourceNotFound
		}
		return nil, err
	}

	return &r, nil
}

func (s *restaurantService) GetRestaurantsOwnedByUser(ctx context.Context, userID string) ([]*model.Restaurant, error) {
	var r []*model.Restaurant

	if err := s.db.SelectContext(ctx, &r, `SELECT * FROM restaurant WHERE owner_id=$1`, userID); err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	return r, nil
}

func NewRestaurantService(db *sqlx.DB) RestaurantService {
	return &restaurantService{db}
}
