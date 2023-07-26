run:
	go run ./cmd/aggregator/*.go
build_dockerize:
	docker compose up --build -d
dockerize:
	docker compose up -d
test:
	go test ./... -cover