package base

import (
	"math"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

var _ = Describe("tuple tests", func() {
	It("creates tuples", func() {
		expTuple := &Tuple{
			xAxis: 1,
			yAxis: 2,
			zAxis: 3,
			w:     4,
		}
		Expect(NewTuple(1, 2, 3, 4)).To(Equal(expTuple))
	})

	It("creates vectors", func() {
		expVector := &Tuple{
			xAxis: 1,
			yAxis: 2,
			zAxis: 3,
			w:     vector,
		}
		Expect(NewVector(1, 2, 3)).To(Equal(expVector))
	})

	It("creates points", func() {
		expPoint := &Tuple{
			xAxis: 1,
			yAxis: 2,
			zAxis: 3,
			w:     point,
		}
		Expect(NewPoint(1, 2, 3)).To(Equal(expPoint))
	})

	It("returns whether a tuple is a vector or point", func() {
		p := NewPoint(1, 2, 3)
		Expect(p.IsPoint()).To(BeTrue())
		Expect(p.IsVector()).To(BeFalse())

		v := NewVector(1, 2, 3)
		Expect(v.IsVector()).To(BeTrue())
		Expect(v.IsPoint()).To(BeFalse())
	})

	It("returns each axis", func() {
		t := NewTuple(1, 2, 3, 4)
		Expect(t.GetX()).To(Equal(1.0))
		Expect(t.GetY()).To(Equal(2.0))
		Expect(t.GetZ()).To(Equal(3.0))
	})

	It("sets the w field", func() {
		t := NewTuple(1, 2, 3, 4)
		Expect(t.w).To(Equal(4.0))
		t.SetW(5)
		Expect(t.w).To(Equal(5.0))
	})

	It("adds tuples", func() {
		// vector + vector
		v1 := NewVector(1, 2, 3)
		v2 := NewVector(4, 5, 6)
		expVector := NewVector(5, 7, 9)
		Expect(v1.Add(v2)).To(Equal(expVector))

		// point + vector
		p := NewPoint(6, 7, -8)
		v := NewVector(1, -2, 3)
		expPoint := NewPoint(7, 5, -5)
		Expect(p.Add(v)).To(Equal(expPoint))

		// point + point
		p1 := NewPoint(1, 2, 3)
		p2 := NewPoint(4, 5, 6)
		Expect(p1.Add(p2)).To(BeNil())
	})

	It("subtracts tuples", func() {
		// vector - vector
		v1 := NewVector(1, 5, 6)
		v2 := NewVector(4, 1, 6)
		expVector := NewVector(-3, 4, 0)
		Expect(v1.Subtract(v2)).To(Equal(expVector))

		// point - vector
		p := NewPoint(6, 7, -8)
		v := NewVector(5, -2, 3)
		expPoint := NewPoint(1, 9, -11)
		Expect(p.Subtract(v)).To(Equal(expPoint))

		// point - point
		p1 := NewPoint(1, 2, 8)
		p2 := NewPoint(-2, 5, 6)
		expVector = NewVector(3, -3, 2)
		Expect(p1.Subtract(p2)).To(Equal(expVector))

		// vector - point
		Expect(v.Subtract(p)).To(BeNil())
	})

	It("multiplies tuples", func() {
		val := 2.5
		p := NewPoint(2, -3, 5)
		expPoint := NewPoint(5, -7.5, 12.5)
		Expect(p.Multiply(val)).To(Equal(expPoint))

		val = .5
		v := NewVector(2, 3, 8)
		expVector := NewVector(1, 1.5, 4)
		Expect(v.Multiply(val)).To(Equal(expVector))
	})

	It("divides tuples", func() {
		val := 2.0
		p := NewPoint(2, 3, 6)
		expPoint := NewPoint(1, 1.5, 3)
		Expect(p.Divide(val)).To(Equal(expPoint))

		val = .5
		v := NewVector(2, 3, -8)
		expVector := NewVector(4, 6, -16)
		Expect(v.Divide(val)).To(Equal(expVector))
	})

	It("checks the equivalence of tuples", func() {
		v1 := NewVector(1.001, 2, -3.345)
		v2 := NewVector(1.001, 2, -3.345)
		Expect(v1.Equals(v2)).To(BeTrue())

		v1 = NewVector(1.002, 2, -3.346)
		v2 = NewVector(1.002, 2, -3.345)
		Expect(v1.Equals(v2)).To(BeFalse())

		p1 := NewPoint(-1, 2.5, 3.0000001)
		p2 := NewPoint(-1, 2.50, 3.0000001)
		Expect(p1.Equals(p2)).To(BeTrue())

		p1 = NewPoint(-1, 2.5, 3.0000001)
		p2 = NewPoint(1, 2.5, 8)
		Expect(p1.Equals(p2)).To(BeFalse())

		p1 = NewPoint(-1, 2.5, 3.0000001)
		p2 = NewPoint(-1, 2.50001, 3.0000001)
		Expect(p1.Equals(p2)).To(BeFalse())

		Expect(p1.Equals(v1)).To(BeFalse())
	})

	It("checks whether tuples are greater or less than each other", func() {
		t1 := NewPoint(0, 0, 0)
		t2 := NewPoint(1, 1, 1)
		Expect(t1.LessThan(t2)).To(BeTrue())
		Expect(t1.GreaterThan(t2)).To(BeFalse())

		t2 = NewPoint(-1, 1, 1)
		Expect(t1.LessThan(t2)).To(BeFalse())

		t2 = NewPoint(1, -1, 1)
		Expect(t1.LessThan(t2)).To(BeFalse())

		t2 = NewPoint(1, 1, -1)
		Expect(t1.LessThan(t2)).To(BeFalse())

		t2 = NewPoint(-1, -1, -1)
		Expect(t1.GreaterThan(t2)).To(BeTrue())

		t2 = NewPoint(-1, 1, -1)
		Expect(t1.GreaterThan(t2)).To(BeFalse())

		t2 = NewPoint(-1, -1, 1)
		Expect(t1.GreaterThan(t2)).To(BeFalse())

		t1 = NewPoint(0, 0, 0)
		t2 = NewPoint(0, 0, 0)
		Expect(t1.LessThan(t2)).To(BeFalse())
		Expect(t1.GreaterThan(t2)).To(BeFalse())
	})

	It("negates tuples", func() {
		v := NewVector(1, 2, -3)
		expVector := NewVector(-1, -2, 3)
		Expect(v.Negate()).To(Equal(expVector))

		p := NewPoint(-1, 2, 3)
		expPoint := NewPoint(1, -2, -3)
		Expect(p.Negate()).To(Equal(expPoint))
	})

	It("gets the magnitude of vectors", func() {
		v := NewVector(1, 0, 0)
		Expect(v.Magnitude()).To(Equal(1.0))

		v = NewVector(0, 1, 0)
		Expect(v.Magnitude()).To(Equal(1.0))

		v = NewVector(0, 0, 1)
		Expect(v.Magnitude()).To(Equal(1.0))

		v = NewVector(2, 7, -4)
		Expect(v.Magnitude()).To(Equal(8.306623862918075))

		v = NewVector(2, 7, 4)
		Expect(v.Magnitude()).To(Equal(8.306623862918075))
	})

	It("normalizes vectors", func() {
		v := NewVector(4, 0, 0)
		expVector := NewVector(1, 0, 0)
		Expect(v.Normalize()).To(Equal(expVector))

		v = NewVector(1, 2, 3)
		expVector = NewVector(0.2672612419124244, 0.5345224838248488, 0.8017837257372732)
		Expect(v.Normalize()).To(Equal(expVector))
	})

	It("gets the dot product of two vectors", func() {
		v1 := NewVector(1, 2, 3)
		v2 := NewVector(2, 3, 4)
		Expect(v1.DotProduct(v2)).To(Equal(20.0))
	})

	It("gets the cross product of two vectors", func() {
		v1 := NewVector(1, 2, 3)
		v2 := NewVector(2, 3, 4)
		expVector := NewVector(-1, 2, -1)
		Expect(v1.CrossProduct(v2)).To(Equal(expVector))

		expVector = NewVector(1, -2, 1)
		Expect(v2.CrossProduct(v1)).To(Equal(expVector))
	})

	It("reflects vectors around a normal", func() {
		// reflecting a vector approaching at 45 degrees
		v := NewVector(1, -1, 0)
		normal := NewVector(0, 1, 0)
		expVector := NewVector(1, 1, 0)
		Expect(v.Reflect(normal).Equals(expVector)).To(BeTrue())

		// reflecting a vector off a slanted surface
		v = NewVector(0, -1, 0)
		normal = NewVector(math.Sqrt(2)/2, math.Sqrt(2)/2, 0)
		expVector = NewVector(1, 0, 0)
		Expect(v.Reflect(normal).Equals(expVector)).To(BeTrue())
	})
})
