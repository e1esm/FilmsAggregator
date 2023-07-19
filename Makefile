run:
	go run ./cmd/aggregator/*.go
dockerize:
	docker compose up --build -d