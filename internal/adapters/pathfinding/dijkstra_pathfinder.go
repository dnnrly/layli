package pathfinding

import (
	"fmt"

	"github.com/dnnrly/layli/internal/adapters"
	"github.com/dnnrly/layli/internal/domain"
	"github.com/dnnrly/layli/layout"
	"github.com/dnnrly/layli/pathfinder/dijkstra"
)

type DijkstraPathfinder struct{}

func NewDijkstraPathfinder() *DijkstraPathfinder {
	return &DijkstraPathfinder{}
}

func (p *DijkstraPathfinder) FindPaths(diagram *domain.Diagram) error {
	cfg := adapters.ToLayoutConfigWithFullPaths(diagram)

	finder := func(start, end dijkstra.Point) layout.PathFinder {
		return createPathfinder(start, end, cfg.Path)
	}

	layoutObj, err := layout.NewLayoutFromConfig(finder, &cfg)
	if err != nil {
		return fmt.Errorf("finding paths: %w", err)
	}

	for i := range diagram.Edges {
		lp := findMatchingPath(layoutObj.Paths, diagram.Edges[i])
		if lp == nil {
			return fmt.Errorf("pathfinder could not calculate path for edge: %s -> %s", diagram.Edges[i].From, diagram.Edges[i].To)
		}

		points := make([]domain.Position, len(lp.Points))
		for j, pt := range lp.Points {
			points[j] = domain.Position{X: int(pt.X), Y: int(pt.Y)}
		}
		diagram.Edges[i].Path = &domain.Path{Points: points}
	}

	return nil
}

// createPathfinder creates the appropriate pathfinding algorithm based on configuration
func createPathfinder(start, end dijkstra.Point, pathConfig layout.ConfigPath) layout.PathFinder {
	switch pathConfig.Algorithm {
	case "astar":
		heuristic := getHeuristic(pathConfig.Heuristic)
		return dijkstra.NewAStarPathFinder(start, end, heuristic)
	case "bidirectional":
		return dijkstra.NewBidirectionalPathFinder(start, end)
	default:
		// Default to Dijkstra for backward compatibility
		return dijkstra.NewPathFinder(start, end)
	}
}

// getHeuristic returns the appropriate heuristic function for A*
func getHeuristic(heuristicName string) dijkstra.HeuristicFunction {
	switch heuristicName {
	case "manhattan":
		return dijkstra.ManhattanDistance
	default:
		// Default to Euclidean distance
		return dijkstra.EuclideanDistance
	}
}

func findMatchingPath(paths layout.LayoutPaths, edge domain.Edge) *layout.LayoutPath {
	for _, lp := range paths {
		if lp.From == edge.From && lp.To == edge.To {
			return &lp
		}
	}
	return nil
}

