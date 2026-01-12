package wbproduct

import (
	"context"
	"core-consumer/internal/app/gen/model"
)

type FindByUrlAndPriceInCatalogNotificationParams struct {
	NotificationID int64
	URL            string
	Price          int32
}

func (r *Repository) FindByUrlAndPriceInCatalogNotification(
	params FindByUrlAndPriceInCatalogNotificationParams,
) (*model.WbCatalogNotificationProduct, error) {
	return r.q.WbCatalogNotificationProduct.WithContext(context.TODO()).
		Where(
			r.q.WbCatalogNotificationProduct.URL.Eq(params.URL),
			r.q.WbCatalogNotificationProduct.Price.Eq(params.Price),
			r.q.WbCatalogNotificationProduct.WbCatalogNotificationID.Eq(params.NotificationID),
		).
		First()
}
