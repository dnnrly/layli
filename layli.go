package layli

import (
	"strings"

	svg "github.com/ajstarks/svgo"
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
