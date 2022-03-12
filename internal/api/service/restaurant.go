package service

import (
	"context"
	"database/sql"
	"github.com/ashah360/foodquest-api/internal/api/model"
	"github.com/jmoiron/sqlx"
)

type RestaurantService interface {
	GetRestaurantsOwnedByUser(ctx context.Context, userID string) ([]*model.Restaurant, error)
	CreateRestaurant(ctx context.Context, userID string, opts *CreateRestaurantOpts) (*model.Restaurant, error)
}

type restaurantService struct {
	db *sqlx.DB
}

func (s *restaurantService) GetRestaurantsOwnedByUser(ctx context.Context, userID string) ([]*model.Restaurant, error) {
	var r []*model.Restaurant

	if err := s.db.SelectContext(ctx, &r, `SELECT * FROM restaurant WHERE owner_id=$1`, userID); err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	return r, nil
}

type CreateRestaurantOpts struct {
}
