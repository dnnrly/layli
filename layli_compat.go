package layli

// This file provides backward compatibility for the root package.
// All layout/rendering logic has been moved to internal/layout/.
//
// These are type and function aliases for any code that imports from
// the github.com/dnnrly/layli root package.

import layoutpkg "github.com/dnnrly/layli/internal/layout"

// Legacy type aliases
type OutputFunc = layoutpkg.OutputFunc
type Diagram = layoutpkg.Diagram
type Config = layoutpkg.Config
type ConfigNode = layoutpkg.ConfigNode
type ConfigNodes = layoutpkg.ConfigNodes
type ConfigEdge = layoutpkg.ConfigEdge
type ConfigEdges = layoutpkg.ConfigEdges
type ConfigPath = layoutpkg.ConfigPath
type ConfigStyles = layoutpkg.ConfigStyles
type Layout = layoutpkg.Layout
type LayoutNode = layoutpkg.LayoutNode
type LayoutNodes = layoutpkg.LayoutNodes
type LayoutPath = layoutpkg.LayoutPath
type LayoutPaths = layoutpkg.LayoutPaths
type LayoutDrawer = layoutpkg.LayoutDrawer
type PathFinder = layoutpkg.PathFinder
type CreateFinder = layoutpkg.CreateFinder
type VertexMap = layoutpkg.VertexMap
type Arc = layoutpkg.Arc
type Arcs = layoutpkg.Arcs
type Position = layoutpkg.Position
type Point = layoutpkg.Point
type Points = layoutpkg.Points

// Legacy function aliases
var (
	NewConfigFromFile              = layoutpkg.NewConfigFromFile
	NewLayoutFromConfig            = layoutpkg.NewLayoutFromConfig
	NewLayout                       = layoutpkg.NewLayout
	NewLayoutNode                   = layoutpkg.NewLayoutNode
	NewVertexMap                    = layoutpkg.NewVertexMap
	BuildVertexMap                  = layoutpkg.BuildVertexMap
	PythagoreanDistance             = layoutpkg.PythagoreanDistance
	LayoutFlowSquare                = layoutpkg.LayoutFlowSquare
	LayoutTopologicalSort           = layoutpkg.LayoutTopologicalSort
	LayoutTarjan                    = layoutpkg.LayoutTarjan
	LayoutRandomShortestSquare      = layoutpkg.LayoutRandomShortestSquare
	LayoutAbsolute                  = layoutpkg.LayoutAbsolute
	AbsoluteFromSVG                 = layoutpkg.AbsoluteFromSVG
)
