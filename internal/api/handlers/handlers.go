package handlers

import "github.com/jmoiron/sqlx"

type HandlerGroup struct {
	db *sqlx.DB
}

func NewHandlerGroup(db *sqlx.DB) *HandlerGroup {
	return &HandlerGroup{
		db,
	}
}
