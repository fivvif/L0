package repository

import (
	"L0"
	"github.com/jmoiron/sqlx"
)

type Orders interface {
	SaveOrder(order L0.Order) error
	RecoverCache() ([]L0.Order, error)
}

type Repository struct {
	Orders
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{Orders: NewOrdersPostgres(db)}
}
