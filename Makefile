include .env
export $(shell sed 's/=.*//' .env)

build:
	go build -o core-consumer ./cmd/app

start:
	nohup ./core-consumer > nohup.log 2>&1 &

fmt:
	go fmt ./...

generate-db-models:
	go tool gorm.io/gen/tools/gentool -db postgres -dsn "host=localhost port=5432 user=sail password=password dbname=laravel sslmode=disable" -outPath "internal/app/gen/query" -fieldNullable true



generate-di-container:
	go tool github.com/google/wire/cmd/wire ./internal/generated/container

generate-api:
	go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml --package=api docs/api.yml


