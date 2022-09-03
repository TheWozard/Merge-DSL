.PHONY: test clean start coverage

test:
	 go test ./...

coverage:
	go test -coverprofile=coverage.out  ./...
	go tool cover -html=coverage.out

clean:
	git clean -fXd

start:
	GIN_MODE=release PORT=8080 go run cmd/server/main.go

