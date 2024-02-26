package layli

import (
	"strings"

	svg "github.com/ajstarks/svgo"
	"gopkg.in/yaml.v3"
)

type OutputFunc func(output string) error

type Diagram struct {
	Output   OutputFunc
	Config   Config
	Layout   *Layout
	ShowGrid bool
}

// Draw turns the diagram in to an image
func (d *Diagram) Draw() error {
	w := strings.Builder{}

	canvas := svg.New(&w)
	canvas.Start(
		(d.Layout.LayoutWidth()-1)*d.Config.Spacing,
		(d.Layout.LayoutHeight()-1)*d.Config.Spacing,
		"style=\"background-color: white;\"",
	)
	if len(d.Config.Styles) != 0 {
		canvas.Style("text/css", d.Config.Styles.toCSS())
	}
	canvas.Gstyle("text-anchor:middle;font-family:sans;fill:none;stroke:black")
	canvas.Def()
	canvas.Marker("arrow", 10, 5, 7, 7,
		`viewBox="0 0 10 10"`,
		`fill="black"`,
		`orient="auto-start-reverse"`)
	canvas.Path("M 0 0 L 10 5 L 0 10 z")
	canvas.MarkerEnd()
	canvas.DefEnd()

	if d.ShowGrid {
		d.Layout.ShowGrid(canvas, d.Config.Spacing)
	}

	d.Layout.Draw(canvas, d.Config.Spacing)

	canvas.Gend()

	canvas.End()
	return d.Output(w.String())
}

// AbsuluteFromSVG parses a string of an SVG and turns it in to a Layli configuration
// with with absulute layout that can represent the same SVG
func AbsoluteFromSVG(svg string, output OutputFunc) error {
	config := &Config{
		Nodes: ConfigNodes{
			ConfigNode{Id: "1", Position: Position{X: 3, Y: 3}},
			ConfigNode{Id: "2", Position: Position{X: 3, Y: 7}},
		},
	}

	buf := strings.Builder{}
	yaml.NewEncoder(&buf).Encode(config)

	return output(buf.String())
}
