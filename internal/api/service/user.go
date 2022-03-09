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

func NewUserService(db *sqlx.DB) UserService {
	return &userService{db}
}
