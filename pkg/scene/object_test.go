package scene

import (
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
)

var _ = Describe("object tests", func() {
	It("creates objects", func() {
		o := newObject(nil, nil)
		Expect(o.GetTransform()).To(Equal(&base.Identity))
		Expect(o.GetMaterial()).To(Equal(&defaultMaterial))

		t := base.TranslationMatrix(2, 3, 4)
		o.SetTransform(t)
		Expect(o.transform).To(Equal(*t))

		m := defaultMaterial
		m.ambient = 1
		o.SetMaterial(&m)
		Expect(o.GetMaterial()).To(Equal(&m))
	})

	Context("spheres", func() {
		It("creates spheres", func() {
			s := NewSphere()
			Expect(s.GetTransform()).To(Equal(&base.Identity))
			Expect(s.GetMaterial()).To(Equal(&defaultMaterial))

			t := base.TranslationMatrix(2, 3, 4)
			s.SetTransform(t)
			Expect(s.transform).To(Equal(*t))

			m := defaultMaterial
			m.ambient = 1
			s.SetMaterial(&m)
			Expect(s.GetMaterial()).To(Equal(&m))
		})

		It("calculates a sphere intersection", func() {
			r := NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
			s := NewSphere()

			ints := s.intersect(r)
			Expect(len(ints)).To(Equal(2))
			Expect(ints[0].value).To(Equal(4.0))
			Expect(ints[1].value).To(Equal(6.0))
			Expect(ints[0].GetObject()).To(Equal(s.object))
			Expect(ints[1].GetObject()).To(Equal(s.object))

			// tangent
			r = NewRay(base.NewPoint(0, 1, -5), base.NewVector(0, 0, 1))
			s = NewSphere()

			ints = s.intersect(r)
			Expect(len(ints)).To(Equal(2))
			Expect(ints[0].value).To(Equal(5.0))
			Expect(ints[1].value).To(Equal(5.0))

			// too high, no  intersection
			r = NewRay(base.NewPoint(0, 2, -5), base.NewVector(0, 0, 1))
			s = NewSphere()

			ints = s.intersect(r)
			Expect(len(ints)).To(Equal(0))

			// ray starts within the sphere
			r = NewRay(base.Origin, base.NewVector(0, 0, 1))
			s = NewSphere()

			ints = s.intersect(r)
			Expect(len(ints)).To(Equal(2))
			Expect(ints[0].value).To(Equal(-1.0))
			Expect(ints[1].value).To(Equal(1.0))

			// ray starts past the sphere
			r = NewRay(base.NewPoint(0, 0, 5), base.NewVector(0, 0, 1))
			s = NewSphere()

			ints = s.intersect(r)
			Expect(len(ints)).To(Equal(2))
			Expect(ints[0].value).To(Equal(-6.0))
			Expect(ints[1].value).To(Equal(-4.0))

			//  intersect a scaled sphere
			r = NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
			s = NewSphere()
			s.SetTransform(base.ScalingMatrix(2, 2, 2))
			ints = s.intersect(r)
			Expect(len(ints)).To(Equal(2))
			Expect(ints[0].value).To(Equal(3.0))
			Expect(ints[1].value).To(Equal(7.0))

			r = NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
			s = NewSphere()
			s.SetTransform(base.TranslationMatrix(5, 0, 0))
			ints = s.intersect(r)
			Expect(len(ints)).To(Equal(0))
		})

		It("computes the surface normal", func() {
			// x axis
			s := NewSphere()
			n := s.normalAt(base.NewPoint(1, 0, 0))
			Expect(n).To(Equal(base.NewVector(1, 0, 0)))

			// y axis
			n = s.normalAt(base.NewPoint(0, 1, 0))
			Expect(n).To(Equal(base.NewVector(0, 1, 0)))

			// z axis
			n = s.normalAt(base.NewPoint(0, 0, 1))
			Expect(n).To(Equal(base.NewVector(0, 0, 1)))

			// non axis
			n = s.normalAt(base.NewPoint(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3))
			Expect(n).To(Equal(base.NewVector(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3)))

			// surface normal is a normalized vector
			Expect(n).To(Equal(n.Normalize()))

			// translated sphere
			s.SetTransform(base.TranslationMatrix(0, 1, 0))
			n = s.normalAt(base.NewPoint(0, 1.70711, -0.70711))
			Expect(n).To(Equal(base.NewVector(0, 0.7071067811865475, -0.7071067811865476)))

			// scaled/rotated sphere
			m := base.ScalingMatrix(1, 0.5, 1).Multiply(base.ZRotationMatrix(math.Pi / 5))
			s.SetTransform(m)
			n = s.normalAt(base.NewPoint(0, math.Sqrt(2)/2, -math.Sqrt(2)/2))
			Expect(n).To(Equal(base.NewVector(0, 0.9701425001453319, -0.24253562503633294)))
		})
	})
})
