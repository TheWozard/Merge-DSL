.PHONY: test clean start

test:
	 go test ./...

clean:
	git clean -fXd

start:
	GIN_MODE=release PORT=8080 go run cmd/server/main.go

