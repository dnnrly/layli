package domain

import (
	"math"
	"testing"
)

func TestPositionDistance(t *testing.T) {
	tests := []struct {
		name     string
		from     Position
		to       Position
		expected float64
	}{
		{
			name:     "same position",
			from:     Position{X: 0, Y: 0},
			to:       Position{X: 0, Y: 0},
			expected: 0,
		},
		{
			name:     "3-4-5 triangle",
			from:     Position{X: 0, Y: 0},
			to:       Position{X: 3, Y: 4},
			expected: 5,
		},
		{
			name:     "unit distance horizontal",
			from:     Position{X: 0, Y: 0},
			to:       Position{X: 1, Y: 0},
			expected: 1,
		},
		{
			name:     "unit distance vertical",
			from:     Position{X: 0, Y: 0},
			to:       Position{X: 0, Y: 1},
			expected: 1,
		},
		{
			name:     "negative coordinates",
			from:     Position{X: -3, Y: -4},
			to:       Position{X: 0, Y: 0},
			expected: 5,
		},
		{
			name:     "distance is symmetric",
			from:     Position{X: 5, Y: 5},
			to:       Position{X: 2, Y: 1},
			expected: 5, // sqrt(9 + 16)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.from.Distance(tt.to)
			if math.Abs(d-tt.expected) > 1e-10 {
				t.Errorf("Distance() = %v, want %v", d, tt.expected)
			}
		})
	}
}

func TestPositionAdd(t *testing.T) {
	tests := []struct {
		name   string
		pos    Position
		dx     int
		dy     int
		expect Position
	}{
		{
			name:   "add positive values",
			pos:    Position{X: 0, Y: 0},
			dx:     3,
			dy:     4,
			expect: Position{X: 3, Y: 4},
		},
		{
			name:   "add negative values",
			pos:    Position{X: 5, Y: 5},
			dx:     -2,
			dy:     -3,
			expect: Position{X: 3, Y: 2},
		},
		{
			name:   "add zero",
			pos:    Position{X: 10, Y: 20},
			dx:     0,
			dy:     0,
			expect: Position{X: 10, Y: 20},
		},
		{
			name:   "add mixed signs",
			pos:    Position{X: 5, Y: 5},
			dx:     3,
			dy:     -2,
			expect: Position{X: 8, Y: 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pos.Add(tt.dx, tt.dy)
			if result != tt.expect {
				t.Errorf("Add(%d, %d) = %v, want %v", tt.dx, tt.dy, result, tt.expect)
			}
		})
	}
}

func TestPositionString(t *testing.T) {
	tests := []struct {
		name   string
		pos    Position
		expect string
	}{
		{name: "origin", pos: Position{X: 0, Y: 0}, expect: "(0,0)"},
		{name: "positive", pos: Position{X: 3, Y: 4}, expect: "(3,4)"},
		{name: "negative", pos: Position{X: -5, Y: -10}, expect: "(-5,-10)"},
		{name: "mixed", pos: Position{X: 5, Y: -3}, expect: "(5,-3)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pos.String()
			if result != tt.expect {
				t.Errorf("String() = %q, want %q", result, tt.expect)
			}
		})
	}
}

func TestBoundsWidth(t *testing.T) {
	tests := []struct {
		name   string
		bounds Bounds
		expect int
	}{
		{
			name:   "simple bounds",
			bounds: Bounds{Min: Position{X: 0, Y: 0}, Max: Position{X: 10, Y: 20}},
			expect: 10,
		},
		{
			name:   "non-zero origin",
			bounds: Bounds{Min: Position{X: 5, Y: 5}, Max: Position{X: 15, Y: 25}},
			expect: 10,
		},
		{
			name:   "negative coordinates",
			bounds: Bounds{Min: Position{X: -10, Y: -10}, Max: Position{X: 0, Y: 0}},
			expect: 10,
		},
		{
			name:   "zero width",
			bounds: Bounds{Min: Position{X: 5, Y: 5}, Max: Position{X: 5, Y: 10}},
			expect: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.bounds.Width()
			if w != tt.expect {
				t.Errorf("Width() = %d, want %d", w, tt.expect)
			}
		})
	}
}

func TestBoundsHeight(t *testing.T) {
	tests := []struct {
		name   string
		bounds Bounds
		expect int
	}{
		{
			name:   "simple bounds",
			bounds: Bounds{Min: Position{X: 0, Y: 0}, Max: Position{X: 10, Y: 20}},
			expect: 20,
		},
		{
			name:   "non-zero origin",
			bounds: Bounds{Min: Position{X: 5, Y: 5}, Max: Position{X: 15, Y: 25}},
			expect: 20,
		},
		{
			name:   "negative coordinates",
			bounds: Bounds{Min: Position{X: -10, Y: -10}, Max: Position{X: 0, Y: 0}},
			expect: 10,
		},
		{
			name:   "zero height",
			bounds: Bounds{Min: Position{X: 5, Y: 5}, Max: Position{X: 10, Y: 5}},
			expect: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.bounds.Height()
			if h != tt.expect {
				t.Errorf("Height() = %d, want %d", h, tt.expect)
			}
		})
	}
}

