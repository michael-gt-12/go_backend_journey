package day2

import (
	"math"
	"testing"
)

const epsilon = 1e-9

func TestShapes(t *testing.T) {
	tests := []struct {
		name          string
		shape         Shape
		wantArea      float64
		wantPerimeter float64
	}{
		{"rectangle 3x4", Rectangle{Width: 3, Height: 4}, 12, 14},
		{"square 5x5", Rectangle{Width: 5, Height: 5}, 25, 20},
		{"zero rectangle", Rectangle{}, 0, 0},
		{"circle r=1", Circle{Radius: 1}, math.Pi, 2 * math.Pi},
		{"circle r=2.5", Circle{Radius: 2.5}, math.Pi * 6.25, 5 * math.Pi},
		{"zero circle", Circle{}, 0, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.shape.Area(); math.Abs(got-tc.wantArea) > epsilon {
				t.Errorf("Area() = %v, want %v", got, tc.wantArea)
			}
			if got := tc.shape.Perimeter(); math.Abs(got-tc.wantPerimeter) > epsilon {
				t.Errorf("Perimeter() = %v, want %v", got, tc.wantPerimeter)
			}
		})
	}
}
