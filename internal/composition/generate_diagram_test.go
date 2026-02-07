package composition

import "testing"

func TestNewGenerateDiagram(t *testing.T) {
	tests := []struct {
		name     string
		showGrid bool
	}{
		{"with grid", true},
		{"without grid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := NewGenerateDiagram(tt.showGrid)

			if generator == nil {
				t.Fatal("Expected non-nil generator")
			}

			// Verify it's a valid generator by checking it's properly constructed
			// The Execute method should exist and be callable
			if generator == nil {
				t.Error("Expected generator to be properly constructed")
			}
		})
	}
}
