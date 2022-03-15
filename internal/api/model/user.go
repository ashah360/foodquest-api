package model

import (
	"time"
)

type User struct {
	ID          string `json:"id" db:"id"`
	Email       string `json:"email" db:"email"`
	Password    string `json:"password" db:"password"`
	FirstName   string `json:"firstName" db:"first_name"`
	LastName    string `json:"lastName" db:"last_name"`
	PhoneNumber string `json:"phoneNumber" db:"phone_number"`

	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	LastLogin *time.Time `json:"lastLogin" db:"last_login"`
}

type PublicUser struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`

	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	LastLogin *time.Time `json:"lastLogin" db:"last_login"`
}

type UserOrders struct {
	ID          string `json:"id" db:"id"`
	Email       string `json:"email" db:"email"`
	FirstName   string `json:"firstName" db:"first_name"`
	LastName    string `json:"lastName" db:"last_name"`
	PhoneNumber string `json:"phoneNumber" db:"phone_number"`
	Orders      int    `json:"orders" db:"orders"`
}

func (u *User) ToPublicUser() *PublicUser {
	return &PublicUser{
		ID:          u.ID,
		Email:       u.Email,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		PhoneNumber: u.PhoneNumber,
		CreatedAt:   u.CreatedAt,
		LastLogin:   u.LastLogin,
	}
}
