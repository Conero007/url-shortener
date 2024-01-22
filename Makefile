build:
	@go build -o bin/shorten

run: build
	@./bin/shorten

test:
	@go test ./... -v