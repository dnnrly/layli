package layout

import (
	"testing"
)

func TestGetLayoutOptions(t *testing.T) {
	layouts := GetLayoutOptions()

	if len(layouts) != 5 {
		t.Errorf("expected 5 layouts, got %d", len(layouts))
	}

	expected := []string{"flow-square", "topo-sort", "tarjan", "absolute", "random-shortest-square"}
	for i, exp := range expected {
		if i >= len(layouts) || layouts[i] != exp {
			t.Errorf("expected layout %s at index %d", exp, i)
		}
	}
}

func TestGetPathfindingAlgorithms(t *testing.T) {
	algorithms := GetPathfindingAlgorithms()

	if len(algorithms) != 3 {
		t.Errorf("expected 3 algorithms, got %d", len(algorithms))
	}

	expected := []string{"dijkstra", "astar", "bidirectional"}
	for i, exp := range expected {
		if i >= len(algorithms) || algorithms[i] != exp {
			t.Errorf("expected algorithm %s at index %d", exp, i)
		}
	}
}

func TestGetHeuristics(t *testing.T) {
	heuristics := GetHeuristics()

	if len(heuristics) != 2 {
		t.Errorf("expected 2 heuristics, got %d", len(heuristics))
	}

	expected := []string{"euclidean", "manhattan"}
	for i, exp := range expected {
		if i >= len(heuristics) || heuristics[i] != exp {
			t.Errorf("expected heuristic %s at index %d", exp, i)
		}
	}
}
