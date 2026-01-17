package wbstocknotification

import (
	"core-consumer/internal/app/gen/query"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
	q  *query.Query
}

func New(
	q *query.Query,
	db *gorm.DB,
) *Repository {
	return &Repository{
		q:  q,
		db: db,
	}
}
