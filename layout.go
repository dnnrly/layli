package layli

import (
	"fmt"
	"math"
)

type LayoutDrawer interface {
	Circle(x int, y int, r int, s ...string)
	Path(d string, s ...string)
	Roundrect(x int, y int, w int, h int, rx int, ry int, s ...string)
	Textspan(x int, y int, t string, s ...string)
	TextEnd()
}

type Layout struct {
	Nodes LayoutNodes
	Paths LayoutPaths

	nodeHeight   int // Height of a node in path unites
	nodeWidth    int // Width of a node in path units
	nodeMargin   int // Spare around nodes in path units
	layoutBorder int // Space at edge of layout in path units

	pathSpacing int // Length of a path unit in pixels

	vertexMap VertexMap // The state of all of the positions on the grid
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

	l := &Layout{
		Nodes: LayoutNodes{},

		nodeWidth:    c.NodeWidth,
		nodeHeight:   c.NodeHeight,
		nodeMargin:   c.Margin,
		layoutBorder: c.Border,
	}

	pos := 0
	for y := 0; y < size && pos < numNodes; y++ {
		for x := 0; x < size && pos < numNodes; x++ {
			l.Nodes = append(
				l.Nodes,
				NewLayoutNode(
					c.Nodes[pos].Id,
					c.Nodes[pos].Contents,
					l.layoutBorder+
						l.nodeMargin+
						(x*l.nodeWidth)+
						(x*(l.nodeMargin*2)),
					l.layoutBorder+
						l.nodeMargin+
						(y*l.nodeHeight)+
						(y*(l.nodeMargin*2)),
					l.nodeWidth, l.nodeHeight,
				),
			)

			pos++
		}
	}

	for _, p := range c.Edges {
		l.AddPath(p.From, p.To)
	}

	return l
}

// LayoutHeight is the height in path units
func (l *Layout) LayoutHeight() int {
	numNodes := len(l.Nodes)

	root := math.Sqrt(float64(numNodes))
	rows := numNodes / int(root)
	if numNodes%int(root) == 0 {
		rows--
	}

	return l.layoutBorder*2 +
		(rows * l.nodeHeight) +
		(rows * l.nodeMargin * 2)
}

// LayoutWidth is the width in path units
func (l *Layout) LayoutWidth() int {
	numNodes := len(l.Nodes)
	columns := numNodes

	if numNodes >= 4 {
		root := math.Sqrt(float64(numNodes))
		columns = int(math.Ceil(root))
	}

	return l.layoutBorder*2 +
		(columns * l.nodeWidth) +
		(columns * l.nodeMargin * 2)
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
	for x := 0; x <= l.LayoutWidth(); x++ {
		for y := 0; y <= l.LayoutHeight(); y++ {
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

	for _, p := range l.Paths {
		p.Draw(canvas, spacing)
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

func (nodes LayoutNodes) ByID(id string) *LayoutNode {
	for _, n := range nodes {
		if n.Id == id {
			return &n
		}
	}
	return nil
}

func NewLayoutNode(id, contents string, left, top, width, height int) LayoutNode {
	return LayoutNode{
		Id:       id,
		Contents: contents,

		width:  width,
		height: height,

		top:    top,
		bottom: top + height - 1,
		left:   left,
		right:  left + width - 1,
	}
}

func (n *LayoutNode) IsInside(x, y int) bool {
	if y >= n.top && y <= n.bottom && x >= n.left && x <= n.right {
		return true
	}

	return false
}

func (n *LayoutNode) IsPort(x, y int) bool {
	if y == n.top && x > n.left && x < n.right {
		return true
	}

	if y == n.bottom && x > n.left && x < n.right {
		return true
	}

	if x == n.left && y > n.top && y < n.bottom {
		return true
	}

	if x == n.right && y > n.top && y < n.bottom {
		return true
	}

	return false
}

func (n *LayoutNode) GetPorts() Points {
	ports := Points{}

	for i := n.left + 1; i < n.right; i++ {
		ports = append(ports,
			Point{X: float64(i), Y: float64(n.top)},
			Point{X: float64(i), Y: float64(n.bottom)},
		)
	}

	for i := n.top + 1; i < n.bottom; i++ {
		ports = append(ports,
			Point{X: float64(n.left), Y: float64(i)},
			Point{X: float64(n.right), Y: float64(i)},
		)
	}

	return ports
}

func (n *LayoutNode) GetCentre() Point {
	return Point{
		X: float64(n.left) + float64(n.width)/2,
		Y: float64(n.top) + float64(n.height)/2,
	}
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

type LayoutPath struct {
	points Points
}

func (p *LayoutPath) Draw(canvas LayoutDrawer, spacing int) {
	canvas.Path(p.points.Path(spacing), `class="path-line"`)
}

type LayoutPaths []LayoutPath

func (paths *LayoutPaths) Draw(canvas LayoutDrawer, spacing int) {
	fmt.Printf("Drawing %d paths\n", len(*paths))
	for i, p := range *paths {
		fmt.Printf("Drawing path %d of %d\n", i+1, len(*paths))
		p.Draw(canvas, spacing)
	}
}