package layli

import (
	"os"
	"strings"
	"testing"

	"github.com/dnnrly/layli/pathfinder/dijkstra"
)

// func FuzzConfig(f *testing.F) {
// 	dir, _ := os.ReadDir("./examples")
// 	for _, d := range dir {
// 		if !d.IsDir() && strings.HasSuffix(d.Name(), ".layli") {
// 			config, err := os.ReadFile("./examples/" + d.Name())
// 			if err != nil {
// 				panic(err)
// 			}
// 			f.Add(string(config)) // Use f.Add to provide a seed corpus
// 		}
// 	}

// 	f.Fuzz(func(t *testing.T, orig string) {
// 		_, _ = NewConfigFromFile(strings.NewReader(orig))
// 	})
// }

func FuzzDrawing(f *testing.F) {
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
		config, err := NewConfigFromFile(strings.NewReader(orig))
		if err != nil || config.Layout == "tarjan" {
			t.Skip()
		}

		layout, err := NewLayoutFromConfig(func(start, end dijkstra.Point) PathFinder {
			return dijkstra.NewPathFinder(start, end)
		}, config)

		if err != nil {
			t.Skip()
		}

		d := Diagram{
			Output:   func(data string) error { return nil },
			ShowGrid: false,
			Config:   *config,
			Layout:   layout,
		}

		_ = d.Draw()
	})
}
