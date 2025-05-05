GOLANGCI_LINT_VERSION = v2.1.6

build:
	go build -o gtracer

test:
	go test ./... -race -shuffle=on -coverprofile=test-coverage.out
	go tool cover -html=test-coverage.out -o test-coverage.html

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run --fix

deps:
	@go mod tidy && go mod verify && go mod download
