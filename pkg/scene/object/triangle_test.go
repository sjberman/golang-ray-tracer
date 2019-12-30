package object

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

var _ = Describe("triangle tests", func() {
	It("creates triangles", func() {
		p1 := base.NewPoint(0, 1, 0)
		p2 := base.NewPoint(-1, 0, 0)
		p3 := base.NewPoint(1, 0, 0)

		t := NewTriangle(p1, p2, p3)
		testNewObject(t)

		Expect(t.P1).To(Equal(p1))
		Expect(t.P2).To(Equal(p2))
		Expect(t.P3).To(Equal(p3))
		Expect(t.edge1).To(Equal(base.NewVector(-1, -1, 0)))
		Expect(t.edge2).To(Equal(base.NewVector(1, -1, 0)))
		Expect(t.normal).To(Equal(base.NewVector(0, 0, -1)))

		Expect(t.Bounds()).To(Equal(&Bounds{
			Minimum: base.NewPoint(-1, 0, 0),
			Maximum: base.NewPoint(1, 1, 0)}))
	})

	It("calculates a triangle intersection", func() {
		// ray is parallel (misses)
		t := NewTriangle(base.NewPoint(0, 1, 0), base.NewPoint(-1, 0, 0), base.NewPoint(1, 0, 0))
		r := ray.NewRay(base.NewPoint(0, -1, -2), base.NewVector(0, 1, 0))
		ints := t.Intersect(r)
		Expect(len(ints)).To(BeZero())

		// ray misses the P1-P3 edge
		r = ray.NewRay(base.NewPoint(1, 1, -2), base.NewVector(0, 0, 1))
		ints = t.Intersect(r)
		Expect(len(ints)).To(BeZero())

		// ray misses the P1-P2 edge
		r = ray.NewRay(base.NewPoint(-1, 1, -2), base.NewVector(0, 0, 1))
		ints = t.Intersect(r)
		Expect(len(ints)).To(BeZero())

		// ray misses the P2-P3 edge
		r = ray.NewRay(base.NewPoint(0, -1, -2), base.NewVector(0, 0, 1))
		ints = t.Intersect(r)
		Expect(len(ints)).To(BeZero())

		// hit
		r = ray.NewRay(base.NewPoint(0, 0.5, -2), base.NewVector(0, 0, 1))
		ints = t.Intersect(r)
		Expect(len(ints)).To(Equal(1))
		Expect(ints[0].Value).To(Equal(2.0))
	})

	It("computes the surface normal", func() {
		t := NewTriangle(base.NewPoint(0, 1, 0), base.NewPoint(-1, 0, 0), base.NewPoint(1, 0, 0))
		n1 := t.NormalAt(base.NewPoint(0, 0.5, 0), nil)
		n2 := t.NormalAt(base.NewPoint(-0.5, 0.75, 0), nil)
		n3 := t.NormalAt(base.NewPoint(0.5, 0.25, 0), nil)

		Expect(n1).To(Equal(t.normal))
		Expect(n2).To(Equal(t.normal))
		Expect(n3).To(Equal(t.normal))
	})

	Context("smooth triangles", func() {
		var tri *SmoothTriangle
		var p1, p2, p3, n1, n2, n3 *base.Tuple

		BeforeEach(func() {
			p1 = base.NewPoint(0, 1, 0)
			p2 = base.NewPoint(-1, 0, 0)
			p3 = base.NewPoint(1, 0, 0)
			n1 = base.NewVector(0, 1, 0)
			n2 = base.NewVector(-1, 0, 0)
			n3 = base.NewVector(1, 0, 0)
			tri = NewSmoothTriangle(p1, p2, p3, n1, n2, n3)
		})

		It("creates smooth triangles", func() {
			Expect(tri.P1).To(Equal(p1))
			Expect(tri.P2).To(Equal(p2))
			Expect(tri.P3).To(Equal(p3))
			Expect(tri.N1).To(Equal(n1))
			Expect(tri.N2).To(Equal(n2))
			Expect(tri.N3).To(Equal(n3))
		})

		It("uses u/v to interpolate the normal", func() {
			intersection := NewIntersection(1, tri)
			intersection.u = 0.45
			intersection.v = 0.25

			n := tri.NormalAt(base.NewPoint(0, 0, 0), intersection)
			Expect(n).To(Equal(base.NewVector(-0.5547001962252291, 0.8320502943378437, 0)))
		})
	})
})
