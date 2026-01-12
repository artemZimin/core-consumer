package wbproduct

import (
	"core-consumer/internal/app/gen/query"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
	q  *query.Query
}

func New(
	db *gorm.DB,
	q *query.Query,
) *Repository {
	return &Repository{
		db: db,
		q:  q,
	}
}
