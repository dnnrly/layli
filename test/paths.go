package test

import (
	"math"
	"strconv"
	"strings"
)

type Point struct {
	X, Y float64
}

type Segment struct {
	Start, End Point
}

type Segments []Segment

func NewSegments(path string) Segments {
	path = strings.ReplaceAll(path, "M ", "")
	parts := strings.Split(path, " L ")

	points := []Point{}
	for _, p := range parts {
		p = strings.Trim(p, " ")
		axes := strings.Split(p, " ")
		x, _ := strconv.Atoi(axes[0])
		y, _ := strconv.Atoi(axes[1])

		points = append(points, Point{X: float64(x), Y: float64(y)})
	}

	segments := []Segment{}
	for len(points) > 1 {
		segments = append(segments, Segment{
			Start: points[0],
			End:   points[1],
		})
		points = points[1:]
	}

	return segments
}

func doSegmentsIntersectOrCoincident(p1, p2, p3, p4 Point) bool {
	denom := (p4.Y-p3.Y)*(p2.X-p1.X) - (p4.X-p3.X)*(p2.Y-p1.Y)
	numA := (p4.X-p3.X)*(p1.Y-p3.Y) - (p4.Y-p3.Y)*(p1.X-p3.X)
	numB := (p2.X-p1.X)*(p1.Y-p3.Y) - (p2.Y-p1.Y)*(p1.X-p3.X)

	if denom == 0 {
		// The segments are parallel
		if numA == 0 && numB == 0 {
			// The segments are coincident
			// Check if the segments overlap
			if p1.X <= math.Max(p3.X, p4.X) && p1.X >= math.Min(p3.X, p4.X) &&
				p1.Y <= math.Max(p3.Y, p4.Y) && p1.Y >= math.Min(p3.Y, p4.Y) {
				return true
			}
			if p2.X <= math.Max(p3.X, p4.X) && p2.X >= math.Min(p3.X, p4.X) &&
				p2.Y <= math.Max(p3.Y, p4.Y) && p2.Y >= math.Min(p3.Y, p4.Y) {
				return true
			}
		}
		// The segments are parallel but non-coincident
		return false
	}

	uA := numA / denom
	uB := numB / denom

	return uA >= 0 && uA <= 1 && uB >= 0 && uB <= 1
}

func doSegmentsMeet(p1, p2, p3, p4 Point) bool {
	return doSegmentsIntersectOrCoincident(p1, p2, p3, p4)
}

func (s Segments) Crosses(targetPath Segments) bool {
	for _, s1 := range s {
		for _, s2 := range targetPath {
			if doSegmentsMeet(s1.Start, s1.End, s2.Start, s2.End) {
				return true
			}
		}
	}

	return false
}
