package wbstocknotification

func (r *Repository) CountInStatus(status string) (int64, error) {
	var res int64

	if err := r.db.Raw(
		"select count(*) from wb_stock_notifications wsn where wsn.status = ?;",
		status,
	).Scan(&res).Error; err != nil {
		return 0, err
	}

	return res, nil
}
