package useragent

import (
	"core-consumer/internal/app/gen/model"
)

func (p *Repository) FindRandom() (*model.UserAgent, error) {
	var result model.UserAgent

	if err := p.db.Raw(`
		select * from user_agents ua
		order by random()
		limit 1;
		`).Scan(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
