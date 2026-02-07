package domain

import (
	"testing"
)

func TestDiagramValidate_Success(t *testing.T) {
	tests := []struct {
		name    string
		diagram Diagram
	}{
		{
			name: "single node diagram",
			diagram: Diagram{
				Nodes: []Node{
					{ID: "a", Width: 5, Height: 5},
				},
				Edges: []Edge{},
				Config: DiagramConfig{
					NodeWidth:      5,
					NodeHeight:     5,
					Border:         1,
					Margin:         1,
					PathAttempts:   100,
					LayoutAttempts: 100,
				},
			},
		},
		{
			name: "two nodes with edge",
			diagram: Diagram{
				Nodes: []Node{
					{ID: "a", Width: 5, Height: 5},
					{ID: "b", Width: 5, Height: 5},
				},
				Edges: []Edge{
					{ID: "e1", From: "a", To: "b"},
				},
				Config: DiagramConfig{
					NodeWidth:      5,
					NodeHeight:     5,
					Margin:         1,
					PathAttempts:   100,
					LayoutAttempts: 100,
				},
			},
		},
		{
			name: "multiple edges",
			diagram: Diagram{
				Nodes: []Node{
					{ID: "a", Width: 5, Height: 5},
					{ID: "b", Width: 5, Height: 5},
					{ID: "c", Width: 5, Height: 5},
				},
				Edges: []Edge{
					{ID: "e1", From: "a", To: "b"},
					{ID: "e2", From: "b", To: "c"},
					{ID: "e3", From: "a", To: "c"},
				},
				Config: DiagramConfig{
					NodeWidth:      5,
					NodeHeight:     5,
					Margin:         1,
					PathAttempts:   100,
					LayoutAttempts: 100,
				},
			},
		},
		{
			name: "maximum margin",
			diagram: Diagram{
				Nodes: []Node{
					{ID: "a", Width: 5, Height: 5},
				},
				Edges: []Edge{},
				Config: DiagramConfig{
					NodeWidth:      5,
					NodeHeight:     5,
					Margin:         10,
					PathAttempts:   100,
					LayoutAttempts: 100,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.diagram.Validate(); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestDiagramValidate_NoNodes(t *testing.T) {
	d := Diagram{
		Nodes: []Node{},
		Edges: []Edge{},
		Config: DiagramConfig{
			NodeWidth:      5,
			NodeHeight:     5,
			PathAttempts:   100,
			LayoutAttempts: 100,
		},
	}
	if err := d.Validate(); err == nil {
		t.Fatal("expected error for no nodes")
	}
}

func TestDiagramValidate_InvalidNodeDimensions(t *testing.T) {
	tests := []struct {
		name      string
		nodeWidth int
		nodeHeight int
	}{
		{"zero width", 0, 5},
		{"zero height", 5, 0},
		{"negative width", -1, 5},
		{"negative height", 5, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Diagram{
				Nodes: []Node{
					{ID: "a", Width: 5, Height: 5},
				},
				Edges: []Edge{},
				Config: DiagramConfig{
					NodeWidth:      tt.nodeWidth,
					NodeHeight:     tt.nodeHeight,
					PathAttempts:   100,
					LayoutAttempts: 100,
				},
			}
			if err := d.Validate(); err == nil {
				t.Fatal("expected error for invalid node dimensions")
			}
		})
	}
}

func TestDiagramValidate_InvalidMargin(t *testing.T) {
	tests := []struct {
		name   string
		margin int
	}{
		{"negative margin", -1},
		{"margin too large", 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Diagram{
				Nodes: []Node{
					{ID: "a", Width: 5, Height: 5},
				},
				Edges: []Edge{},
				Config: DiagramConfig{
					NodeWidth:      5,
					NodeHeight:     5,
					Margin:         tt.margin,
					PathAttempts:   100,
					LayoutAttempts: 100,
				},
			}
			if err := d.Validate(); err == nil {
				t.Fatalf("expected error for invalid margin: %d", tt.margin)
			}
		})
	}
}

func TestDiagramValidate_InvalidPathAttempts(t *testing.T) {
	tests := []struct {
		name         string
		pathAttempts int
	}{
		{"zero attempts", 0},
		{"negative attempts", -100},
		{"too many attempts", 10001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Diagram{
				Nodes: []Node{
					{ID: "a", Width: 5, Height: 5},
				},
				Edges: []Edge{},
				Config: DiagramConfig{
					NodeWidth:      5,
					NodeHeight:     5,
					PathAttempts:   tt.pathAttempts,
					LayoutAttempts: 100,
				},
			}
			if err := d.Validate(); err == nil {
				t.Fatalf("expected error for invalid path attempts: %d", tt.pathAttempts)
			}
		})
	}
}

func TestDiagramValidate_InvalidLayoutAttempts(t *testing.T) {
	tests := []struct {
		name            string
		layoutAttempts int
	}{
		{"zero attempts", 0},
		{"negative attempts", -100},
		{"too many attempts", 10001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Diagram{
				Nodes: []Node{
					{ID: "a", Width: 5, Height: 5},
				},
				Edges: []Edge{},
				Config: DiagramConfig{
					NodeWidth:      5,
					NodeHeight:     5,
					PathAttempts:   100,
					LayoutAttempts: tt.layoutAttempts,
				},
			}
			if err := d.Validate(); err == nil {
				t.Fatalf("expected error for invalid layout attempts: %d", tt.layoutAttempts)
			}
		})
	}
}

func TestDiagramValidate_InvalidNode(t *testing.T) {
	d := Diagram{
		Nodes: []Node{
			{ID: "", Width: 5, Height: 5}, // Invalid node (empty ID)
		},
		Edges: []Edge{},
		Config: DiagramConfig{
			NodeWidth:      5,
			NodeHeight:     5,
			PathAttempts:   100,
			LayoutAttempts: 100,
		},
	}
	if err := d.Validate(); err == nil {
		t.Fatal("expected error for invalid node")
	}
}

func TestDiagramValidate_EdgeReferencesNonexistentFromNode(t *testing.T) {
	d := Diagram{
		Nodes: []Node{
			{ID: "a", Width: 5, Height: 5},
			{ID: "b", Width: 5, Height: 5},
		},
		Edges: []Edge{
			{ID: "e1", From: "nonexistent", To: "b"},
		},
		Config: DiagramConfig{
			NodeWidth:      5,
			NodeHeight:     5,
			PathAttempts:   100,
			LayoutAttempts: 100,
		},
	}
	if err := d.Validate(); err == nil {
		t.Fatal("expected error for nonexistent from node")
	}
}

func TestDiagramValidate_EdgeReferencesNonexistentToNode(t *testing.T) {
	d := Diagram{
		Nodes: []Node{
			{ID: "a", Width: 5, Height: 5},
			{ID: "b", Width: 5, Height: 5},
		},
		Edges: []Edge{
			{ID: "e1", From: "a", To: "nonexistent"},
		},
		Config: DiagramConfig{
			NodeWidth:      5,
			NodeHeight:     5,
			PathAttempts:   100,
			LayoutAttempts: 100,
		},
	}
	if err := d.Validate(); err == nil {
		t.Fatal("expected error for nonexistent to node")
	}
}

func TestDiagramValidate_InvalidEdge(t *testing.T) {
	d := Diagram{
		Nodes: []Node{
			{ID: "a", Width: 5, Height: 5},
		},
		Edges: []Edge{
			{ID: "e1", From: "a", To: "a"}, // Invalid: self-loop
		},
		Config: DiagramConfig{
			NodeWidth:      5,
			NodeHeight:     5,
			PathAttempts:   100,
			LayoutAttempts: 100,
		},
	}
	if err := d.Validate(); err == nil {
		t.Fatal("expected error for self-loop edge")
	}
}
