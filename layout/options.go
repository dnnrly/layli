package layout

import (
	"github.com/dnnrly/layli/internal/domain"
)

// GetLayoutOptions returns all available layout algorithm names
//
// NOTE: When adding a new layout algorithm, remember to:
// 1. Add constant to internal/domain/diagram.go
// 2. Implement function in arrangements.go
// 3. Register in selectArrangement() in arrangements.go
// 4. Register in selectArranger() in internal/adapters/layout/engine.go
// 5. Add it to this function (so it's discoverable)
// See CONTRIBUTING_LAYOUTS.md for detailed steps.
func GetLayoutOptions() []string {
	return []string{
		string(domain.LayoutFlowSquare),
		string(domain.LayoutTopoSort),
		string(domain.LayoutTarjan),
		string(domain.LayoutAbsolute),
		string(domain.LayoutRandomShortest),
	}
}

// GetPathfindingAlgorithms returns all available pathfinding algorithm names
//
// NOTE: When adding a new pathfinding algorithm, add it to:
// 1. internal/domain/diagram.go (PathfindingAlgorithm constant)
// 2. internal/adapters/config/yaml_parser.go (validation map)
// 3. This function (so it's discoverable)
func GetPathfindingAlgorithms() []string {
	return []string{
		string(domain.PathfindingDijkstra),
		string(domain.PathfindingAStar),
		string(domain.PathfindingBidirectional),
	}
}

// GetHeuristics returns all available heuristic function names
//
// NOTE: When adding a new heuristic, add it to:
// 1. internal/domain/diagram.go (PathfindingHeuristic constant)
// 2. internal/adapters/config/yaml_parser.go (validation map)
// 3. This function (so it's discoverable)
func GetHeuristics() []string {
	return []string{
		string(domain.HeuristicEuclidean),
		string(domain.HeuristicManhattan),
	}
}
