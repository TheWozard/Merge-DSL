.PHONY: test bench bench-diff brew clean start coverage lint
tests = ./...
n = 5
duration = 1s

test:
	 go test -count=1 $(tests)

bench:
	@if [ -f benchmark.out  ] ; then \
		mv benchmark.out benchmark.prev.out; \
	fi
	go test -bench=. -run=^# -benchtime=$(duration) -benchmem -count=$(n) $(tests) | tee benchmark.out
	@if [ -f benchmark.prev.out  ] ; then \
		echo ; \
		echo ; \
		echo "📋 \033[0;34m=[ benchstat ]===========================================================\033[0m"; \
		benchstat -split pkg -sort delta benchmark.prev.out benchmark.out; \
	fi

deps:
	go install golang.org/x/perf/cmd/benchstat@latest
	brew install golangci-lint

coverage:
	go test -coverprofile=coverage.out  $(tests)
	go tool cover -html=coverage.out

clean:
	git clean -fXd

lint:
	golangci-lint run
