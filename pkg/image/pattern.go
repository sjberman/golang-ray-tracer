package image

import (
	"math"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
)

// Pattern represents a color pattern
type Pattern interface {
	GetColors() []*Color
	GetTransform() *base.Matrix
	SetTransform(*base.Matrix)
	PatternAt(*base.Tuple) *Color
}

// PatternObject is the base implementation of the pattern interface
type PatternObject struct {
	color1        *Color
	color2        *Color
	transform     *base.Matrix
	patternAtFunc func(*base.Tuple, *PatternObject) *Color
}

// newPattern returns a new pattern object
func newPattern(
	color1, color2 *Color,
	patternFunc func(*base.Tuple, *PatternObject) *Color,
) *PatternObject {
	return &PatternObject{
		color1:        color1,
		color2:        color2,
		patternAtFunc: patternFunc,
		transform:     &base.Identity,
	}
}

// GetColors returns a list of colors for the pattern
func (p *PatternObject) GetColors() []*Color {
	return []*Color{p.color1, p.color2}
}

// GetTransform returns the pattern's transform matrix
func (p *PatternObject) GetTransform() *base.Matrix {
	return p.transform
}

// SetTransform sets the pattern's transform matrix
func (p *PatternObject) SetTransform(matrix *base.Matrix) {
	p.transform = matrix
}

// PatternAt returns the color at a specific point, based on the pattern
func (p *PatternObject) PatternAt(point *base.Tuple) *Color {
	return p.patternAtFunc(point, p)
}

// StripePattern represents a striped color pattern
type StripePattern struct {
	*PatternObject
}

// NewStripePattern returns a new StripePattern object
func NewStripePattern(color1, color2 *Color) *StripePattern {
	return &StripePattern{
		PatternObject: newPattern(color1, color2, stripeAt),
	}
}

// stripeAt returns the stripe color at a specific point
func stripeAt(point *base.Tuple, p *PatternObject) *Color {
	m := mod(point.GetX(), 2)
	if m >= 0 && m < 1 {
		return p.color1
	}
	return p.color2
}

// GradientPattern represents a gradient color pattern
type GradientPattern struct {
	*PatternObject
}

// NewGradientPattern returns a new GradientPattern object
func NewGradientPattern(color1, color2 *Color) *GradientPattern {
	return &GradientPattern{
		PatternObject: newPattern(color1, color2, gradientAt),
	}
}

// gradientAt returns the gradient color at a specific point
func gradientAt(point *base.Tuple, p *PatternObject) *Color {
	colors := p.GetColors()
	distance := colors[1].Subtract(colors[0])
	fraction := point.GetX() - math.Floor(point.GetX())
	return colors[0].Add(distance.Multiply(fraction))
}

// RingPattern represents a ring color pattern
type RingPattern struct {
	*PatternObject
}

// NewRingPattern returns a new RingPattern object
func NewRingPattern(color1, color2 *Color) *RingPattern {
	return &RingPattern{
		PatternObject: newPattern(color1, color2, ringAt),
	}
}

// ringAt returns the ring color at a specific point
func ringAt(point *base.Tuple, p *PatternObject) *Color {
	distance := math.Sqrt(math.Pow(point.GetX(), 2) + math.Pow(point.GetZ(), 2))
	if mod(math.Floor(distance), 2) == 0 {
		return p.color1
	}
	return p.color2
}

// CheckerPattern represents a 3D checker pattern
type CheckerPattern struct {
	*PatternObject
}

// NewCheckerPattern returns a new CheckerPattern object
func NewCheckerPattern(color1, color2 *Color) *CheckerPattern {
	return &CheckerPattern{
		PatternObject: newPattern(color1, color2, checkerAt),
	}
}

// checkerAt returns the checker color at a specific point
func checkerAt(point *base.Tuple, p *PatternObject) *Color {
	sum := math.Floor(point.GetX()) + math.Floor(point.GetY()) + math.Floor(point.GetZ())
	if base.EqualFloats(mod(sum, 2), 0) {
		return p.color1
	}
	return p.color2
}

// properly handles modding with negative numbers (returns positive)
func mod(x, y float64) float64 {
	res := math.Mod(x, y)
	if (res < 0 && y > 0) || (res > 0 && y < 0) {
		return res + y
	}
	return res
}