.PHONY: build run help

build:
	go build -o bin/nukex main.go

run:
	go run *.go

help:
	@echo "Available targets:"
	@echo "  build   - Compile the project and output the binary to the bin/ directory"
	@echo "  run     - Run the application"
	@echo "  help    - Display this help message"