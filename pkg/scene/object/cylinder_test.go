package object

import (
	"fmt"
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

var _ = Describe("cylinder tests", func() {
	It("creates cylinders", func() {
		c := NewCylinder()
		testNewObject(c)
		Expect(c.Minimum).To(Equal(math.Inf(-1)))
		Expect(c.Maximum).To(Equal(math.Inf(0)))
		Expect(c.Closed).To(BeFalse())
	})

	It("calculates a cylinder intersection", func() {
		// misses
		c := NewCylinder()
		r := ray.NewRay(base.NewPoint(1, 0, 0), base.NewVector(0, 1, 0))
		ints := c.Intersect(r)
		Expect(len(ints)).To(BeZero())

		r = ray.NewRay(base.NewPoint(0, 0, 0), base.NewVector(0, 1, 0))
		ints = c.Intersect(r)
		Expect(len(ints)).To(BeZero())

		r = ray.NewRay(base.NewPoint(0, 0, -5), base.NewVector(1, 1, 1))
		ints = c.Intersect(r)
		Expect(len(ints)).To(BeZero())

		// hits
		testCases := []data{
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

		for _, t := range testCases {
			ints := c.Intersect(t.ray)
			Expect(len(ints)).To(Equal(2),
				fmt.Sprintf("origin: %v, direction: %v", t.ray.Origin, t.ray.Direction))
			Expect(ints[0].Value).To(Equal(t.expVal1))
			Expect(ints[1].Value).To(Equal(t.expVal2))
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
			Expect(len(ints)).To(BeZero(),
				fmt.Sprintf("origin: %v, direction: %v", r.Origin, r.Direction))
		}
		r = ray.NewRay(base.NewPoint(0, 1.5, -2), base.NewVector(0, 0, 1))
		Expect(len(c.Intersect(r))).To(Equal(2))

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
			Expect(len(ints)).To(Equal(2),
				fmt.Sprintf("origin: %v, direction: %v", r.Origin, r.Direction))
		}
	})

	It("computes the surface normal", func() {
		c := NewCylinder()

		n := c.NormalAt(base.NewPoint(1, 0, 0))
		Expect(n).To(Equal(base.NewVector(1, 0, 0)))

		n = c.NormalAt(base.NewPoint(0, 5, -1))
		Expect(n).To(Equal(base.NewVector(0, 0, -1)))

		n = c.NormalAt(base.NewPoint(0, -2, 1))
		Expect(n).To(Equal(base.NewVector(0, 0, 1)))

		n = c.NormalAt(base.NewPoint(-1, 1, 0))
		Expect(n).To(Equal(base.NewVector(-1, 0, 0)))

		// cylinder caps
		c.Minimum = 1
		c.Maximum = 2
		c.Closed = true
		n = c.NormalAt(base.NewPoint(0, 1, 0))
		Expect(n).To(Equal(base.NewVector(0, -1, 0)))

		n = c.NormalAt(base.NewPoint(0.5, 1, 0))
		Expect(n).To(Equal(base.NewVector(0, -1, 0)))

		n = c.NormalAt(base.NewPoint(0, 1, 0.5))
		Expect(n).To(Equal(base.NewVector(0, -1, 0)))

		n = c.NormalAt(base.NewPoint(0, 2, 0))
		Expect(n).To(Equal(base.NewVector(0, 1, 0)))

		n = c.NormalAt(base.NewPoint(0.5, 2, 0))
		Expect(n).To(Equal(base.NewVector(0, 1, 0)))

		n = c.NormalAt(base.NewPoint(0, 2, 0.5))
		Expect(n).To(Equal(base.NewVector(0, 1, 0)))
	})
})
