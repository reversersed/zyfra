test:
	@go test -v -coverprofile=tests/coverage -coverpkg=./... ./... | findstr /V mocks && go tool cover -func=tests/coverage -o tests/coverage.func && go tool cover -html=tests/coverage -o tests/coverage.html

gen:
	@go install github.com/golang/mock/mockgen@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	
	@go generate ./...
	@swag init --parseDependency -d ./internal/handlers -g ../app/app.go -o ./docs

test-verbose: gen
	@go test -v ./... | grep -v mocks

check: gen
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run
