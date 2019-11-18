package types

// Color is a tuple representing an rgb value
type Color struct {
	red   float64
	green float64
	blue  float64
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
func (c *Color) Add(c2 *Color) (*Color, error) {
	return &Color{
		red:   c.red + c2.red,
		green: c.green + c2.green,
		blue:  c.blue + c2.blue,
	}, nil
}

// Subtract returns the difference between two colors
func (c *Color) Subtract(c2 *Color) (*Color, error) {
	return &Color{
		red:   c.red - c2.red,
		green: c.green - c2.green,
		blue:  c.blue - c2.blue,
	}, nil
}

// Multiply returns a color multiplied by a value
func (c *Color) Multiply(val float64) *Color {
	return &Color{
		red:   c.red * val,
		green: c.green * val,
		blue:  c.blue * val,
	}
}

// MultiplyColor multiplies two colors by each other and returns the result
func (c *Color) MultiplyColor(c2 *Color) *Color {
	return &Color{
		red:   c.red * c2.red,
		green: c.green * c2.green,
		blue:  c.blue * c2.blue,
	}
}
