package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/charmbracelet/huh/spinner"
)

type DirInfo struct {
	Path          string
	Size          int64
	FormattedSize string
}

func main() {
	// print banner
	showBanner()

	if len(os.Args) < 2 {
		fmt.Println(errorStyle.Render("❌ Error: Please provide a directory to scan"))
		fmt.Println(dimStyle.Render("Usage: nukex <directory>"))
		os.Exit(1)
	}

	baseDir := os.Args[1]

	// validate directory exists
	info, err := os.Stat(baseDir)
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("❌ Directory '%s' does not exist", baseDir)))
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Println(errorStyle.Render("❌ Provided path is not a valid directory"))
		os.Exit(1)
	}

	// TODO: load config values from cli
	cfg := NewConfig("", 0, []string{})
	run(baseDir, cfg)
}

func run(baseDir string, cfg *Config) {
	var results []DirInfo
	// render scanning message while scanning for directories
	spin := spinner.New().Title(fmt.Sprintf("Scanning %s", baseDir)).Action(func() {
		results = scanDirectories(baseDir, cfg)
	},
	)
	if err := spin.Run(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("💯 Nukex completed scanning...")
	fmt.Println()
	// return early if no directories were found
	if len(results) == 0 {
		fmt.Println(warningStyle.Render("🤷 No node_modules directories found!"))
		fmt.Println(successStyle.Render("Your project is already clean! ✨"))
		return
	}

	// render results
	showResults(results)

	// select directory for deletion
	selectedDirs, err := selectDirectoriesToDelete(results)
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("❌ Error during selection: %v", err)))
		return
	}

	if len(selectedDirs) == 0 {
		fmt.Println(dimStyle.Render("👋 No directories selected. See you next time!"))
		return
	}

	// take user final confirmation
	confirmed, err := confirmDeletion(selectedDirs)
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("❌ Failed to proceed with deletion: %v", err)))
		return
	}

	if !confirmed {
		fmt.Println(warningStyle.Render("🚫 Deletion cancelled. You still have large node modules!"))
	}
	// proceed with deletion
	deleteDirectories(selectedDirs)

	// show final message
	showCompletionMessage(selectedDirs)
}

func scanDirectories(baseDir string, cfg *Config) []DirInfo {
	// get immediate children of base directory
	dirs, err := os.ReadDir(baseDir)
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("❌ Failed to access directory: %v", err)))
		os.Exit(1)
	}

	// store results
	var results []DirInfo
	var mu sync.Mutex

	// setup semaphore of default size 'max workers' to avoid running multiple goroutines at peak
	semaphore := make(chan struct{}, cfg.MaxWorkers)
	var wg sync.WaitGroup
	errCh := make(chan error, len(dirs))

	for _, dir := range dirs {
		// setup worker and acquire semaphore
		semaphore <- struct{}{}
		wg.Add(1)

		go func(dir os.DirEntry) {
			// release semaphore
			defer wg.Done()
			defer func() { <-semaphore }()

			// compose full path and walk through directory
			fullPath := filepath.Join(baseDir, dir.Name())
			err := filepath.WalkDir(fullPath, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				// skip unwanted directories
				if slices.Contains(cfg.SubDirsToSkip, filepath.Base(path)) {
					return filepath.SkipDir
				}

				if d.IsDir() && strings.Contains(filepath.Base(path), cfg.Pattern) {
					size, err := getDirSize(path)
					if err != nil {
						size = 0
					}

					// write results to slice in a thread-safe manner
					mu.Lock()
					results = append(results, DirInfo{Path: path, Size: size, FormattedSize: formatSize(size)})
					mu.Unlock()

					return filepath.SkipDir
				}

				return nil
			})

			if err != nil {
				errCh <- fmt.Errorf("error processing path %q: %v", dir.Name(), err)
			}
		}(dir)
	}

	// wait for workers to complete
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// handle errors
	for err := range errCh {
		fmt.Println(warningStyle.Render(fmt.Sprintf("⚠️  Warning: %s", err)))
	}
	return results
}
