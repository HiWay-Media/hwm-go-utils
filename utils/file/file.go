package file_utils

import (
	"os"
)

// WriteToFile creates or truncates file and writes s to it.
func WriteToFile(s []byte, file string) error {
	return os.WriteFile(file, s, 0666)
}

// FileExists checks if a file exists and is not a directory.
// It returns false when the file cannot be stat'ed (e.g. permission denied).
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
