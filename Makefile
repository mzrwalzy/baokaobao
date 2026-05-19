.PHONY: build run test clean docker

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server -config config/config.yaml

test:
	go test -v ./...

clean:
	rm -rf bin/

docker-build:
	docker build -t baokaobao:latest .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

deps:
	go mod download
	go mod tidy
