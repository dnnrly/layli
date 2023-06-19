package layli

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewDiagramFromFile_Simple(t *testing.T) {
	r := strings.NewReader(`
nodes:
  - id: node-1
    contents: "Some content here"
  - id: node-2
    contents: "More contents"
`)

	d, err := NewDiagramFromFile(io.NopCloser(r), func(output string) error {
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, len(d.config.Nodes))
}
