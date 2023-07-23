//go:build docs

package layli_test

import (
	"io/ioutil"
	"os/exec"
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

func TestReadme_ImageUpToDate(t *testing.T) {
	err := exec.Command("./layli", "./demo.layli", "--show-grid").Run()
	require.NoError(t, err)

	result, err := exec.Command("git", "status", "demo.svg").Output()
	require.NoError(t, err)

	assert.NotContains(t, string(result), "demo.svg")
}
