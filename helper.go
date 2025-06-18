package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// formatSize returns a formatted version of the input bytes into exponents of bytes
func formatSize(size int64) string {
	totalSize := float64(size)
	units := []string{"B", "KB", "MB", "GB", "TB"}
	ptr := 0

	for totalSize >= 1024 && ptr < len(units)-1 {
		ptr++
		totalSize /= 1024
	}

	// use whole number for bytes
	if ptr == 0 {
		return fmt.Sprintf("%.0f %s", totalSize, units[ptr])
	}
	return fmt.Sprintf("%.2f %s", totalSize, units[ptr])
}

// getDirSize calculates the total size of all contents in a directory recursively
func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			size += info.Size()
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return size, nil
}

// shortenPath returns a shortened relative path of the given input path
func shortenPath(path string) string {
	const maxLength = 50

	// get relative path from current directory
	if cwd, err := os.Getwd(); err == nil {
		if rel, err := filepath.Rel(cwd, path); err == nil && len(rel) < len(path) {
			path = rel
		}
	}

	if len(path) <= maxLength {
		return path
	}

	// truncate from the middle
	start := path[:15]
	end := path[len(path)-30:]
	return start + "..." + end
}

// calculateTotalSize returns the formatted cumulative sum of the total directories
func calculateTotalSize(dirs []DirInfo) string {
	var totalSize int64
	for _, dir := range dirs {
		totalSize += dir.Size
	}
	return formatSize(totalSize)
}
