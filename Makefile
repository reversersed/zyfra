test:
	@go test -v -coverprofile=tests/coverage -coverpkg=./... ./... | findstr /V mocks && go tool cover -func=tests/coverage -o tests/coverage.func && go tool cover -html=tests/coverage -o tests/coverage.html

gen:
	@swag init --parseDependency -d ./internal/handlers -g ../app/app.go -o ./docs