package image

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("color tests", func() {
	It("creates colors", func() {
		color := NewColor(1, 2, 3)
		Expect(color.red).To(Equal(1.0))
		Expect(color.green).To(Equal(2.0))
		Expect(color.blue).To(Equal(3.0))
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
		expColor := NewColor(0.9, 0.2, 0.04)
		res := c1.MultiplyColor(c2)
		Expect(res.Equals(expColor)).To(BeTrue())
	})

	It("checks the equivalence of colors", func() {
		c1 := NewColor(1.001, 2, -3.345)
		c2 := NewColor(1.001, 2, -3.345)
		Expect(c1.Equals(c2)).To(BeTrue())

		c1 = NewColor(1.002, 2, -3.346)
		c2 = NewColor(2, 2, -3.346)
		Expect(c1.Equals(c2)).To(BeFalse())

		c1 = NewColor(1.002, 2, -3.346)
		c2 = NewColor(1.002, 2.0001, -3.346)
		Expect(c1.Equals(c2)).To(BeFalse())

		c1 = NewColor(1.002, 2, -3.346)
		c2 = NewColor(1.002, 2, -3.345)
		Expect(c1.Equals(c2)).To(BeFalse())
	})
})
