package test_test

import (
	"github.com/antchfx/xmlquery"
)

// getNodeIds takes an *xmlquery.Node and returns the "id" attribute of all elements
func getNodeIds(node *xmlquery.Node) []string {
	var ids []string

	// Iterate over child nodes
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == xmlquery.ElementNode {
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
