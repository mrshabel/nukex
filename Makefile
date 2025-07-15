BINARY_NAME=nukex
BUILD_DIR=bin

.PHONY: build run help clean build-all build-windows build-mac-intel build-mac-arm build-linux build-linux-arm

# create build directory
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# build for current platform
build: $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME).exe *.go

# cross-compile targets by disabling cgo and setting target os
build-windows: $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe *.go

build-mac-intel: $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 *.go

build-mac-arm: $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 *.go

build-linux: $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 *.go

build-linux-arm: $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 *.go

# build all platforms
build-all: build-windows build-mac-intel build-mac-arm build-linux build-linux-arm

# clean build directory
clean:
	rm -rf $(BUILD_DIR)

run:
	go run *.go

help:
	@echo "Available targets:"
	@echo "  build        - Compile the project and output the binary to the build directory. [nukex.exe]"
	@echo "  build-all    - Build binaries for all supported platforms (Windows, macOS, Linux)"
	@echo "  build-windows - Build for Windows 64-bit"
	@echo "  build-mac-intel - Build for macOS Intel"
	@echo "  build-mac-arm - Build for macOS Apple Silicon (M1/M2)"
	@echo "  build-linux  - Build for Linux 64-bit"
	@echo "  build-linux-arm - Build for Linux ARM64"
	@echo "  run          - Run the application"
	@echo "  clean        - Remove build directory and all binaries"
	@echo "  help         - Display this help message"