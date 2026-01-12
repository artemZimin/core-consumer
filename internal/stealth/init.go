package stealth

import (
	"core-consumer/internal/app/gen/query"
	"core-consumer/internal/stealth/repository/proxy"
	useragent "core-consumer/internal/stealth/repository/user_agent"

	"gorm.io/gorm"
)

type Module struct {
	ProxyRepo     *proxy.Repository
	UserAgentRepo *useragent.Repository
}

func Init(
	db *gorm.DB,
	q *query.Query,
) *Module {
	return &Module{
		ProxyRepo:     proxy.New(q, db),
		UserAgentRepo: useragent.New(q, db),
	}
}
