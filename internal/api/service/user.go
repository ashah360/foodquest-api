package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ashah360/foodquest-api/internal/api/cerror"
	"github.com/ashah360/foodquest-api/internal/api/model"
	"github.com/jmoiron/sqlx"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]*model.User, error)
	GetUserByID(ctx context.Context, ID string) (*model.User, error)
	GetUserOrders(ctx context.Context, id string) ([]*model.OrderInvoice, error)
}

type userService struct {
	db *sqlx.DB
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	var users []*model.User

	if err := s.db.SelectContext(ctx, &users, fmt.Sprintf("select * from %s", "users")); err != nil {
		if err == sql.ErrNoRows {
			return users, nil
		}

		return nil, err
	}

	return users, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var u model.User

	if err := s.db.GetContext(ctx, &u, fmt.Sprintf("select * from %s where id=$1", "users"), id); err != nil {
		if err == sql.ErrNoRows {
			return nil, cerror.ErrUserDoesNotExist
		}

		return nil, err
	}

	return &u, nil
}

func (s *userService) GetUserOrders(ctx context.Context, id string) ([]*model.OrderInvoice, error) {
	var ord []*model.OrderInvoice

	q := `SELECT O.id, user_id, restaurant_id, restaurant_name, created_at, total, confirmed, status, estimated_mins 
FROM order_details O
LEFT JOIN restaurant R
ON R.id = O.restaurant_id 
WHERE O.user_id=$1
ORDER BY created_at DESC`

	if err := s.db.SelectContext(ctx, &ord, q, id); err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	q = `SELECT order_id, menu_item_id, title, item_description, quantity, price, image_url FROM line_item LEFT JOIN menu_item ON line_item.menu_item_id = menu_item.id WHERE order_id IN (SELECT id FROM order_details WHERE order_details.id=$1)`

	for _, o := range ord {
		if err := s.db.SelectContext(ctx, &(o.LineItems), q, o.ID); err != nil && err != sql.ErrNoRows {
			panic(err)
		}
		if o.LineItems == nil {
			o.LineItems = []*model.OrderLineItem{}
		}

	}

	return ord, nil
}

func NewUserService(db *sqlx.DB) UserService {
	return &userService{db}
}
