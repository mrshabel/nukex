.PHONY: build run help

build:
	go build -o nukex.exe *.go

run:
	go run *.go

help:
	@echo "Available targets:"
	@echo "  build   - Compile the project and output the binary to the base directory. [nukex.exe]"
	@echo "  run     - Run the application"
	@echo "  help    - Display this help message"