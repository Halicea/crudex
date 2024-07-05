build: tidy
	go build ./...

run: build
	go run main.go

test: build
	go test -v  ./...

coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

watch: build
	air -c .air.toml

example: build
	cd examples/simple && go build -o simple
	cd examples/simple && ./simple

tidy:
	go mod tidy

generate:
	go generate ./...

lint:
	golangci-lint run
