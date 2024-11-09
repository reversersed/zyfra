test:
	@go test -v -coverprofile=tests/coverage -coverpkg=./... ./... | findstr /V mocks && go tool cover -func=tests/coverage -o tests/coverage.func && go tool cover -html=tests/coverage -o tests/coverage.html

gen:
	@go generate ./...
	@swag init --parseDependency -d ./internal/handlers -g ../app/app.go -o ./docs

test-verbose:
	@go test -v ./... | grep -v mocks

check:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run
