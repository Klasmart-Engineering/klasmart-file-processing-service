package test

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"testing"
)

func Checksum(t *testing.T, filePath, sha string) {
	f, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		t.Fatal(err)
	}

	if sha != fmt.Sprintf("%x", h.Sum(nil)) {
		t.Logf("sha is %s expected is %s", fmt.Sprintf("%x", h.Sum(nil)), sha)
		t.Fail()
	}
}

func RemoveFile(t *testing.T, filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		t.Fatal(err)
	}
}
