package image

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("color tests", func() {
	It("creates colors", func() {
		color := NewColor(1, 2, 3)
		Expect(color.red).To(Equal(float64(1)))
		Expect(color.green).To(Equal(float64(2)))
		Expect(color.blue).To(Equal(float64(3)))
	})

	It("adds colors", func() {
		c1 := NewColor(0.9, 0.6, 0.75)
		c2 := NewColor(0.7, 0.1, 0.25)
		expColor := NewColor(1.6, 0.7, 1.0)
		Expect(c1.Add(c2)).To(Equal(expColor))
	})

	It("subtracts colors", func() {
		c1 := NewColor(0.9, 0.6, 0.75)
		c2 := NewColor(0.7, 0.1, 0.25)
		// Floating point precision...
		expColor := NewColor(0.20000000000000007, 0.5, 0.5)
		Expect(c1.Subtract(c2)).To(Equal(expColor))
	})

	It("multiplies colors by values", func() {
		val := 2.0
		c := NewColor(0.2, 0.3, 0.4)
		expColor := NewColor(0.4, 0.6, 0.8)
		Expect(c.Multiply(val)).To(Equal(expColor))
	})

	It("multiplies colors by colors", func() {
		c1 := NewColor(1, 0.2, 0.4)
		c2 := NewColor(0.9, 1, 0.1)
		// Floating point precision...
		expColor := NewColor(0.9, 0.2, 0.04000000000000001)
		Expect(c1.MultiplyColor(c2)).To(Equal(expColor))
	})
})
