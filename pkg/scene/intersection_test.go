package scene

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("intersection tests", func() {
	It("creates intersections", func() {
		s := NewSphere()
		i := NewIntersection(3.5, s)
		Expect(i.GetValue()).To(Equal(3.5))
		Expect(i.GetObject()).To(Equal(s))
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
})
