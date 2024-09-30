.PHONY: compose.up consume produce.json
compose.up:
	docker compose up -d
consume:
	go run cmd/consumer/main.go
produce.json:
	go run cmd/producer/main.go
