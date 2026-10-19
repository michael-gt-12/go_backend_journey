package day2

import "math"

// Shape is implemented by any type that can report its area and perimeter.
// In Go, a type satisfies an interface implicitly: no "implements" keyword.
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle is a shape with a width and a height.
type Rectangle struct {
	Width  float64
	Height float64
}

// Area returns width * height.
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter returns 2 * (width + height).
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Circle is a shape defined by its radius.
type Circle struct {
	Radius float64
}

// Area returns pi * r^2.
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Perimeter returns the circumference, 2 * pi * r.
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Compile-time checks that both types really satisfy Shape.
var (
	_ Shape = Rectangle{}
	_ Shape = Circle{}
)
