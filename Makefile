build:
	go build

run:
	APP_ENV=dev go run main.go

lint:
	golangci-lint run -c .dev/.golangci.yml

unit-test:
	go clean --testcache && go test ./...

code-coverage:
	go test `go list ./... | grep -v /tilt_modules` -coverprofile cover.out
	go tool cover -html=cover.out -o coverage.html
	echo `go tool cover -func cover.out | grep total`

generate-mocks:
	mockgen -source internal/gateway/handler.go -destination internal/gateway/mock_handler_test.go -package gateway_test
	mockgen -source internal/gateway/service.go -destination internal/gateway/mock_service_test.go -package gateway_test

