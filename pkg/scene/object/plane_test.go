package object

import (
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

var _ = Describe("plane tests", func() {
	It("creates planes", func() {
		p := NewPlane()
		testNewObject(p)
		Expect(p.Bounds()).To(Equal(&bounds{
			minimum: base.NewPoint(math.Inf(-1), 0, math.Inf(-1)),
			maximum: base.NewPoint(math.Inf(1), 0, math.Inf(1))}))
	})

	It("calculates a plane intersection", func() {
		// ray parallel to plane
		p := NewPlane()
		r := ray.NewRay(base.NewPoint(0, 10, 0), base.NewVector(0, 0, 1))
		ints := p.Intersect(r)
		Expect(len(ints)).To(BeZero())

		// coplanar ray
		r = ray.NewRay(base.NewPoint(0, 0, 0), base.NewVector(0, 0, 1))
		ints = p.Intersect(r)
		Expect(len(ints)).To(BeZero())

		// intersects plane from above
		r = ray.NewRay(base.NewPoint(0, 1, 0), base.NewVector(0, -1, 0))
		ints = p.Intersect(r)
		Expect(len(ints)).To(Equal(1))
		Expect(ints[0].Value).To(Equal(1.0))
		Expect(ints[0].Object).To(Equal(p))

		// intersects plane from below
		r = ray.NewRay(base.NewPoint(0, -1, 0), base.NewVector(0, 1, 0))
		ints = p.Intersect(r)
		Expect(len(ints)).To(Equal(1))
		Expect(ints[0].Value).To(Equal(1.0))
		Expect(ints[0].Object).To(Equal(p))
	})

	It("computes the surface normal", func() {
		p := NewPlane()
		constVector := base.NewVector(0, 1, 0)

		n := p.NormalAt(base.NewPoint(0, 0, 0))
		Expect(n).To(Equal(constVector))
		n = p.NormalAt(base.NewPoint(10, 0, -10))
		Expect(n).To(Equal(constVector))
		n = p.NormalAt(base.NewPoint(-5, 0, 150))
		Expect(n).To(Equal(constVector))
	})
})
