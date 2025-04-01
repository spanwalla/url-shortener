include .env
export

compose-up: ### Run docker-compose
	docker-compose up --build -d && docker-compose logs -f
.PHONY: compose-up

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

docker-rm-volume: ### Remove docker volume
	docker volume rm url-shortener_postgres_data
.PHONY: docker-rm-volume

migrate-create:  ### Create new migration
	migrate create -ext sql -dir migrations 'url-shortener'
.PHONY: migrate-create

migrate-up: ### Migration up
	migrate -path migrations -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:{$(POSTGRES_PORT)}/$(POSTGRES_DB)?sslmode=disable' up
.PHONY: migrate-up

migrate-down: ### Migration down
	echo "y" | migrate -path migrations -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:{$(POSTGRES_PORT)}/$(POSTGRES_DB)?sslmode=disable' down
.PHONY: migrate-down

linter-golangci: ### Check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

swag: ### Generate swagger docs
	go tool github.com/swaggo/swag/cmd/swag init -g 'internal/app/app.go' --parseInternal --parseDependency
.PHONY: swag

test: ### Run test
	go test -v './...'
.PHONY: test

mockgen: ### Generate mock
	go tool go.uber.org/mock/mockgen -source='internal/service/service.go'       -destination='internal/mocks/service/mock.go'    -package=servicemocks
	go tool go.uber.org/mock/mockgen -source='internal/repository/repository.go' -destination='internal/mocks/repository/mock.go' -package=repomocks
	go tool go.uber.org/mock/mockgen -source='pkg/encoder/encoder.go'            -destination='internal/mocks/encoder/mock.go'    -package=encodermocks
.PHONY: mockgen

bin-deps: ### Install binary dependencies
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.0.2
.PHONY: bin-deps