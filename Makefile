build:
	go build ./...

run:
	go run main.go

test:
	go test -v ./...

tidy:
	go mod tidy
	cd examples/simple && go mod tidy
	
watch:
	air -c .air.toml

example:
	cd examples/simple && go run main.go
