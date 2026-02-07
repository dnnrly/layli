//go:build fuzz

package main_test

import (
	"os"
	"strings"
	"testing"

	"github.com/antchfx/xmlquery"
	"github.com/dnnrly/layli/layout"
	"github.com/dnnrly/layli/pathfinder/dijkstra"
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
		_, _ = layout.NewConfigFromFile(strings.NewReader(orig))
	})
}

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
		config, err := layout.NewConfigFromFile(strings.NewReader(orig))
		if err != nil || config.Layout == "tarjan" {
			t.Skip()
		}

		l, err := layout.NewLayoutFromConfig(func(start, end dijkstra.Point) layout.PathFinder {
			return dijkstra.NewPathFinder(start, end)
		}, config)

		if err != nil {
			t.Skip()
		}

		d := layout.Diagram{
			Output:   func(data string) error { return nil },
			ShowGrid: false,
			Config:   *config,
			Layout:   l,
		}

		_ = d.Draw()
	})
}

func FuzzToAbsolute(f *testing.F) {
	dir, _ := os.ReadDir("./examples")
	for _, d := range dir {
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".svg") {
			config, err := os.ReadFile("./examples/" + d.Name())
			if err != nil {
				panic(err)
			}
			f.Add(string(config)) // Use f.Add to provide a seed corpus
		}
	}

	f.Fuzz(func(t *testing.T, orig string) {
		_, err := xmlquery.Parse(strings.NewReader(orig))
		if err != nil {
			t.Skip()
		}
		_ = layout.AbsoluteFromSVG(orig, func(string) error { return nil })
	})
}
