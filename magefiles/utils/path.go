package utils

import (
	"fmt"
	"os"
	"strings"
)

// RepoRoot Returns the path to the repository root.
func RepoRoot() (string, error) {
	cur, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}

	return cur, nil
}

// Exists checks if a file exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// PathExists checks if a path exists.
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// RemoveDebugBinFiles removes all __debug_bin files.
func RemoveDebugBinFiles(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("read dir: %w", err)
	}

	for _, file := range files {
		absPath := path + "/" + file.Name()
		if file.IsDir() {
			if err := RemoveDebugBinFiles(absPath); err != nil {
				return fmt.Errorf("remove debug bin files: %w", err)
			}
		}

		if file.IsDir() {
			continue
		}

		if !Exists(absPath) || !isDebugBinFile(absPath) {
			continue
		}

		if err := os.Remove(absPath); err != nil {
			return fmt.Errorf("remove debug bin files: %w", err)
		}
	}

	return nil
}

func isDebugBinFile(path string) bool {
	filename := path[strings.LastIndex(path, "/")+1:]
	binI := len("__debug_bin")
	if len(filename) < binI {
		return false
	}
	return filename[:binI] == "__debug_bin"
}
