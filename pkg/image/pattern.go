package image

import (
	"math"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	grtMath "github.com/sjberman/golang-ray-tracer/pkg/math"
)

// Pattern represents a color pattern.
type Pattern interface {
	GetColors() []*Color
	GetTransform() *base.Matrix
	SetTransform(...*base.Matrix)
	PatternAt(*base.Tuple) *Color
}

// PatternObject is the base implementation of the pattern interface.
type PatternObject struct {
	color1        *Color
	color2        *Color
	transform     base.Matrix
	patternAtFunc func(*base.Tuple, *PatternObject) *Color
}

// NewPattern returns a new pattern object.
func NewPattern(
	color1, color2 *Color,
	patternFunc func(*base.Tuple, *PatternObject) *Color,
) *PatternObject {
	return &PatternObject{
		color1:        color1,
		color2:        color2,
		patternAtFunc: patternFunc,
		transform:     base.Identity,
	}
}

// GetColors returns a list of colors for the pattern.
func (p *PatternObject) GetColors() []*Color {
	return []*Color{p.color1, p.color2}
}

// GetTransform returns the pattern's transform matrix.
func (p *PatternObject) GetTransform() *base.Matrix {
	return &p.transform
}

// SetTransform sets the pattern's transform matrix.
func (p *PatternObject) SetTransform(matrix ...*base.Matrix) {
	t := base.Identity
	for _, m := range matrix {
		t = *t.Multiply(m)
	}
	p.transform = t
}

// PatternAt returns the color at a specific point, based on the pattern.
func (p *PatternObject) PatternAt(point *base.Tuple) *Color {
	return p.patternAtFunc(point, p)
}

// StripePattern represents a striped color pattern.
type StripePattern struct {
	*PatternObject
}

// NewStripePattern returns a new StripePattern object.
func NewStripePattern(color1, color2 *Color) *StripePattern {
	return &StripePattern{
		PatternObject: NewPattern(color1, color2, stripeAt),
	}
}

// stripeAt returns the stripe color at a specific point.
func stripeAt(point *base.Tuple, p *PatternObject) *Color {
	m := grtMath.Mod(point.GetX(), 2)
	if m >= 0 && m < 1 {
		return p.color1
	}

	return p.color2
}

// GradientPattern represents a gradient color pattern.
type GradientPattern struct {
	*PatternObject
}

// NewGradientPattern returns a new GradientPattern object.
func NewGradientPattern(color1, color2 *Color) *GradientPattern {
	return &GradientPattern{
		PatternObject: NewPattern(color1, color2, gradientAt),
	}
}

// gradientAt returns the gradient color at a specific point.
func gradientAt(point *base.Tuple, p *PatternObject) *Color {
	colors := p.GetColors()
	distance := colors[1].Subtract(colors[0])
	fraction := point.GetX() - math.Floor(point.GetX())

	return colors[0].Add(distance.Multiply(fraction))
}

// RingPattern represents a ring color pattern.
type RingPattern struct {
	*PatternObject
}

// NewRingPattern returns a new RingPattern object.
func NewRingPattern(color1, color2 *Color) *RingPattern {
	return &RingPattern{
		PatternObject: NewPattern(color1, color2, ringAt),
	}
}

// ringAt returns the ring color at a specific point.
func ringAt(point *base.Tuple, p *PatternObject) *Color {
	distance := math.Sqrt(math.Pow(point.GetX(), 2) + math.Pow(point.GetZ(), 2))
	if grtMath.Mod(math.Floor(distance), 2) == 0 {
		return p.color1
	}

	return p.color2
}

// CheckerPattern represents a 3D checker pattern.
type CheckerPattern struct {
	*PatternObject
}

// NewCheckerPattern returns a new CheckerPattern object.
func NewCheckerPattern(color1, color2 *Color) *CheckerPattern {
	return &CheckerPattern{
		PatternObject: NewPattern(color1, color2, checkerAt),
	}
}

// checkerAt returns the checker color at a specific point.
func checkerAt(point *base.Tuple, p *PatternObject) *Color {
	sum := math.Floor(point.GetX()) + math.Floor(point.GetY()) + math.Floor(point.GetZ())
	if base.EqualFloats(grtMath.Mod(sum, 2), 0) {
		return p.color1
	}

	return p.color2
}

// MockPattern is a mock pattern object for unit testing.
type MockPattern struct {
	*PatternObject
}

// NewMockPattern returns a new MockPattern object.
func NewMockPattern() *MockPattern {
	return &MockPattern{
		PatternObject: NewPattern(nil, nil, mockFunc),
	}
}

func mockFunc(point *base.Tuple, p *PatternObject) *Color {
	return NewColor(point.GetX(), point.GetY(), point.GetZ())
}
