package wbcatalognotification

import (
	"context"
	"core-consumer/internal/app/gen/model"
)

func (r *Repository) FindByID(id int64) (*model.WbCatalogNotification, error) {
	return r.q.WithContext(context.TODO()).WbCatalogNotification.Where(
		r.q.WbCatalogNotification.ID.Eq(id),
	).First()

}
