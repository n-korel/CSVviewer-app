.PHONY: up down test

up:
	docker-compose up --build

down:
	docker-compose down

test:
	go test ./... -v -cover

generate-csv:
	go run ./utils/generate_csv.go