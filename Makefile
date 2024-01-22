build:
	@go build -o bin/shorten

run: build
	@./bin/shorten

test: build
	@go test ./... -v