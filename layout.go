package layli

import (
	"errors"
	"fmt"
	"strings"
)

type LayoutDrawer interface {
	Circle(x int, y int, r int, s ...string)
	Path(d string, s ...string)
	Roundrect(x int, y int, w int, h int, rx int, ry int, s ...string)
	Textspan(x int, y int, t string, s ...string)
	TextEnd()
}

type Layout struct {
	Nodes        LayoutNodes
	Paths        LayoutPaths
	CreateFinder CreateFinder

	nodeHeight   int // Height of a node in path unites
	nodeWidth    int // Width of a node in path units
	nodeMargin   int // Spare around nodes in path units
	layoutBorder int // Space at edge of layout in path units

	pathSpacing int // Length of a path unit in pixels
}

func NewLayoutFromConfig(finder CreateFinder, c *Config) (*Layout, error) {
	arranger, err := selectArrangement(c)
	if err != nil {
		return nil, err
	}

	pathStrategy, err := selectPathStrategy(c)
	if err != nil {
		return nil, err
	}

	l := &Layout{
		Nodes:        arranger(c),
		CreateFinder: finder,

		nodeWidth:    c.NodeWidth,
		nodeHeight:   c.NodeHeight,
		nodeMargin:   c.Margin,
		layoutBorder: c.Border,
	}

	err = pathStrategy(*c, &l.Paths, l.FindPath)
	if err != nil {
		return nil, err
	}

	return l, nil
}

// LayoutHeight is the height in path units
func (l *Layout) LayoutHeight() int {
	maxBottom := 0
	for _, n := range l.Nodes {
		if maxBottom <= n.bottom {
			maxBottom = n.bottom
		}
	}

	return maxBottom + l.nodeMargin + 1 + l.layoutBorder
}

// LayoutWidth is the width in path units
func (l *Layout) LayoutWidth() int {
	maxRight := 0
	for _, n := range l.Nodes {
		if maxRight <= n.right {
			maxRight = n.right
		}
	}

	return maxRight + l.nodeMargin + 1 + l.layoutBorder
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
	vm := BuildVertexMap(l)
	for _, v := range vm.GetVertexPoints() {
		canvas.Circle(int(v.X)*spacing, int(v.Y)*spacing, 1, `class="path-dot"`)
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

func (nodes LayoutNodes) String() string {
	var buf []string
	for _, n := range nodes {
		buf = append(buf, n.Id)
	}
	return fmt.Sprintf("[%s]", strings.Join(buf, ", "))
}

func (nodes LayoutNodes) ByID(id string) *LayoutNode {
	for _, n := range nodes {
		if n.Id == id {
			return &n
		}
	}
	return nil
}

func (n LayoutNodes) ConnectionDistances(connections ConfigEdges) (float64, error) {
	dist := 0.0

	for _, c := range connections {
		f := n.ByID(c.From)
		if f == nil {
			return 0, errors.New("cannot find node " + c.From)
		}

		to := n.ByID(c.To)
		if to == nil {
			return 0, errors.New("cannot find node " + c.To)
		}

		dist += f.GetCentre().Distance(to.GetCentre())
	}

	return dist, nil
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
		(n.width-1)*spacing, (n.height-1)*spacing,
		3, 3,
		fmt.Sprintf(`id="%s"`, n.Id),
	)
	d.Textspan(
		n.left*spacing+(((n.width-1)*spacing)/2),
		n.top*spacing+(((n.height-1)*spacing)/2),
		n.Contents,
		fmt.Sprintf(`id="%s-text"`, n.Id),
		"font-size:10px",
	)
	d.TextEnd()
}

type LayoutPath struct {
	Points Points
}

func (p *LayoutPath) Draw(canvas LayoutDrawer, spacing int) {
	canvas.Path(p.Points.Path(spacing), `class="path-line"`, `marker-end="url(#arrow)"`)
}

func (paths *LayoutPath) Length() float64 {
	if len(paths.Points) <= 1 {
		return 0.0
	}

	length := 0.0
	for i := 1; i < len(paths.Points); i++ {
		a := paths.Points[i-1]
		b := paths.Points[i]

		length += a.Distance(b)
	}

	return length
}

type LayoutPaths []LayoutPath

func (paths *LayoutPaths) Draw(canvas LayoutDrawer, spacing int) {
	for _, p := range *paths {
		p.Draw(canvas, spacing)
	}
}

func (paths LayoutPaths) Length() float64 {
	totalLength := float64(0)
	for _, p := range paths {
		totalLength += p.Length()
	}
	return totalLength
}
