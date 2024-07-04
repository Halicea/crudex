build: tidy generate lint
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
	cd examples/simple && go run main.go

tidy:
	go mod tidy

generate:
	go generate ./...

lint:
	golangci-lint run

