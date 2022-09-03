.PHONY: test clean start coverage

test:
	 go test -count=1 -v ./...

coverage:
	go test -coverprofile=coverage.out  ./...
	go tool cover -html=coverage.out

clean:
	git clean -fXd
