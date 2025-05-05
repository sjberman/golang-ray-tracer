package image

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
)

func testNewPattern(g *WithT, p Pattern) {
	g.Expect(p.GetColors()[0]).To(Equal(White))
	g.Expect(p.GetColors()[1]).To(Equal(Black))
	g.Expect(p.GetTransform()).To(Equal(&base.Identity))

	t := base.Translate(2, 3, 4)
	p.SetTransform(t)
	g.Expect(p.GetTransform()).To(Equal(t))
}

func TestNewPattern(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	testNewPattern(g, NewPattern(White, Black, nil))
}

func TestNewStripePattern(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	testNewPattern(g, NewStripePattern(White, Black))
}

func TestPatternAt_Stripe(t *testing.T) {
	t.Parallel()

	sp := NewStripePattern(White, Black)

	tests := []struct {
		name  string
		point *base.Tuple
		color *Color
	}{
		{
			name:  "all zeros",
			point: base.NewPoint(0, 0, 0),
			color: White,
		},
		{
			name:  "constant in y; 1",
			point: base.NewPoint(0, 1, 0),
			color: White,
		},
		{
			name:  "constant in y; 2",
			point: base.NewPoint(0, 2, 0),
			color: White,
		},
		{
			name:  "constant in z; 1",
			point: base.NewPoint(0, 0, 1),
			color: White,
		},
		{
			name:  "constant in z; 2",
			point: base.NewPoint(0, 0, 2),
			color: White,
		},
		{
			name:  "alternates in x; 0.9",
			point: base.NewPoint(0.9, 0, 0),
			color: White,
		},
		{
			name:  "alternates in x; 1",
			point: base.NewPoint(1, 0, 0),
			color: Black,
		},
		{
			name:  "alternates in x; -0.1",
			point: base.NewPoint(-0.1, 0, 0),
			color: Black,
		},
		{
			name:  "alternates in x; -1",
			point: base.NewPoint(-1, 0, 0),
			color: Black,
		},
		{
			name:  "alternates in x; -1.1",
			point: base.NewPoint(-1.1, 0, 0),
			color: White,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			stripe := sp.PatternAt(test.point)
			g.Expect(stripe).To(Equal(test.color))
		})
	}
}

func TestNewGradientPattern(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	testNewPattern(g, NewGradientPattern(White, Black))
}

func TestPatternAt_Gradient(t *testing.T) {
	t.Parallel()

	gp := NewGradientPattern(White, Black)

	tests := []struct {
		name  string
		point *base.Tuple
		color *Color
	}{
		{
			name:  "all zeros",
			point: base.NewPoint(0, 0, 0),
			color: White,
		},
		{
			name:  "0.25 in x",
			point: base.NewPoint(0.25, 0, 0),
			color: NewColor(0.75, 0.75, 0.75),
		},
		{
			name:  "0.5 in x",
			point: base.NewPoint(0.5, 0, 0),
			color: NewColor(0.5, 0.5, 0.5),
		},
		{
			name:  "0.75 in x",
			point: base.NewPoint(0.75, 0, 0),
			color: NewColor(0.25, 0.25, 0.25),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			gradient := gp.PatternAt(test.point)
			g.Expect(gradient).To(Equal(test.color))
		})
	}
}

func TestNewRingPattern(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	testNewPattern(g, NewRingPattern(White, Black))
}

func TestPatternAt_Ring(t *testing.T) {
	t.Parallel()

	rp := NewRingPattern(White, Black)

	tests := []struct {
		name  string
		point *base.Tuple
		color *Color
	}{
		{
			name:  "all zeros",
			point: base.NewPoint(0, 0, 0),
			color: White,
		},
		{
			name:  "1 in x",
			point: base.NewPoint(1, 0, 0),
			color: Black,
		},
		{
			name:  "1 in z",
			point: base.NewPoint(0, 0, 1),
			color: Black,
		},
		{
			name:  "sqrt(2)/2",
			point: base.NewPoint(0.708, 0, 0.708),
			color: Black,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			ring := rp.PatternAt(test.point)
			g.Expect(ring).To(Equal(test.color))
		})
	}
}

func TestNewCheckerPattern(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	testNewPattern(g, NewCheckerPattern(White, Black))
}

func TestPatternAt_Checker(t *testing.T) {
	t.Parallel()

	cp := NewCheckerPattern(White, Black)

	tests := []struct {
		name  string
		point *base.Tuple
		color *Color
	}{
		{
			name:  "all zeros",
			point: base.NewPoint(0, 0, 0),
			color: White,
		},
		{
			name:  "repeat in x; 0.99",
			point: base.NewPoint(0.99, 0, 0),
			color: White,
		},
		{
			name:  "repeat in x; 1.01",
			point: base.NewPoint(1.01, 0, 0),
			color: Black,
		},
		{
			name:  "repeat in y; 0.99",
			point: base.NewPoint(0, 0.99, 0),
			color: White,
		},
		{
			name:  "repeat in y; 1.01",
			point: base.NewPoint(0, 1.01, 0),
			color: Black,
		},
		{
			name:  "repeat in z; 0.99",
			point: base.NewPoint(0, 0, 0.99),
			color: White,
		},
		{
			name:  "repeat in z; 1.01",
			point: base.NewPoint(0, 0, 1.01),
			color: Black,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			checker := cp.PatternAt(test.point)
			g.Expect(checker).To(Equal(test.color))
		})
	}
}
