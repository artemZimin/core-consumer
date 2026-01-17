package wbstocknotification

import "context"

func (r *Repository) UpdateIsInStock(id int64, isInStock bool) error {
	m, err := r.q.WbStockNotification.WithContext(context.TODO()).Where(
		r.q.WbStockNotification.ID.Eq(id),
	).First()
	if err != nil {
		return err
	}

	m.IsInStock = isInStock

	return r.db.Save(m).Error
}
