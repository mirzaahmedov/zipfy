include .env

build:
	@go build -o server cmd/main.go

run: build
	@exec ./server

migrate_up:
	@migrate -path ./migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PWD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DBNAME)?sslmode=$(POSTGRES_SSL)" up
