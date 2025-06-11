# NUKEX🧹

Simple CLI tool to search and list node_modules directories with their sizes.

## Why NUKEX?

-   🔍 Quickly find all node_modules directories eating up disk space
-   📏 Get accurate size measurements of each directory
-   🚀 Process large directory trees efficiently with parallel search

## Setup

Build from source:

```bash
make build
```

or use pre-built binary from the `bin` directory

```bash
./bin/nukex
```

## Usage

```bash
nukex <directory>
```

Example:

```bash
nukex ~/projects


Found node_modules at: ~/projects/my-app/node_modules
Size: 2.56 GB

Found node_modules at: ~/projects/another-app/node_modules
Size: 156.78 MB
```

## Features

-   🏷️ Controlled parallel processing with semaphores
-   🚀 Fast parallel directory search
-   📊 Shows directory sizes in human readable format
-   🎯 Targets only node_modules directories

## Contributing

Pull requests are welcome. For major changes, please open an issue first.
