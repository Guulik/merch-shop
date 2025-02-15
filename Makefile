PHONY: lint

lint:
	golangci-lint run -c ./.golangci.yml

start:
	docker-compose up -d

down:
	docker-compose down

stop:
	docker-compose stop

test:
	go test ./...

test_integration:
	go test -tags integration ./tests/integration

coverage_cli:
	go test -coverprofile="coverage.out" ./...
	go tool cover -func="coverage.out"

coverage_html:
	go test -coverprofile="coverage.out" ./...
	go tool cover -html="coverage.out"

