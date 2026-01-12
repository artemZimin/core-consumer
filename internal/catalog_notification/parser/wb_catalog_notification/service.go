package wbcatalognotification

import (
	browserstorage "core-consumer/internal/app/storage/browser_storage"
)

type Service struct {
	browserStorage *browserstorage.Storage
}

func New(
	browserStorage *browserstorage.Storage,
) *Service {
	return &Service{
		browserStorage: browserStorage,
	}
}
