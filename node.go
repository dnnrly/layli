package layli

import "fmt"

type NodeDrawer interface {
	Roundrect(x int, y int, w int, h int, rx int, ry int, s ...string)
	Textspan(x int, y int, t string, s ...string)
	TextEnd()
}

type Node struct {
	Id       string `yaml:"id"`
	Contents string `yaml:"contents"`

	X int `yaml:"-"`
	Y int `yaml:"-"`
}

func (n *Node) Draw(d NodeDrawer, spacing, width, Height int) {
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
