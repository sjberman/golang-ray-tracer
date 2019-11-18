package types

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTypes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

var _ = Describe("tuple tests", func() {
	It("creates vectors", func() {
		expVector := &Tuple{
			xAxis:     1,
			yAxis:     2,
			zAxis:     3,
			tupleType: vector,
		}
		Expect(CreateVector(1, 2, 3)).To(Equal(expVector))
	})

	It("creates points", func() {
		expPoint := &Tuple{
			xAxis:     1,
			yAxis:     2,
			zAxis:     3,
			tupleType: point,
		}
		Expect(CreatePoint(1, 2, 3)).To(Equal(expPoint))
	})

	It("returns whether a tuple is a vector or point", func() {
		p := CreatePoint(1, 2, 3)
		Expect(p.IsPoint()).To(BeTrue())
		Expect(p.IsVector()).To(BeFalse())

		v := CreateVector(1, 2, 3)
		Expect(v.IsVector()).To(BeTrue())
		Expect(v.IsPoint()).To(BeFalse())
	})

	It("adds tuples", func() {
		// vector + vector
		v1 := CreateVector(1, 2, 3)
		v2 := CreateVector(4, 5, 6)
		expVector := CreateVector(5, 7, 9)
		Expect(v1.Add(v2)).To(Equal(expVector))

		// point + vector
		p := CreatePoint(6, 7, -8)
		v := CreateVector(1, -2, 3)
		expPoint := CreatePoint(7, 5, -5)
		Expect(p.Add(v)).To(Equal(expPoint))

		// point + point
		p1 := CreatePoint(1, 2, 3)
		p2 := CreatePoint(4, 5, 6)
		_, err := p1.Add(p2)
		Expect(err).To(HaveOccurred())
	})

	It("subtracts tuples", func() {
		// vector - vector
		v1 := CreateVector(1, 5, 6)
		v2 := CreateVector(4, 1, 6)
		expVector := CreateVector(-3, 4, 0)
		Expect(v1.Subtract(v2)).To(Equal(expVector))

		// point - vector
		p := CreatePoint(6, 7, -8)
		v := CreateVector(5, -2, 3)
		expPoint := CreatePoint(1, 9, -11)
		Expect(p.Subtract(v)).To(Equal(expPoint))

		// point - point
		p1 := CreatePoint(1, 2, 8)
		p2 := CreatePoint(-2, 5, 6)
		expVector = CreateVector(3, -3, 2)
		Expect(p1.Subtract(p2)).To(Equal(expVector))

		// vector - point
		_, err := v.Subtract(p)
		Expect(err).To(HaveOccurred())
	})

	It("multiplies tuples", func() {
		val := 2.5
		p := CreatePoint(2, -3, 5)
		expPoint := CreatePoint(5, -7.5, 12.5)
		Expect(p.Multiply(val)).To(Equal(expPoint))

		val = .5
		v := CreateVector(2, 3, 8)
		expVector := CreateVector(1, 1.5, 4)
		Expect(v.Multiply(val)).To(Equal(expVector))
	})

	It("divides tuples", func() {
		val := 2.0
		p := CreatePoint(2, 3, 6)
		expPoint := CreatePoint(1, 1.5, 3)
		Expect(p.Divide(val)).To(Equal(expPoint))

		val = .5
		v := CreateVector(2, 3, -8)
		expVector := CreateVector(4, 6, -16)
		Expect(v.Divide(val)).To(Equal(expVector))
	})

	It("checks the equivalence of tuples", func() {
		v1 := CreateVector(1.001, 2, -3.345)
		v2 := CreateVector(1.001, 2, -3.345)
		Expect(v1.Equals(v2)).To(BeTrue())

		v1 = CreateVector(1.002, 2, -3.346)
		v2 = CreateVector(1.002, 2, -3.345)
		Expect(v1.Equals(v2)).To(BeFalse())

		p1 := CreatePoint(-1, 2.5, 3.0000001)
		p2 := CreatePoint(-1, 2.50, 3.0000001)
		Expect(p1.Equals(p2)).To(BeTrue())

		p1 = CreatePoint(-1, 2.5, 3.0000001)
		p2 = CreatePoint(1, 2.5, 8)
		Expect(p1.Equals(p2)).To(BeFalse())

		p1 = CreatePoint(-1, 2.5, 3.0000001)
		p2 = CreatePoint(-1, 2.500000000000001, 3.0000001)
		Expect(p1.Equals(p2)).To(BeFalse())

		Expect(p1.Equals(v1)).To(BeFalse())
	})

	It("negates tuples", func() {
		v := CreateVector(1, 2, -3)
		expVector := CreateVector(-1, -2, 3)
		Expect(v.Negate()).To(Equal(expVector))

		p := CreatePoint(-1, 2, 3)
		expPoint := CreatePoint(1, -2, -3)
		Expect(p.Negate()).To(Equal(expPoint))
	})

	It("gets the magnitude of vectors", func() {
		v := CreateVector(1, 0, 0)
		Expect(v.Magnitude()).To(Equal(float64(1)))

		v = CreateVector(0, 1, 0)
		Expect(v.Magnitude()).To(Equal(float64(1)))

		v = CreateVector(0, 0, 1)
		Expect(v.Magnitude()).To(Equal(float64(1)))

		v = CreateVector(2, 7, -4)
		Expect(v.Magnitude()).To(Equal(float64(8.306623862918075)))

		v = CreateVector(2, 7, 4)
		Expect(v.Magnitude()).To(Equal(float64(8.306623862918075)))
	})

	It("normalizes vectors", func() {
		v := CreateVector(4, 0, 0)
		expVector := CreateVector(1, 0, 0)
		Expect(v.Normalize()).To(Equal(expVector))

		v = CreateVector(1, 2, 3)
		expVector = CreateVector(0.2672612419124244, 0.5345224838248488, 0.8017837257372732)
		Expect(v.Normalize()).To(Equal(expVector))
	})

	It("gets the dot product of two vectors", func() {
		v1 := CreateVector(1, 2, 3)
		v2 := CreateVector(2, 3, 4)
		Expect(v1.DotProduct(v2)).To(Equal(float64(20)))
	})

	It("gets the cross product of two vectors", func() {
		v1 := CreateVector(1, 2, 3)
		v2 := CreateVector(2, 3, 4)
		expVector := CreateVector(-1, 2, -1)
		Expect(v1.CrossProduct(v2)).To(Equal(expVector))

		expVector = CreateVector(1, -2, 1)
		Expect(v2.CrossProduct(v1)).To(Equal(expVector))
	})
})
