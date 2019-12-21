package image

import (
	"github.com/sjberman/golang-ray-tracer/pkg/base"
)

// Color is a tuple representing an rgb value
type Color struct {
	red   float64
	green float64
	blue  float64
}

var Black = &Color{
	red:   0,
	green: 0,
	blue:  0,
}

var White = &Color{
	red:   1,
	green: 1,
	blue:  1,
}

// NewColor returns a new Color object
func NewColor(r, g, b float64) *Color {
	return &Color{
		red:   r,
		green: g,
		blue:  b,
	}
}

// Add adds two colors together and returns the result
func (c *Color) Add(c2 *Color) *Color {
	return NewColor(c.red+c2.red, c.green+c2.green, c.blue+c2.blue)
}

// Subtract returns the difference between two colors
func (c *Color) Subtract(c2 *Color) *Color {
	return NewColor(c.red-c2.red, c.green-c2.green, c.blue-c2.blue)
}

// Multiply returns a color multiplied by a value
func (c *Color) Multiply(val float64) *Color {
	return NewColor(c.red*val, c.green*val, c.blue*val)
}

// MultiplyColor multiplies two colors by each other and returns the result
func (c *Color) MultiplyColor(c2 *Color) *Color {
	return NewColor(c.red*c2.red, c.green*c2.green, c.blue*c2.blue)
}

// Equals returns whether or not two colors are equal to each other
func (c *Color) Equals(c2 *Color) bool {
	if !base.EqualFloats(c.red, c2.red) {
		return false
	}
	if !base.EqualFloats(c.green, c2.green) {
		return false
	}
	return base.EqualFloats(c.blue, c2.blue)
}
