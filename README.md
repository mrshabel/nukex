# NUKEX🧹

A fast and interactive CLI tool to find and clean up node_modules directories that are consuming disk space.

## Why NUKEX?

-   🔍 Quickly find all node_modules directories eating up disk space
-   📏 Get accurate size of each directory
-   🚀 Process large directory trees efficiently with parallel processing
-   ⚡ Smart directory skipping (.git, .venv, .yarn)
-   🎯 Interactive selection for which directories to remove
-   ✨ Safe deletion with confirmation prompts

## Installation

### Prerequisites

-   Go: 1.23 or higher

### Build from Source

1. Clone the repository:

```bash
git clone github.com/mrshabel/nukex
cd nukex
```

2. Build the project:

```bash
make build
```

This will create a `bin/` directory with the binary for your current platform.

### Pre-built Binaries

Download the appropriate binary from the releases page:

-   **Windows**: `nukex-windows-amd64.exe`
-   **macOS Intel**: `nukex-darwin-amd64`
-   **macOS Apple Silicon**: `nukex-darwin-arm64`
-   **Linux**: `nukex-linux-amd64`
-   **Linux ARM64**: `nukex-linux-arm64`

### Cross-Platform Building

To build for all supported platforms:

```bash
make build-all
```

Or build for a specific platform:

```bash
# windows 64-bit
make build-windows
# macOS intel
make build-mac-intel
# macOS apple silicon
make build-mac-arm
# linux 64-bit
make build-linux
# linux arm64
make build-linux-arm
```

## Usage

```bash
nukex <directory>
```

Example:

```bash
nukex ~/projects


💫 Scanning ~/projects...
💯 Nukex completed scanning...

Found node_modules directories:
📁 ~/projects/try/node_modules (2.56 GB)
📁 ~/projects/test/node_modules (156.78 MB)

? Select directories to clean up:
  [x] ~/projects/my-app/node_modules
  [ ] ~/projects/another-app/node_modules

? Are you sure you want to delete these directories? [y/N]
```

## Features

-   🏷️ Controlled parallel processing with semaphores
-   🚀 Fast parallel directory search
-   📊 Shows directory sizes in human readable format
-   🎯 Targets only node_modules directories

## Contributing

Pull requests are welcome. For major changes, please open an issue first.

## TODO

-   Add advanced processing flags with viper
-   Add size-based filtering options
