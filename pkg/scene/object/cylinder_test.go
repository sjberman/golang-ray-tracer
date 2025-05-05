package object

import (
	"math"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

func TestNewCylinder(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	c := NewCylinder()
	testNewObject(g, c)
	g.Expect(c.Minimum).To(Equal(math.Inf(-1)))
	g.Expect(c.Maximum).To(Equal(math.Inf(1)))
	g.Expect(c.Closed).To(BeFalse())
	g.Expect(c.Bounds()).To(Equal(&Bounds{
		Minimum: base.NewPoint(-1, c.Minimum, -1),
		Maximum: base.NewPoint(1, c.Maximum, 1),
	}))
}

func TestCylinderIntersect(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// misses
	c := NewCylinder()
	r := ray.NewRay(base.NewPoint(1, 0, 0), base.NewVector(0, 1, 0))
	ints := c.Intersect(r)
	g.Expect(len(ints)).To(BeZero())

	r = ray.NewRay(base.NewPoint(0, 0, 0), base.NewVector(0, 1, 0))
	ints = c.Intersect(r)
	g.Expect(len(ints)).To(BeZero())

	r = ray.NewRay(base.NewPoint(0, 0, -5), base.NewVector(1, 1, 1))
	ints = c.Intersect(r)
	g.Expect(len(ints)).To(BeZero())

	// hits
	tests := []struct {
		name             string
		ray              *ray.Ray
		expVal1, expVal2 float64
	}{
		{
			ray:     ray.NewRay(base.NewPoint(1, 0, -5), base.NewVector(0, 0, 1)),
			expVal1: 5.0,
			expVal2: 5.0,
		},
		{
			ray:     ray.NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1)),
			expVal1: 4.0,
			expVal2: 6.0,
		},
		{
			ray:     ray.NewRay(base.NewPoint(0.5, 0, -5), base.NewVector(0.1, 1, 1).Normalize()),
			expVal1: 6.80798191702732,
			expVal2: 7.088723439378861,
		},
	}

	for _, test := range tests {
		ints := c.Intersect(test.ray)
		g.Expect(len(ints)).To(Equal(2))
		g.Expect(ints[0].Value).To(BeNumerically("~", test.expVal1))
		g.Expect(ints[1].Value).To(BeNumerically("~", test.expVal2))
	}

	// constrained cylinder
	c.Minimum = 1
	c.Maximum = 2
	rays := []*ray.Ray{
		ray.NewRay(base.NewPoint(0, 1.5, 0), base.NewVector(0.1, 1, 0).Normalize()),
		ray.NewRay(base.NewPoint(0, 3, -5), base.NewVector(0, 0, 1)),
		ray.NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1)),
		ray.NewRay(base.NewPoint(0, 2, -5), base.NewVector(0, 0, 1)),
		ray.NewRay(base.NewPoint(0, 1, -5), base.NewVector(0, 0, 1)),
	}

	for _, r := range rays {
		ints := c.Intersect(r)
		g.Expect(len(ints)).To(BeZero())
	}
	r = ray.NewRay(base.NewPoint(0, 1.5, -2), base.NewVector(0, 0, 1))
	g.Expect(len(c.Intersect(r))).To(Equal(2))

	// intersection of caps on a closed cylinder
	c.Closed = true
	rays = []*ray.Ray{
		ray.NewRay(base.NewPoint(0, 3, 0), base.NewVector(0, -1, 0)),
		ray.NewRay(base.NewPoint(0, 3, -2), base.NewVector(0, -1, 2)),
		ray.NewRay(base.NewPoint(0, 4, -2), base.NewVector(0, -1, 1)),
		ray.NewRay(base.NewPoint(0, 0, -2), base.NewVector(0, 1, 2)),
		ray.NewRay(base.NewPoint(0, -1, -2), base.NewVector(0, 1, 1)),
	}

	for _, r := range rays {
		ints := c.Intersect(r)
		g.Expect(len(ints)).To(Equal(2))
	}
}

func TestCylinderNormalAt(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	c := NewCylinder()

	n := c.NormalAt(base.NewPoint(1, 0, 0), nil)
	g.Expect(n).To(Equal(base.NewVector(1, 0, 0)))

	n = c.NormalAt(base.NewPoint(0, 5, -1), nil)
	g.Expect(n).To(Equal(base.NewVector(0, 0, -1)))

	n = c.NormalAt(base.NewPoint(0, -2, 1), nil)
	g.Expect(n).To(Equal(base.NewVector(0, 0, 1)))

	n = c.NormalAt(base.NewPoint(-1, 1, 0), nil)
	g.Expect(n).To(Equal(base.NewVector(-1, 0, 0)))

	// cylinder caps
	c.Minimum = 1
	c.Maximum = 2
	c.Closed = true
	n = c.NormalAt(base.NewPoint(0, 1, 0), nil)
	g.Expect(n).To(Equal(base.NewVector(0, -1, 0)))

	n = c.NormalAt(base.NewPoint(0.5, 1, 0), nil)
	g.Expect(n).To(Equal(base.NewVector(0, -1, 0)))

	n = c.NormalAt(base.NewPoint(0, 1, 0.5), nil)
	g.Expect(n).To(Equal(base.NewVector(0, -1, 0)))

	n = c.NormalAt(base.NewPoint(0, 2, 0), nil)
	g.Expect(n).To(Equal(base.NewVector(0, 1, 0)))

	n = c.NormalAt(base.NewPoint(0.5, 2, 0), nil)
	g.Expect(n).To(Equal(base.NewVector(0, 1, 0)))

	n = c.NormalAt(base.NewPoint(0, 2, 0.5), nil)
	g.Expect(n).To(Equal(base.NewVector(0, 1, 0)))
}
