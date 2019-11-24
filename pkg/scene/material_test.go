package scene

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/image"
)

var _ = Describe("material tests", func() {
	It("creates materials", func() {
		m := NewMaterial(image.NewColor(1, 1, 1), 1, 2, 3, 100)
		Expect(m.color).To(Equal(image.NewColor(1, 1, 1)))
		Expect(m.ambient).To(Equal(1.0))
		Expect(m.diffuse).To(Equal(2.0))
		Expect(m.specular).To(Equal(3.0))
		Expect(m.shininess).To(Equal(100.0))
	})

	It("sets fields", func() {
		m := NewMaterial(image.NewColor(1, 1, 1), 1, 2, 3, 100)

		m.SetColor(image.NewColor(0, 0, 0))
		Expect(m.color).To(Equal(image.NewColor(0, 0, 0)))

		m.SetAmbient(5)
		Expect(m.ambient).To(Equal(5.0))

		m.SetDiffuse(7)
		Expect(m.diffuse).To(Equal(7.0))

		m.SetSpecular(10)
		Expect(m.specular).To(Equal(10.0))

		m.SetShininess(12)
		Expect(m.shininess).To(Equal(12.0))
	})
})
