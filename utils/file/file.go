package file_utils

import (
	"os"
    "log"
)

func WriteToFile(s []byte, file string) {
	//
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	f.Write(s)
	defer f.Close()
}

// FileExists checks if a file exists and is not a directory.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}