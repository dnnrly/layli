package domain

import (
	"testing"
)

func TestEdgeValidate_Success(t *testing.T) {
	tests := []struct {
		name string
		edge Edge
	}{
		{
			name: "simple edge",
			edge: Edge{ID: "e1", From: "a", To: "b"},
		},
		{
			name: "edge with style",
			edge: Edge{ID: "e2", From: "node1", To: "node2", Style: "stroke:red"},
		},
		{
			name: "edge with class",
			edge: Edge{ID: "e3", From: "src", To: "dst", Class: "class-1"},
		},
		{
			name: "edge with path",
			edge: Edge{
				ID:   "e4",
				From: "a",
				To:   "b",
				Path: &Path{Points: []Position{{X: 0, Y: 0}, {X: 1, Y: 1}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.edge.Validate(); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestEdgeValidate_EmptyFrom(t *testing.T) {
	e := Edge{ID: "e1", From: "", To: "b"}
	if err := e.Validate(); err == nil {
		t.Fatal("expected error for empty from node")
	}
}

func TestEdgeValidate_EmptyTo(t *testing.T) {
	e := Edge{ID: "e1", From: "a", To: ""}
	if err := e.Validate(); err == nil {
		t.Fatal("expected error for empty to node")
	}
}

func TestEdgeValidate_SelfLoop(t *testing.T) {
	e := Edge{ID: "e1", From: "a", To: "a"}
	if err := e.Validate(); err == nil {
		t.Fatal("expected error for self-loop edge")
	}
}

func TestEdgeString(t *testing.T) {
	e := Edge{ID: "e1", From: "node1", To: "node2"}
	result := e.String()
	if result != "Edge(node1 -> node2)" {
		t.Errorf("String() = %q, want %q", result, "Edge(node1 -> node2)")
	}
}

func TestPathLength_Empty(t *testing.T) {
	if p := (*Path)(nil); p.Length() != 0 {
		t.Errorf("empty path length = %v, want 0", p.Length())
	}
}

func TestPathLength_SinglePoint(t *testing.T) {
	p := &Path{Points: []Position{{X: 0, Y: 0}}}
	if p.Length() != 0 {
		t.Errorf("single point path length = %v, want 0", p.Length())
	}
}

func TestPathLength_TwoPoints(t *testing.T) {
	p := &Path{Points: []Position{{X: 0, Y: 0}, {X: 3, Y: 4}}}
	length := p.Length()
	expected := 5.0 // 3-4-5 triangle
	if length != expected {
		t.Errorf("two point path length = %v, want %v", length, expected)
	}
}

func TestPathLength_MultiplePoints(t *testing.T) {
	p := &Path{Points: []Position{
		{X: 0, Y: 0},
		{X: 3, Y: 4},
		{X: 3, Y: 8},
	}}
	length := p.Length()
	expected := 9.0 // 5 + 4 = 9
	if length != expected {
		t.Errorf("multiple point path length = %v, want %v", length, expected)
	}
}

func TestPathCorners_Nil(t *testing.T) {
	if p := (*Path)(nil); p.Corners() != 0 {
		t.Errorf("nil path corners = %d, want 0", p.Corners())
	}
}

func TestPathCorners_LessThanThreePoints(t *testing.T) {
	tests := []struct {
		name   string
		path   *Path
		expect int
	}{
		{"empty path", &Path{Points: []Position{}}, 0},
		{"one point", &Path{Points: []Position{{X: 0, Y: 0}}}, 0},
		{"two points", &Path{Points: []Position{{X: 0, Y: 0}, {X: 1, Y: 1}}}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if corners := tt.path.Corners(); corners != tt.expect {
				t.Errorf("corners = %d, want %d", corners, tt.expect)
			}
		})
	}
}

func TestPathCorners_StraightLine(t *testing.T) {
	// Straight horizontal line: no corners
	p := &Path{Points: []Position{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 2, Y: 0},
	}}
	if corners := p.Corners(); corners != 0 {
		t.Errorf("straight line corners = %d, want 0", corners)
	}
}

func TestPathCorners_RightAngle(t *testing.T) {
	// L-shaped path: 1 corner
	p := &Path{Points: []Position{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 1, Y: 1},
	}}
	if corners := p.Corners(); corners != 1 {
		t.Errorf("L-shaped path corners = %d, want 1", corners)
	}
}

func TestPathCorners_ZigZag(t *testing.T) {
	// Zig-zag path: 2 corners
	p := &Path{Points: []Position{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 1, Y: 1},
		{X: 2, Y: 1},
	}}
	if corners := p.Corners(); corners != 2 {
		t.Errorf("zig-zag path corners = %d, want 2", corners)
	}
}

func TestPathString_Nil(t *testing.T) {
	if p := (*Path)(nil); p.String() != "Path(empty)" {
		t.Errorf("nil path String() = %q, want %q", p.String(), "Path(empty)")
	}
}

func TestPathString_Empty(t *testing.T) {
	p := &Path{Points: []Position{}}
	if p.String() != "Path(empty)" {
		t.Errorf("empty path String() = %q, want %q", p.String(), "Path(empty)")
	}
}

func TestPathString_WithPoints(t *testing.T) {
	p := &Path{Points: []Position{
		{X: 0, Y: 0},
		{X: 1, Y: 1},
		{X: 2, Y: 2},
	}}
	result := p.String()
	expected := "Path(3 points)"
	if result != expected {
		t.Errorf("path String() = %q, want %q", result, expected)
	}
}
