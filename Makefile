build:
	go build -o gtracer

test:
	go test ./... -v -coverprofile=test-coverage.out
	go tool cover -func=test-coverage.out

lint:
	golangci-lint run

deps:
	go mod tidy
	go mod verify
