package layli

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_validate(t *testing.T) {
	check := func(config string, contained string) {
		_, err := NewConfigFromFile(strings.NewReader(config))

		assert.Error(t, err)
		assert.ErrorContains(t, err, contained)
	}

	check(`edges:
- 0:`, "all edges must have a from and a to")
	check(`edges:
    - from: 0
      to: 00`, "all edges must have a from and a to that are valid node ids")
	check(`path:
  strategy: random
  attempts: 1e12a`, "reading config file: yaml: unmarshal errors")
	check(`layout: random-shortest-square
layout-attempts: 1e12n`, "reading config file: yaml: unmarshal errors")
	check(`nodes:
    - 00: 0
    - id: b
    - id: c
    - id: d
    - id: e
edges:
    - from: b
      to: b
    - from: b
      to: b
    - from: b
      to: b
    - from: b
      to: b
    - from: b
      to: b`, "all nodes must have an id")
}
