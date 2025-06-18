# NUKEX🧹

A fast and interactive CLI tool to find and clean up node_modules directories that are consuming disk space.

## Why NUKEX?

-   🔍 Quickly find all node_modules directories eating up disk space
-   📏 Get accurate size of each directory
-   🚀 Process large directory trees efficiently with parallel processing
-   ⚡ Smart directory skipping (.git, .venv, .yarn)
-   🎯 Interactive selection for which directories to remove
-   ✨ Safe deletion with confirmation prompts

## Setup

Build from source:

```bash
make build
```

or use pre-built binary

```bash
./nukex.exe
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
