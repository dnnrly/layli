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
	Nodes LayoutNodes

	nodeHeight int // Height of a node in path unites
	nodeWidth  int // Width of a node in path units

	pathSpacing int // Length of a path unit in pixels
}

func NewLayoutFromConfig(c Config) *Layout {
	numNodes := len(c.Nodes)

	root := math.Sqrt(float64(numNodes))
	size := int(math.Ceil(root))

	if size < int(root) {
		size++
	}
	if numNodes < 4 {
		size = 2
	}
	pathSpacing := 20

	l := &Layout{
		Nodes: LayoutNodes{},

		pathSpacing: pathSpacing,
		nodeWidth:   5,
		nodeHeight:  3,
	}

	pos := 0
	for y := 0; y < size && pos < numNodes; y++ {
		for x := 0; x < size && pos < numNodes; x++ {
			l.Nodes = append(
				l.Nodes,
				NewLayoutNode(
					c.Nodes[pos].Id,
					c.Nodes[pos].Contents,
					x*(l.nodeWidth+1),
					y*(l.nodeHeight+1),
					l.nodeWidth, l.nodeHeight,
				),
			)

			pos++
		}
	}

	return l
}

func (l *Layout) LayoutHeight() int {
	numNodes := len(l.Nodes)

	root := math.Sqrt(float64(numNodes))
	rows := numNodes / int(root)
	if numNodes%int(root) == 0 {
		rows--
	}

	return (1 + (rows * l.nodeHeight) + rows) * l.pathSpacing
}

func (l *Layout) LayoutWidth() int {
	numNodes := len(l.Nodes)
	columns := numNodes

	if numNodes >= 4 {
		root := math.Sqrt(float64(numNodes))
		columns = int(math.Ceil(root))
	}

	return (1 + (columns * l.nodeWidth) + columns) * l.pathSpacing
}

func (l *Layout) InsideAny(x, y int) bool {
	for _, n := range l.Nodes {
		if n.IsInside(x, y) {
			return true
		}
	}
	return false
}

func (l *Layout) IsAnyPort(x, y int) bool {
	for _, n := range l.Nodes {
		if n.IsPort(x, y) {
			return true
		}
	}
	return false
}

func (l *Layout) ShowGrid(canvas LayoutDrawer, spacing int) {
	for y := 0; y <= l.LayoutHeight(); y++ {
		for x := 0; x <= l.LayoutWidth(); x++ {
			gridX := x
			gridY := y

			if !l.InsideAny(gridX, gridY) || l.IsAnyPort(gridX, gridY) {
				canvas.Circle(gridX*spacing, gridY*spacing, 1, `class="path-dot"`)
			}
		}
	}
}

func (l *Layout) Draw(canvas LayoutDrawer, spacing int) {
	for _, n := range l.Nodes {
		n.Draw(canvas, spacing)
	}
}

type LayoutNode struct {
	Id       string
	Contents string

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

func NewLayoutNode(id, contents string, left, top, width, height int) LayoutNode {
	return LayoutNode{
		Id:       id,
		Contents: contents,

		width:  width,
		height: height,

		top:    top,
		bottom: top + height,
		left:   left,
		right:  left + width,
	}
}

func (n *LayoutNode) IsInside(x, y int) bool {
	if y >= n.top && y <= n.bottom && x >= n.left && x <= n.right {
		return true
	}

	return false
}

func (n *LayoutNode) IsPort(x, y int) bool {
	// Points not on any of the borders are rejected
	if x != n.left && x != n.right && y != n.top && y != n.bottom {
		return false
	}

	// All corners are rejected
	if x == n.left && y == n.top ||
		x == n.left && y == n.bottom ||
		x == n.right && y == n.top ||
		x == n.right && y == n.bottom {
		return false
	}

	if y < n.top && y > n.bottom {
		return false
	}

	if x < n.left && x > n.right {
		return true
	}

	return true
}

func (n *LayoutNode) Draw(d LayoutDrawer, spacing int) {
	d.Roundrect(
		n.left*spacing, n.top*spacing,
		n.width*spacing, n.height*spacing,
		3, 3,
		fmt.Sprintf(`id="%s"`, n.Id),
	)
	d.Textspan(
		n.left*spacing+((n.width*spacing)/2),
		n.top*spacing+((n.height*spacing)/2),
		n.Contents,
		fmt.Sprintf(`id="%s-text"`, n.Id),
		"font-size:10px",
	)
	d.TextEnd()
}
