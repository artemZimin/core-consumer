include .env
export $(shell sed 's/=.*//' .env)

build:
	go build ./cmd/app

fmt:
	go fmt ./...

generate-db-models:
	go tool github.com/xo/xo schema -o=internal/generated/models -k=field -e=created_at "$(DB_STRING)"

generate-di-container:
	go tool github.com/google/wire/cmd/wire ./internal/generated/container

generate-api:
	go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml --package=api docs/api.yml


