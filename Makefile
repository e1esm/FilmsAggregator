run:
	go run ./cmd/aggregator/*.go
dockerize:
	docker compose up --build -d
test:
	go test ./... -cover