func TestBoundsContains(t *testing.T) {
	bounds := Bounds{Min: Position{X: 0, Y: 0}, Max: Position{X: 10, Y: 10}}

	tests := []struct {
		name   string
		pos    Position
		expect bool
	}{
		{name: "point at min", pos: Position{X: 0, Y: 0}, expect: true},
		{name: "point at max", pos: Position{X: 10, Y: 10}, expect: true},
		{name: "point inside", pos: Position{X: 5, Y: 5}, expect: true},
		{name: "point on edge", pos: Position{X: 0, Y: 5}, expect: true},
		{name: "point outside left", pos: Position{X: -1, Y: 5}, expect: false},
		{name: "point outside right", pos: Position{X: 11, Y: 5}, expect: false},
		{name: "point outside top", pos: Position{X: 5, Y: -1}, expect: false},
		{name: "point outside bottom", pos: Position{X: 5, Y: 11}, expect: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bounds.Contains(tt.pos)
			if result != tt.expect {
				t.Errorf("Contains(%v) = %v, want %v", tt.pos, result, tt.expect)
			}
		})
	}
}

func TestBoundsOverlaps(t *testing.T) {
	tests := []struct {
		name      string
		bounds1   Bounds
		bounds2   Bounds
		expect    bool
	}{
		{
			name:      "identical bounds",
			bounds1:   Bounds{Min: Position{X: 0, Y: 0}, Max: Position{X: 10, Y: 10}},
			bounds2:   Bounds{Min: Position{X: 0, Y: 0}, Max: Position{X: 10, Y: 10}},
			expect:    true,
		},
		{
			name:      "overlapping bounds",
			bounds1:   Bounds{Min: Position{X: 0, Y: 0}, Max: Position{X: 10, Y: 10}},
			bounds2:   Bounds{Min: Position{X: 5, Y: 5}, Max: Position{X: 15, Y: 15}},
			expect:    true,
		},
		{
			name:      "touching bounds (edge)",
			bounds1:   Bounds{Min: Position{X: 0, Y: 0}, Max: Position{X: 10, Y: 10}},
			bounds2:   Bounds{Min: Position{X: 10, Y: 0}, Max: Position{X: 20, Y: 10}},
			expect:    true,
		},
		{
			name:      "disjoint bounds",
			bounds1:   Bounds{Min: Position{X: 0, Y: 0}, Max: Position{X: 10, Y: 10}},
			bounds2:   Bounds{Min: Position{X: 20, Y: 20}, Max: Position{X: 30, Y: 30}},
			expect:    false,
		},
		{
			name:      "one inside the other",
			bounds1:   Bounds{Min: Position{X: 0, Y: 0}, Max: Position{X: 20, Y: 20}},
			bounds2:   Bounds{Min: Position{X: 5, Y: 5}, Max: Position{X: 10, Y: 10}},
			expect:    true,
		},
		{
			name:      "adjacent but not overlapping",
			bounds1:   Bounds{Min: Position{X: 0, Y: 0}, Max: Position{X: 10, Y: 10}},
			bounds2:   Bounds{Min: Position{X: 10, Y: 10}, Max: Position{X: 20, Y: 20}},
			expect:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.bounds1.Overlaps(tt.bounds2)
			if result != tt.expect {
				t.Errorf("Overlaps(%v, %v) = %v, want %v",
					tt.bounds1, tt.bounds2, result, tt.expect)
			}
		})
	}
}

func TestBoundsOverlapsSymmetric(t *testing.T) {
	// Verify that Overlaps is symmetric
	b1 := Bounds{Min: Position{X: 0, Y: 0}, Max: Position{X: 10, Y: 10}}
	b2 := Bounds{Min: Position{X: 5, Y: 5}, Max: Position{X: 15, Y: 15}}

	if b1.Overlaps(b2) != b2.Overlaps(b1) {
		t.Error("Overlaps is not symmetric")
	}
}

func TestBoundsString(t *testing.T) {
	bounds := Bounds{Min: Position{X: 0, Y: 0}, Max: Position{X: 10, Y: 20}}
	result := bounds.String()
	expected := "[(0,0) to (10,20)]"
	if result != expected {
		t.Errorf("String() = %q, want %q", result, expected)
	}
}
