package proxy

import (
	"core-consumer/internal/app/gen/model"
)

func (p *Repository) FindRanbom() (*model.Proxy, error) {
	var result model.Proxy

	if err := p.db.Raw(`
		select p.* from proxies p 
		order by random()
		limit 1;
		`).Scan(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
