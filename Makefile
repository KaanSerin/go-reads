build:
	@go build -o ./bin ./...

run: build
	@./bin/go_reads