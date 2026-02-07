package filesystem

import (
	"os"

	"github.com/dnnrly/layli/internal/usecases"
)

var _ usecases.FileReader = (*OSFileReader)(nil)
var _ usecases.FileWriter = (*OSFileWriter)(nil)

// OSFileReader reads files from the OS filesystem.
type OSFileReader struct{}

func NewOSFileReader() *OSFileReader { return &OSFileReader{} }

func (r *OSFileReader) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// OSFileWriter writes files to the OS filesystem.
type OSFileWriter struct{}

func NewOSFileWriter() *OSFileWriter { return &OSFileWriter{} }

func (w *OSFileWriter) Write(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
