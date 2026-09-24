package db

import (
	"strings"
	"testing"
)

func TestOpenErrorDoesNotLeakPassword(t *testing.T) {
	// nothing listens on port 1: the connection is refused immediately
	_, err := Open("app", "s3cr3t-pass", "127.0.0.1", 1, "shop", "", "")
	if err == nil {
		t.Fatal("expected a connection error")
	}
	if strings.Contains(err.Error(), "s3cr3t-pass") {
		t.Errorf("error leaks the password: %v", err)
	}
	if !strings.Contains(err.Error(), "app@127.0.0.1:1/shop") {
		t.Errorf("error should identify the target: %v", err)
	}
}
