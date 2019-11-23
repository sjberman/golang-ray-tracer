test:
	go test ./... -v -coverprofile=test-coverage.out
	go tool cover -func=test-coverage.out

lint:
	golangci-lint run --enable-all -D gochecknoglobals -D wsl

deps:
	go mod tidy
	go mod verify