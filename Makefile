PROJECT_NAME ?= test_task
PROJECT_VERSION ?= debug

.PHONY: build
build: \
	build-service

.PHONY: build-service
build-service:
	CGO_ENABLED=0 GOOS=linux go build -mod=readonly -a -installsuffix cgo \
		-o app/ \
		-ldflags " \
			-X 'main.ProjectName=${PROJECT_NAME}' \
			-X 'main.ProjectVersion=${PROJECT_VERSION}' \
		" \
		./cmd/service

.PHONY: lint
lint:
	docker run --rm -v .:/src -w /src golangci/golangci-lint:v1.56.2 golangci-lint run

.PHONY: test
test:
	go test ./...

.PHONY: up
up:
	docker-compose up -d
