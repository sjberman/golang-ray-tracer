package scene

import (
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
)

var _ = Describe("sphere tests", func() {
	It("creates spheres", func() {
		s := NewSphere()
		Expect(s.transform).To(Equal(base.Identity))
		Expect(s.GetMaterial()).To(Equal(&defaultMaterial))

		t := base.TranslationMatrix(2, 3, 4)
		s.SetTransform(t)
		Expect(s.transform).To(Equal(*t))

		m := defaultMaterial
		m.ambient = 1
		s.SetMaterial(&m)
		Expect(s.GetMaterial()).To(Equal(&m))
	})

	It("computes the surface normal", func() {
		// x axis
		s := NewSphere()
		n := s.NormalAt(base.NewPoint(1, 0, 0))
		Expect(n).To(Equal(base.NewVector(1, 0, 0)))

		// y axis
		n = s.NormalAt(base.NewPoint(0, 1, 0))
		Expect(n).To(Equal(base.NewVector(0, 1, 0)))

		// z axis
		n = s.NormalAt(base.NewPoint(0, 0, 1))
		Expect(n).To(Equal(base.NewVector(0, 0, 1)))

		// non axis
		n = s.NormalAt(base.NewPoint(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3))
		Expect(n).To(Equal(base.NewVector(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3)))

		// surface normal is a normalized vector
		Expect(n).To(Equal(n.Normalize()))

		// translated sphere
		s.SetTransform(base.TranslationMatrix(0, 1, 0))
		n = s.NormalAt(base.NewPoint(0, 1.70711, -0.70711))
		Expect(n).To(Equal(base.NewVector(0, 0.7071067811865475, -0.7071067811865476)))

		// scaled/rotated sphere
		m := base.ScalingMatrix(1, 0.5, 1).Multiply(base.ZRotationMatrix(math.Pi / 5))
		s.SetTransform(m)
		n = s.NormalAt(base.NewPoint(0, math.Sqrt(2)/2, -math.Sqrt(2)/2))
		Expect(n).To(Equal(base.NewVector(0, 0.9701425001453319, -0.24253562503633294)))
	})
})
