build:
	@go build -o bin/shorten

run: build test
	@./bin/shorten

test: build
	@go test ./app -v