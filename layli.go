package layli

import (
	"fmt"
	"strconv"
	"strings"

	svg "github.com/ajstarks/svgo"
	"github.com/antchfx/xmlquery"
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
	dom, err := xmlquery.Parse(strings.NewReader(svg))
	if err != nil {
		return fmt.Errorf("parsing svg: %w", err)
	}

	config := &Config{
		Nodes: ConfigNodes{},
	}

	for _, n := range xmlquery.Find(dom, "//rect") {
		id := n.SelectAttr("id")
		x, err := strconv.Atoi(n.SelectAttr("data-pos-x"))
		if err != nil {
			return fmt.Errorf("parsing X: %w", err)
		}
		y, err := strconv.Atoi(n.SelectAttr("data-pos-y"))
		if err != nil {
			return fmt.Errorf("parsing Y: %w", err)
		}

		text := xmlquery.FindOne(dom, "//*[@id='"+id+"-text']")
		if text == nil {
			return fmt.Errorf("no text found for node %s", id)
		}

		config.Nodes = append(config.Nodes, ConfigNode{
			Id:       id,
			Contents: text.InnerText(),
			Position: Position{X: x, Y: y},
		})
	}

	for _, e := range xmlquery.Find(dom, "//g/path") {
		config.Edges = append(config.Edges, ConfigEdge{
			From: e.SelectAttr("data-from"),
			To:   e.SelectAttr("data-to"),
		})
	}

	return output(config.String())
}
