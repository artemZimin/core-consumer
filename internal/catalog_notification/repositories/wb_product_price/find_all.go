package wbproductprice

import (
	"context"
	"core-consumer/internal/app/gen/model"
)

func (r *Repository) FindAll() ([]*model.WbProductPrice, error) {
	return r.q.WbProductPrice.WithContext(context.TODO()).Find()
}
