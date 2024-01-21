package layli

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_validate(t *testing.T) {
	check := func(t *testing.T, config string, contained string) {
		_, err := NewConfigFromFile(strings.NewReader(config))

		assert.Error(t, err)
		assert.ErrorContains(t, err, contained)
	}

	t.Run("Edges have from and to", func(t *testing.T) {
		check(t, `nodes:
  - id: a
edges:
  - 0:`, "all edges must have a from and a to")
	})

	t.Run("Nodes require IDs", func(t *testing.T) {
		check(t, `nodes:
  - id: a
edges:
    - from: 0
      to: 00`, "all edges must have a from and a to that are valid node ids")
	})

	t.Run("Non-number for layout attempts", func(t *testing.T) {
		check(t, `nodes:
  - id: a
layout: random-shortest-square
layout-attempts: no-a-number`, "reading config file: yaml: unmarshal errors")
	})

	t.Run("Layout attempts too big", func(t *testing.T) {
		check(t, `nodes:
- id: a
layout: random-shortest-square
layout-attempts: 1e12`, "cannot specify more that 10000 layout attempts")
	})

	t.Run("Nodes require ID", func(t *testing.T) {
		check(t, `nodes:
    - 00: 0
    - id: b
    - id: c
    - id: d
    - id: e
edges:
    - from: a
      to: b
    - from: b
      to: c`, "all nodes must have an id")
	})

	t.Run("Edges cannot have same from and to", func(t *testing.T) {
		check(t, `nodes:
    - id: a
    - id: b
    - id: c
    - id: d
    - id: e
edges:
  - from: a
    to: b
  - from: a
    to: a
  - from: b
    to: c`, "edges cannot have the same from and to")
	})

	t.Run("Require at least 1 node", func(t *testing.T) {
		check(t, `path:
  strategy: random
  attempts: 99`, "must specify at least 1 node")
	})

	t.Run("Too many path attempt", func(t *testing.T) {
		check(t, `path:
  strategy: random
  attempts: 1e12
  nodes:
    - id: a`, "cannot specify more that 10000 path attempts")
	})

	t.Run("Margin too big", func(t *testing.T) {
		check(t, `margin: 20
nodes:
  - id: a`, "margin cannot be larger than 10")
	})
}
