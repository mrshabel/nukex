package main

// Config holds all application specific configurations
type Config struct {
	// pattern is the directory pattern to scan and remove
	Pattern string
	// MaxWorkers is the maximum number of workers to skip for parallel processing
	MaxWorkers int
	// SubDirsToSkip is the directories to skip during processing
	SubDirsToSkip []string
}

// New returns a new configuration object
func NewConfig(pattern string, maxWorkers int, subDirsToSkip []string) *Config {
	if pattern == "" {
		pattern = "node_modules"
	}
	if maxWorkers == 0 {
		maxWorkers = 20
	}
	if len(subDirsToSkip) == 0 {
		subDirsToSkip = []string{".", "..", ".venv", "venv", ".yarn", ".git"}
	}

	return &Config{Pattern: pattern, MaxWorkers: maxWorkers, SubDirsToSkip: subDirsToSkip}
}
