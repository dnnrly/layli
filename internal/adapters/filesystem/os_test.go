package filesystem_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dnnrly/layli/internal/adapters/filesystem"
)

func TestOSFileWriter_and_Reader(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	data := []byte("hello world")

	writer := filesystem.NewOSFileWriter()
	if err := writer.Write(path, data); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	reader := filesystem.NewOSFileReader()
	got, err := reader.Read(path)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if string(got) != string(data) {
		t.Errorf("got %q, want %q", got, data)
	}
}

func TestOSFileReader_NonExistent(t *testing.T) {
	reader := filesystem.NewOSFileReader()
	_, err := reader.Read("/no/such/file/exists.txt")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestOSFileWriter_InvalidPath(t *testing.T) {
	writer := filesystem.NewOSFileWriter()
	err := writer.Write("/no/such/directory/file.txt", []byte("data"))
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestOSFileWriter_Permissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	writer := filesystem.NewOSFileWriter()
	if err := writer.Write(path, []byte("content")); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	perm := info.Mode().Perm()
	if perm != 0644 {
		t.Errorf("got permissions %o, want 0644", perm)
	}
}
