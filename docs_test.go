package layli_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadme_Updated(t *testing.T) {
	readme, err := ioutil.ReadFile("README.md")
	require.NoError(t, err)

	demo, err := ioutil.ReadFile("demo.layli")
	require.NoError(t, err)

	assert.Contains(t, string(readme), string(demo))
}
