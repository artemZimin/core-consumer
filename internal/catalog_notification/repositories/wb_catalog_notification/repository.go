package wbcatalognotification

import "core-consumer/internal/app/gen/query"

type Repository struct {
	q *query.Query
}

func New(q *query.Query) *Repository {
	return &Repository{
		q: q,
	}
}
