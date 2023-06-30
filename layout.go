package layli

import (
	"fmt"
	"math"
)

type LayoutDrawer interface {
	Circle(x int, y int, r int, s ...string)
	Roundrect(x int, y int, w int, h int, rx int, ry int, s ...string)
	Textspan(x int, y int, t string, s ...string)
	TextEnd()
}

type Layout struct {
	Nodes      LayoutNodes
	GridHeight int
	GridWidth  int

	pathSpacing int
	gridSpacing int
	nodeHeight  int
	nodeWidth   int
}

func NewLayoutFromConfig(c Config) *Layout {
	numNodes := len(c.Nodes)

	root := math.Sqrt(float64(numNodes))
	size := int(math.Round(root))

	if size < int(root) {
		size++
	}
	if numNodes < 4 {
		size = 2
	}
	pathSpacing := 20

	l := &Layout{
		Nodes:      LayoutNodes{},
		GridWidth:  size,
		GridHeight: size,

		pathSpacing: pathSpacing,
		gridSpacing: pathSpacing * 7,
		nodeWidth:   pathSpacing * 5,
		nodeHeight:  pathSpacing * 3,
	}

	pos := 0
	for y := 0; y < size && pos < numNodes; y++ {
		for x := 0; x < size && pos < numNodes; x++ {
			l.Nodes = append(l.Nodes, LayoutNode{
				Id:       c.Nodes[pos].Id,
				Contents: c.Nodes[pos].Contents,
				X:        x + 1,
				Y:        y + 1,
			})

			pos++
		}
	}

	return l
}

func (l *Layout) LayoutHeight() int {
	return l.gridSpacing * (l.GridHeight + 1)
}

func (l *Layout) LayoutWidth() int {
	return l.gridSpacing * (l.GridWidth + 1)
}

func (l *Layout) ShowGrid(canvas LayoutDrawer) {
	for y := 0; y < l.gridSpacing*(l.GridHeight+1); y += l.pathSpacing {
		for x := 0; x < l.gridSpacing*(l.GridWidth+1); x += l.pathSpacing {
			canvas.Circle(
				l.pathSpacing/2+x,
				l.pathSpacing/2+y,
				1,
				`class="path-dot"`,
			)
		}
	}
}

func (l *Layout) Draw(canvas LayoutDrawer) {
	for _, n := range l.Nodes {
		n.Draw(canvas, l.gridSpacing, l.nodeWidth, l.nodeHeight)
	}
}

type LayoutNode struct {
	Id       string
	Contents string
	X        int
	Y        int
}

type LayoutNodes []LayoutNode

func (n *LayoutNode) Draw(d LayoutDrawer, spacing, width, Height int) {
	d.Roundrect(
		spacing*n.X-width/2, spacing*n.Y-Height/2,
		width, Height,
		3, 3,
		fmt.Sprintf(`id="%s"`, n.Id),
	)
	d.Textspan(
		spacing*n.X,
		spacing*n.Y,
		n.Contents,
		fmt.Sprintf(`id="%s-text"`, n.Id),
		"font-size:10px",
	)
	d.TextEnd()
}
