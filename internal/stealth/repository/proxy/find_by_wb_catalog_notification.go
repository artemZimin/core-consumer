package proxy

import (
	"core-consumer/internal/app/gen/model"
)

func (p *Repository) FindByWbCatalogNotification(wcnID int64) (*model.Proxy, error) {
	var result model.Proxy

	if err := p.db.Raw(`
		select p.* from proxies p 
		join wb_catalog_notification_proxy wcnp on wcnp.proxy_id = p.id
		join wb_catalog_notifications wcn on wcn.id = wcnp.wb_catalog_notification_id 
		where wcn.id = ?
		order by random()
		limit 1;
		`, wcnID).Scan(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
