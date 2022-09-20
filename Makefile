.PHONY: test bench brew clean start coverage lint

test:
	 go test -count=1 ./...

bench:
	go test -bench=. -run=^# -benchtime=1s ./pkg/example_test.go

brew:
	brew install golangci-lint

coverage:
	go test -coverprofile=coverage.out  ./...
	go tool cover -html=coverage.out

clean:
	git clean -fXd

lint:
	golangci-lint run
