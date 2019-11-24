package ray

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sjberman/golang-ray-tracer/pkg/base"
)

func TestRay(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ray Suite")
}

var _ = Describe("ray tests", func() {
	It("creates rays", func() {
		o := base.NewPoint(1, 2, 3)
		d := base.NewVector(4, 5, 6)
		ray := NewRay(o, d)
		Expect(ray.origin).To(Equal(o))
		Expect(ray.direction).To(Equal(d))
	})

	It("creates spheres", func() {
		s := NewSphere()
		Expect(s.transform).To(Equal(base.Identity))
		t := base.TranslationMatrix(2, 3, 4)
		s.SetTransform(t)
		Expect(s.transform).To(Equal(t))
	})

	It("creates intersections", func() {
		s := NewSphere()
		i := NewIntersection(3.5, s)
		Expect(i.value).To(Equal(3.5))
		Expect(i.object).To(Equal(s))
	})

	It("computes the position of a point on the ray", func() {
		r := NewRay(base.NewPoint(2, 3, 4), base.NewVector(1, 0, 0))
		expPoint := base.NewPoint(2, 3, 4)
		Expect(r.Position(0)).To(Equal(expPoint))

		expPoint = base.NewPoint(3, 3, 4)
		Expect(r.Position(1)).To(Equal(expPoint))

		expPoint = base.NewPoint(1, 3, 4)
		Expect(r.Position(-1)).To(Equal(expPoint))

		expPoint = base.NewPoint(4.5, 3, 4)
		Expect(r.Position(2.5)).To(Equal(expPoint))
	})

	It("computes the intersection points of a ray and sphere", func() {
		r := NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
		s := NewSphere()

		ints := r.Intersect(s)
		Expect(len(ints)).To(Equal(2))
		Expect(ints[0].value).To(Equal(4.0))
		Expect(ints[1].value).To(Equal(6.0))
		Expect(ints[0].object).To(Equal(s))
		Expect(ints[1].object).To(Equal(s))

		// tangent
		r = NewRay(base.NewPoint(0, 1, -5), base.NewVector(0, 0, 1))
		s = NewSphere()

		ints = r.Intersect(s)
		Expect(len(ints)).To(Equal(2))
		Expect(ints[0].value).To(Equal(5.0))
		Expect(ints[1].value).To(Equal(5.0))

		// too high, no intersection
		r = NewRay(base.NewPoint(0, 2, -5), base.NewVector(0, 0, 1))
		s = NewSphere()

		ints = r.Intersect(s)
		Expect(len(ints)).To(Equal(0))

		// ray starts within the sphere
		r = NewRay(base.NewPoint(0, 0, 0), base.NewVector(0, 0, 1))
		s = NewSphere()

		ints = r.Intersect(s)
		Expect(len(ints)).To(Equal(2))
		Expect(ints[0].value).To(Equal(-1.0))
		Expect(ints[1].value).To(Equal(1.0))

		// ray starts past the sphere
		r = NewRay(base.NewPoint(0, 0, 5), base.NewVector(0, 0, 1))
		s = NewSphere()

		ints = r.Intersect(s)
		Expect(len(ints)).To(Equal(2))
		Expect(ints[0].value).To(Equal(-6.0))
		Expect(ints[1].value).To(Equal(-4.0))

		// intersect a scaled sphere
		r = NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
		s = NewSphere()
		s.SetTransform(base.ScalingMatrix(2, 2, 2))
		ints = r.Intersect(s)
		Expect(len(ints)).To(Equal(2))
		Expect(ints[0].value).To(Equal(3.0))
		Expect(ints[1].value).To(Equal(7.0))

		r = NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
		s = NewSphere()
		s.SetTransform(base.TranslationMatrix(5, 0, 0))
		ints = r.Intersect(s)
		Expect(len(ints)).To(Equal(0))
	})

	It("returns a list of intersections", func() {
		s := NewSphere()
		i1 := NewIntersection(1, s)
		i2 := NewIntersection(2, s)
		ints := intersections(i1, i2)
		Expect(len(ints)).To(Equal(2))
		Expect(ints[0].value).To(Equal(1.0))
		Expect(ints[1].value).To(Equal(2.0))
	})

	It("identifies a hit", func() {
		s := NewSphere()
		i1 := NewIntersection(1, s)
		i2 := NewIntersection(2, s)
		ints := intersections(i1, i2)
		i := Hit(ints)
		Expect(i).To(Equal(i1))

		i1 = NewIntersection(-1, s)
		i2 = NewIntersection(1, s)
		ints = intersections(i1, i2)
		i = Hit(ints)
		Expect(i).To(Equal(i2))

		i1 = NewIntersection(-2, s)
		i2 = NewIntersection(-1, s)
		ints = intersections(i1, i2)
		i = Hit(ints)
		Expect(i).To(BeNil())

		i1 = NewIntersection(5, s)
		i2 = NewIntersection(7, s)
		i3 := NewIntersection(-3, s)
		i4 := NewIntersection(2, s)
		ints = intersections(i1, i2, i3, i4)
		i = Hit(ints)
		Expect(i).To(Equal(i4))
	})

	It("transforms a ray", func() {
		r := NewRay(base.NewPoint(1, 2, 3), base.NewVector(0, 1, 0))
		m := base.TranslationMatrix(3, 4, 5)
		r2 := r.Transform(m)
		Expect(r2.origin).To(Equal(base.NewPoint(4, 6, 8)))
		Expect(r2.direction).To(Equal(base.NewVector(0, 1, 0)))

		m = base.ScalingMatrix(2, 3, 4)
		r2 = r.Transform(m)
		Expect(r2.origin).To(Equal(base.NewPoint(2, 6, 12)))
		Expect(r2.direction).To(Equal(base.NewVector(0, 3, 0)))
	})
})
