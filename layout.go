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
			l.Nodes = append(
				l.Nodes,
				NewLayoutNode(
					c.Nodes[pos].Id,
					c.Nodes[pos].Contents,
					(x+1)*l.gridSpacing,
					(y+1)*l.gridSpacing,
					pathSpacing,
					l.nodeWidth, l.nodeHeight,
				),
			)

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

	// location of the centre of the node in pixels
	X int
	Y int

	spacing int // number of pixels between 'path' points on either axis

	// dimensions of the node in pixels
	width  int
	height int

	// the pixel position of the edges of the node
	top    int
	bottom int
	left   int
	right  int
}

type LayoutNodes []LayoutNode

func NewLayoutNode(id, contents string, x, y, spacing, width, height int) LayoutNode {
	return LayoutNode{
		Id:       id,
		Contents: contents,
		X:        x,
		Y:        y,

		spacing: spacing,

		width:  width,
		height: height,

		top:    y - height/2,
		bottom: y + height/2,
		left:   x - width/2,
		right:  x + width/2,
	}
}

func (n *LayoutNode) Draw(d LayoutDrawer, spacing, width, height int) {
	d.Roundrect(
		n.left, n.top,
		width, height,
		3, 3,
		fmt.Sprintf(`id="%s"`, n.Id),
	)
	d.Textspan(
		n.X,
		n.Y,
		n.Contents,
		fmt.Sprintf(`id="%s-text"`, n.Id),
		"font-size:10px",
	)
	d.TextEnd()
}
