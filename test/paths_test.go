package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExtractPath(t *testing.T) {
	segments := NewSegments("M 200 160 L 200 140 L 100 140 L 100 100 z")

	assert.EqualValues(t, []Segment{
		{Start: Point{X: 200, Y: 160}, End: Point{X: 200, Y: 140}},
		{Start: Point{X: 200, Y: 140}, End: Point{X: 100, Y: 140}},
		{Start: Point{X: 100, Y: 140}, End: Point{X: 100, Y: 100}},
	}, segments)
}

func Test_CrossingSegments(t *testing.T) {
	segments := NewSegments("M 200 160 L 200 140 L 100 140 L 100 100")

	assert.False(t, segments.Crosses(NewSegments("M 300 300 L 200 200")))
	assert.True(t, segments.Crosses(NewSegments("M 200 160 L 200 140")))
	assert.True(t, segments.Crosses(NewSegments("M 190 140 L 150 140 L 150 200")))
	assert.True(t, segments.Crosses(NewSegments("M 200 140 L 50 140 L 100 100")))
}
