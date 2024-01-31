BINARY_NAME ?= go-migrate
TAG ?= $(shell git rev-parse HEAD 2>/dev/null || echo "default-value")

build:
	go build -mod vendor -o bin/$(BINARY_NAME) main.go

docker-build:
	docker build -t pbxx/migrate:$(TAG) .

docker-push:
	docker push pbxx/migrate:$(TAG)

docker-build-push: docker-build docker-push
