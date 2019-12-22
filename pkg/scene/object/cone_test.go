package object

import (
	"fmt"
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

var _ = Describe("cone tests", func() {
	It("creates cones", func() {
		c := NewCone()
		testNewObject(c)
		Expect(c.Minimum).To(Equal(math.Inf(-1)))
		Expect(c.Maximum).To(Equal(math.Inf(0)))
		Expect(c.Closed).To(BeFalse())
	})

	It("calculates a cone intersection", func() {
		c := NewCone()
		testCases := []data{
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

		for _, t := range testCases {
			ints := c.Intersect(t.ray)
			Expect(len(ints)).To(Equal(2),
				fmt.Sprintf("origin: %v, direction: %v", t.ray.Origin, t.ray.Direction))
			Expect(ints[0].Value).To(Equal(t.expVal1))
			Expect(ints[1].Value).To(Equal(t.expVal2))
		}

		// ray parallel to one half
		r := ray.NewRay(base.NewPoint(0, 0, -1), base.NewVector(0, 1, 1).Normalize())
		ints := c.Intersect(r)
		Expect(len(ints)).To(Equal(1))
		Expect(ints[0].Value).To(Equal(.3535533905932738))

		// intersection of caps on a closed cone
		c.Minimum = -0.5
		c.Maximum = 0.5
		c.Closed = true
		r = ray.NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 1, 0))
		Expect(len(c.Intersect(r))).To(BeZero())

		r = ray.NewRay(base.NewPoint(0, 0, -0.25), base.NewVector(0, 1, 1).Normalize())
		Expect(len(c.Intersect(r))).To(Equal(2))

		r = ray.NewRay(base.NewPoint(0, 0, -0.25), base.NewVector(0, 1, 0))
		Expect(len(c.Intersect(r))).To(Equal(4))
	})

	It("computes the surface normal", func() {
		c := NewCone()

		n := c.NormalAt(base.NewPoint(0, 0, 0))
		Expect(n).To(Equal(base.NewVector(0, 0, 0)))

		n = c.NormalAt(base.NewPoint(1, 1, 1))
		Expect(n).To(Equal(base.NewVector(0.5, -math.Sqrt(2)/2, 0.5)))

		n = c.NormalAt(base.NewPoint(-1, -1, 0))
		Expect(n).To(Equal(base.NewVector(-1, 1, 0).Normalize()))
	})
})
