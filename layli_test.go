package layli

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigFromFile(t *testing.T) {
	r := strings.NewReader(`
nodes:
  - id: node-1
    contents: "C1"
  - id: node-2
    contents: "C2"
`)

	config, err := NewConfigFromFile(r)
	require.NoError(t, err)
	assert.Equal(t, Config{
		Path: ConfigPath{
			Attempts: 5,
		},
		Nodes: ConfigNodes{
			ConfigNode{
				Id:       "node-1",
				Contents: "C1",
			},
			ConfigNode{
				Id:       "node-2",
				Contents: "C2",
			},
		},
		Spacing:    20,
		Border:     1,
		Margin:     2,
		NodeWidth:  5,
		NodeHeight: 3,
	}, *config)
}

func TestNewConfigFromFile_FailsOnBadYaml(t *testing.T) {
	r := strings.NewReader(`
nodes:
  - id: node-1
-
  `)

	_, err := NewConfigFromFile(r)
	require.Error(t, err)
}
