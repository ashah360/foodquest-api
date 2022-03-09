package model

import (
	"github.com/ashah360/foodquest-api/internal/api/db"
	"time"
)

type User struct {
	ID          string `json:"id" db:"id"`
	Email       string `json:"email" db:"email"`
	Password    string `json:"password" db:"password"`
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`

	// Address
	AddressLine1 db.NullString `json:"address_line_1" db:"address_line_1"`
	AddressLine2 db.NullString `json:"address_line_2" db:"address_line_2"`
	AddressLine3 db.NullString `json:"address_line_3" db:"address_line_3"`
	City         db.NullString `json:"city" db:"city"`
	State        db.NullString `json:"state" db:"state"`
	PostalCode   db.NullString `json:"postal_code" db:"postal_code"`
	Country      db.NullString `json:"country" db:"country"`

	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	LastLogin *time.Time `json:"last_login" db:"last_login"`
}

type SimpleUser struct {
	ID          string `json:"id" db:"id"`
	Email       string `json:"email" db:"email"`
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`

	// Address
	AddressLine1 db.NullString `json:"address_line_1,omitempty" db:"address_line_1"`
	AddressLine2 db.NullString `json:"address_line_2,omitempty" db:"address_line_2"`
	AddressLine3 db.NullString `json:"address_line_3,omitempty" db:"address_line_3"`
	City         db.NullString `json:"city,omitempty" db:"city"`
	State        db.NullString `json:"state,omitempty" db:"state"`
	PostalCode   db.NullString `json:"postal_code,omitempty" db:"postal_code"`
	Country      db.NullString `json:"country,omitempty" db:"country"`

	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	LastLogin *time.Time `json:"last_login" db:"last_login"`
}

func (u *User) ToSimple() *SimpleUser {
	return &SimpleUser{
		ID:           u.ID,
		Email:        u.Email,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		PhoneNumber:  u.PhoneNumber,
		AddressLine1: u.AddressLine1,
		AddressLine2: u.AddressLine2,
		AddressLine3: u.AddressLine3,
		City:         u.City,
		State:        u.State,
		PostalCode:   u.PostalCode,
		Country:      u.Country,
		CreatedAt:    u.CreatedAt,
		LastLogin:    u.LastLogin,
	}
}
