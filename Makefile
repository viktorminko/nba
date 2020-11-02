.PHONY: all
all: deps build lint unit_test system_test

.PHONY: deps
deps:
	go mod tidy
	go mod download
	go mod vendor

.PHONY: build_simulation
build_simulation:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -o artifacts/svc ./cmd/simulation

.PHONY: build_statistic
build_statistic:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -o artifacts/svc ./cmd/statistic

.PHONY: build
build: build_simulation build_statistic

.PHONY: unit_test
unit_test:
	go test -count=1 -v -race -cover -mod=vendor ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: up_simulation
up_simulation:
	docker-compose up --build simulation

.PHONY: up_statistic
up_statistic:
	docker-compose up --build statistic

.PHONY: up
up:
	docker-compose up --build

.PHONY: down
down:
	docker-compose  down --volumes --remove-orphans
	docker-compose  rm --force --stop -v
