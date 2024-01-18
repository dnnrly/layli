package layli

import (
	"os"
	"strings"
	"testing"
)

func FuzzConfig(f *testing.F) {
	dir, _ := os.ReadDir("./examples")
	for _, d := range dir {
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".layli") {
			config, err := os.ReadFile("./examples/" + d.Name())
			if err != nil {
				panic(err)
			}
			f.Add(string(config)) // Use f.Add to provide a seed corpus
		}
	}

	f.Fuzz(func(t *testing.T, orig string) {
		_, _ = NewConfigFromFile(strings.NewReader(orig))
	})
}
