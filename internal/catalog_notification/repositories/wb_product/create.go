package wbproduct

import "core-consumer/internal/app/gen/model"

func (r *Repository) Create(m *model.WbCatalogNotificationProduct) error {
	return r.db.Save(m).Error
}
