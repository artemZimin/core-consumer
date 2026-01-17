package wbstocknotification

import (
	"context"
	"core-consumer/internal/app/gen/model"
)

func (r *Repository) FindByID(id int64) (*model.WbStockNotification, error) {
	return r.q.WithContext(context.TODO()).WbStockNotification.Where(
		r.q.WbStockNotification.ID.Eq(id),
	).First()

}
