.PHONY: test bench clean start coverage

test:
	 go test -count=1 ./...

bench:
	go test -bench=. -run=^# -benchtime=1s ./pkg/example_test.go

coverage:
	go test -coverprofile=coverage.out  ./...
	go tool cover -html=coverage.out

clean:
	git clean -fXd
