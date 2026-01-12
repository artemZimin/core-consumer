package useragent

import (
	"core-consumer/internal/app/gen/query"

	"gorm.io/gorm"
)

type Repository struct {
	q  *query.Query
	db *gorm.DB
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
