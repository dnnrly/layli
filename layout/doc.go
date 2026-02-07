// Package layout provides grid-based layout and rendering algorithms.
//
// This package is the Rendering Domain - a peer to the Application Domain
// in internal/domain/. While internal/domain/ contains business logic and
// entities (Diagram, Node, Edge), this package handles the computational
// details of arranging nodes on a grid and routing edges (paths).
//
// Key types:
//   - Config: Grid/layout configuration (spacing, margins, dimensions)
//   - LayoutNode: A node positioned on the grid (distinct from domain.Node)
//   - LayoutPath: An edge's route through the grid (distinct from domain.Edge.Path)
//   - VertexMap: Grid vertex tracking for pathfinding
//
// Layout algorithms:
//   - LayoutFlowSquare: Arrange nodes in a square grid pattern
//   - LayoutTopologicalSort: Single-row topological ordering
//   - LayoutTarjan: Multi-row layout using Tarjan's algorithm
//   - LayoutAbsolute: Use explicit positions from config
//   - LayoutRandomShortestSquare: Multiple attempts to minimize path length
//
// Internal adapters (internal/adapters/) convert between the Application
// Domain (internal/domain/) and this Rendering Domain.
package layout
