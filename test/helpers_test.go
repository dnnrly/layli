//go:build acceptance

package test_test

import (
	"strconv"
	"strings"

	"github.com/antchfx/xmlquery"
)

// getNodeIds takes an *xmlquery.Node and returns the "id" attribute of all elements
func getNodeIds(node *xmlquery.Node) []string {
	var ids []string

	// Iterate over child nodes
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == xmlquery.ElementNode && child.Data == "rect" {
			// Check if the element has an "id" attribute
			for _, attr := range child.Attr {
				if attr.Name.Local == "id" {
					ids = append(ids, attr.Value)
					break
				}
			}
		}

		// Recursively call getNodeIds on child nodes
		childIds := getNodeIds(child)
		ids = append(ids, childIds...)
	}

	return ids
}

// Helper function to check if two rectangles overlap
func isOverlap(rectA, rectB *xmlquery.Node) bool {
	xA := rectA.SelectAttr("x")
	yA := rectA.SelectAttr("y")
	widthA := rectA.SelectAttr("width")
	heightA := rectA.SelectAttr("height")

	xB := rectB.SelectAttr("x")
	yB := rectB.SelectAttr("y")
	widthB := rectB.SelectAttr("width")
	heightB := rectB.SelectAttr("height")

	leftA := parseFloat(xA)
	rightA := leftA + parseFloat(widthA)
	topA := parseFloat(yA)
	bottomA := topA + parseFloat(heightA)

	leftB := parseFloat(xB)
	rightB := leftB + parseFloat(widthB)
	topB := parseFloat(yB)
	bottomB := topB + parseFloat(heightB)

	return !(rightA <= leftB || leftA >= rightB || bottomA <= topB || topA >= bottomB)
}

// Helper function to parse a float value from string
func parseFloat(value string) float64 {
	result, _ := strconv.ParseFloat(value, 64)
	return result
}

func printPaths(nodes []*xmlquery.Node) string {
	str := []string{}

	for _, n := range nodes {
		str = append(str, n.SelectAttr("d"))
	}

	return strings.Join(str, "\n")
}
