package scene

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/image"
)

var _ = Describe("material tests", func() {
	It("creates materials", func() {
		pattern := image.NewStripePattern(image.White, image.Black)
		m := NewMaterial(image.White, pattern, 1, 2, 3, 100, 4, 5, 6)
		Expect(m.color).To(Equal(image.White))
		Expect(m.pattern).To(Equal(pattern))
		Expect(m.ambient).To(Equal(1.0))
		Expect(m.diffuse).To(Equal(2.0))
		Expect(m.specular).To(Equal(3.0))
		Expect(m.shininess).To(Equal(100.0))
		Expect(m.reflective).To(Equal(4.0))
		Expect(m.transparency).To(Equal(5.0))
		Expect(m.refractiveIndex).To(Equal(6.0))
	})

	It("sets fields", func() {
		pattern := image.NewStripePattern(image.White, image.Black)
		m := NewMaterial(image.White, pattern, 1, 2, 3, 100, 4, 5, 6)

		m.SetColor(image.Black)
		Expect(m.color).To(Equal(image.Black))

		newPattern := image.NewStripePattern(image.Black, image.White)
		m.SetPattern(newPattern)
		Expect(m.pattern).To(Equal(newPattern))

		m.SetAmbient(5)
		Expect(m.ambient).To(Equal(5.0))

		m.SetDiffuse(7)
		Expect(m.diffuse).To(Equal(7.0))

		m.SetSpecular(10)
		Expect(m.specular).To(Equal(10.0))

		m.SetShininess(12)
		Expect(m.shininess).To(Equal(12.0))

		m.SetReflective(15)
		Expect(m.reflective).To(Equal(15.0))

		m.SetTransparency(18)
		Expect(m.transparency).To(Equal(18.0))

		m.SetRefractiveIndex(20)
		Expect(m.refractiveIndex).To(Equal(20.0))
	})
})
