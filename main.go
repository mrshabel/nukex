package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

const (
	pattern    = "node_modules"
	maxWorkers = 20
)

var (
	subDirsToSkip = []string{".", "..", ".venv", "venv", ".yarn", ".git"}
)

type DirInfo struct {
	Path string
	Size string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Base directory is required")
		os.Exit(1)
	}

	// validate directory
	baseDir := string(os.Args[1])
	info, err := os.Stat(baseDir)
	if err != nil {
		fmt.Printf("Failed to access provided directory: %v\n", err)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Println("Provided path is not a valid directory")
		os.Exit(1)
	}

	// get all direct children of base dir
	dirs, err := os.ReadDir(baseDir)
	if err != nil {
		fmt.Printf("Failed to access sub-directory(s): %v\n", err)
		os.Exit(1)
	}

	// process individual directories in parallel
	resultCh := make(chan DirInfo, len(dirs))
	errCh := make(chan error, len(dirs))
	// setup semaphore of default size 'max workers' to avoid running multiple goroutines at peak
	semaphore := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	for _, dir := range dirs {
		// setup worker and acquire semaphore
		semaphore <- struct{}{}
		wg.Add(1)

		go func(dir os.DirEntry) {
			// release semaphore
			defer wg.Done()
			defer func() {
				<-semaphore
			}()

			// compose full path
			fullPath := filepath.Join(baseDir, dir.Name())
			// walk through directory
			err = filepath.WalkDir(fullPath, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					fmt.Printf("Error accessing path: %v\n", path)
					return err
				}

				// skip "." and ".." dirs
				if slices.Contains(subDirsToSkip, filepath.Base(path)) {
					// fmt.Printf("skipping special directories: %+v \n", path)
					return filepath.SkipDir
				}

				// add found folder to results channel and skip processing it
				if d.IsDir() && strings.Contains(filepath.Base(path), pattern) {
					// calculate directory size
					size, err := getDirSize(path)
					if err != nil {
						size = "unknown"
					}
					resultCh <- DirInfo{Path: path, Size: size}
					return filepath.SkipDir
				}

				return nil
			})

			// report error
			if err != nil {
				errCh <- fmt.Errorf("error processing path %q: %v\n", dir.Name(), err)
			}
		}(dir)
	}

	// done channel to coordinate goroutines
	done := make(chan struct{})
	go func() {
		// wait for all workers to complete
		wg.Wait()

		close(semaphore)
		close(resultCh)
		close(errCh)

		// signal completion
		close(done)
	}()

	for {
		select {
		case result, ok := <-resultCh:
			if !ok {
				continue
			}
			fmt.Printf("Found node_modules at: %s\nSize: %v\n\n", result.Path, result.Size)
		case err, ok := <-errCh:
			if !ok {
				continue
			}
			fmt.Printf("Error: %s\n", err)
		case <-done:
			return
		}
	}
}

// getDirSize returns the directory size formatted in exponents of bytes
func getDirSize(path string) (string, error) {
	// calculate directory size
	var size float64
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// get file size
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			size += float64(info.Size())
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	// format size
	units := []string{"", "K", "M", "G", "T"}
	ptr := 0
	for size >= 1024 && ptr < len(units) {
		ptr++
		// compute size in base form
		size /= 1024
	}

	formatted := fmt.Sprintf("%.2f %sB", size, units[ptr])

	return formatted, nil
}
