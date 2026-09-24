package file_utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExistsAndWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "f.txt")

	if FileExists(path) {
		t.Error("missing file reported as existing")
	}
	if err := WriteToFile([]byte("hi"), path); err != nil {
		t.Fatal(err)
	}
	if !FileExists(path) {
		t.Error("written file not found")
	}
	if FileExists(dir) {
		t.Error("directory reported as file")
	}
	if err := WriteToFile([]byte("x"), filepath.Join(dir, "missing", "f.txt")); err == nil {
		t.Error("expected error writing into a missing directory")
	}
}

func TestFileExistsPermissionDeniedDoesNotPanic(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores directory permissions")
	}
	dir := t.TempDir()
	locked := filepath.Join(dir, "locked")
	if err := os.Mkdir(locked, 0o000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(locked, 0o755)

	if FileExists(filepath.Join(locked, "f.txt")) {
		t.Error("unreadable path reported as existing")
	}
}
