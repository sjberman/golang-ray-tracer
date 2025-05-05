package object

import (
	"math"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

func TestNewCone(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	c := NewCone()
	testNewObject(g, c)
	g.Expect(c.Minimum).To(Equal(math.Inf(-1)))
	g.Expect(c.Maximum).To(Equal(math.Inf(0)))
	g.Expect(c.Closed).To(BeFalse())
	g.Expect(c.Bounds()).To(Equal(&Bounds{
		Minimum: base.NewPoint(-1, c.Minimum, -1),
		Maximum: base.NewPoint(1, c.Maximum, 1),
	}))
}

func TestConeIntersect(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	c := NewCone()
	tests := []struct {
		ray              *ray.Ray
		expVal1, expVal2 float64
	}{
		{
			ray:     ray.NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1)),
			expVal1: 5.0,
			expVal2: 5.0,
		},
		{
			ray:     ray.NewRay(base.NewPoint(0, 0, -5), base.NewVector(1, 1, 1).Normalize()),
			expVal1: 8.660254037844386,
			expVal2: 8.660254037844386,
		},
		{
			ray:     ray.NewRay(base.NewPoint(1, 1, -5), base.NewVector(-0.5, -1, 1).Normalize()),
			expVal1: 4.550055679356349,
			expVal2: 49.449944320643645,
		},
	}

	for _, test := range tests {
		ints := c.Intersect(test.ray)
		g.Expect(len(ints)).To(Equal(2))
		g.Expect(ints[0].Value).To(BeNumerically("~", test.expVal1, 0.00001))
		g.Expect(ints[1].Value).To(BeNumerically("~", test.expVal2, 0.00001))
	}

	// ray parallel to one half
	r := ray.NewRay(base.NewPoint(0, 0, -1), base.NewVector(0, 1, 1).Normalize())
	ints := c.Intersect(r)
	g.Expect(len(ints)).To(Equal(3))
	g.Expect(ints[0].Value).To(Equal(0.0))
	g.Expect(ints[1].Value).To(Equal(.3535533905932738))

	// intersection of caps on a closed cone
	c.Minimum = -0.5
	c.Maximum = 0.5
	c.Closed = true
	r = ray.NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 1, 0))
	g.Expect(len(c.Intersect(r))).To(BeZero())

	r = ray.NewRay(base.NewPoint(0, 0, -0.25), base.NewVector(0, 1, 1).Normalize())
	g.Expect(len(c.Intersect(r))).To(Equal(3))

	r = ray.NewRay(base.NewPoint(0, 0, -0.25), base.NewVector(0, 1, 0))
	g.Expect(len(c.Intersect(r))).To(Equal(4))
}

func TestConeNormalAt(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	c := NewCone()

	n := c.NormalAt(base.NewPoint(0, 0, 0), nil)
	g.Expect(n).To(Equal(base.NewVector(0, 0, 0)))

	n = c.NormalAt(base.NewPoint(1, 1, 1), nil)
	g.Expect(n).To(Equal(base.NewVector(0.5, -math.Sqrt(2)/2, 0.5)))

	n = c.NormalAt(base.NewPoint(-1, -1, 0), nil)
	g.Expect(n).To(Equal(base.NewVector(-1, 1, 0).Normalize()))
}
