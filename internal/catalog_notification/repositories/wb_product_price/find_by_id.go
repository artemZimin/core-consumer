package wbproductprice

import (
	"context"
	"core-consumer/internal/app/gen/model"
)

func (r *Repository) FindByID(id int64) (*model.WbProductPrice, error) {
	return r.q.WbProductPrice.WithContext(context.TODO()).Where(
		r.q.WbProductPrice.ID.Eq(id),
	).First()
}
