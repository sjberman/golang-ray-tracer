package base

import (
	"math"
	"testing"

	. "github.com/onsi/gomega"
)

func TestNewTuple(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	expTuple := &Tuple{
		xAxis: 1,
		yAxis: 2,
		zAxis: 3,
		w:     4,
	}
	g.Expect(NewTuple(1, 2, 3, 4)).To(Equal(expTuple))
}

func TestNewVector(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	expVector := &Tuple{
		xAxis: 1,
		yAxis: 2,
		zAxis: 3,
		w:     vector,
	}
	g.Expect(NewVector(1, 2, 3)).To(Equal(expVector))
}

func TestNewPoint(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	expPoint := &Tuple{
		xAxis: 1,
		yAxis: 2,
		zAxis: 3,
		w:     point,
	}
	g.Expect(NewPoint(1, 2, 3)).To(Equal(expPoint))
}

func TestIsVectorOrIsPoint(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	p := NewPoint(1, 2, 3)
	g.Expect(p.IsPoint()).To(BeTrue())
	g.Expect(p.IsVector()).To(BeFalse())

	v := NewVector(1, 2, 3)
	g.Expect(v.IsVector()).To(BeTrue())
	g.Expect(v.IsPoint()).To(BeFalse())
}

func TestGetAxis(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	tuple := NewTuple(1, 2, 3, 4)
	g.Expect(tuple.GetX()).To(Equal(1.0))
	g.Expect(tuple.GetY()).To(Equal(2.0))
	g.Expect(tuple.GetZ()).To(Equal(3.0))
}

func TestSetW(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	tuple := NewTuple(1, 2, 3, 4)
	g.Expect(tuple.w).To(Equal(4.0))
	tuple.SetW(5)
	g.Expect(tuple.w).To(Equal(5.0))
}

func TestTupleAdd(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// vector + vector
	v1 := NewVector(1, 2, 3)
	v2 := NewVector(4, 5, 6)
	expVector := NewVector(5, 7, 9)
	g.Expect(v1.Add(v2)).To(Equal(expVector))

	// point + vector
	p := NewPoint(6, 7, -8)
	v := NewVector(1, -2, 3)
	expPoint := NewPoint(7, 5, -5)
	g.Expect(p.Add(v)).To(Equal(expPoint))

	// point + point
	p1 := NewPoint(1, 2, 3)
	p2 := NewPoint(4, 5, 6)
	g.Expect(p1.Add(p2)).To(BeNil())
}

func TestTupleSubtract(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// vector - vector
	v1 := NewVector(1, 5, 6)
	v2 := NewVector(4, 1, 6)
	expVector := NewVector(-3, 4, 0)
	g.Expect(v1.Subtract(v2)).To(Equal(expVector))

	// point - vector
	p := NewPoint(6, 7, -8)
	v := NewVector(5, -2, 3)
	expPoint := NewPoint(1, 9, -11)
	g.Expect(p.Subtract(v)).To(Equal(expPoint))

	// point - point
	p1 := NewPoint(1, 2, 8)
	p2 := NewPoint(-2, 5, 6)
	expVector = NewVector(3, -3, 2)
	g.Expect(p1.Subtract(p2)).To(Equal(expVector))

	// vector - point
	g.Expect(v.Subtract(p)).To(BeNil())
}

func TestTupleMultiply(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	val := 2.5
	p := NewPoint(2, -3, 5)
	expPoint := NewPoint(5, -7.5, 12.5)
	g.Expect(p.Multiply(val)).To(Equal(expPoint))

	val = .5
	v := NewVector(2, 3, 8)
	expVector := NewVector(1, 1.5, 4)
	g.Expect(v.Multiply(val)).To(Equal(expVector))
}

func TestTupleDivide(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	val := 2.0
	p := NewPoint(2, 3, 6)
	expPoint := NewPoint(1, 1.5, 3)
	g.Expect(p.Divide(val)).To(Equal(expPoint))

	val = .5
	v := NewVector(2, 3, -8)
	expVector := NewVector(4, 6, -16)
	g.Expect(v.Divide(val)).To(Equal(expVector))
}

func TestTupleEquals(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	v1 := NewVector(1.001, 2, -3.345)
	v2 := NewVector(1.001, 2, -3.345)
	g.Expect(v1.Equals(v2)).To(BeTrue())

	v1 = NewVector(1.002, 2, -3.346)
	v2 = NewVector(1.002, 2, -3.345)
	g.Expect(v1.Equals(v2)).To(BeFalse())

	p1 := NewPoint(-1, 2.5, 3.0000001)
	p2 := NewPoint(-1, 2.50, 3.0000001)
	g.Expect(p1.Equals(p2)).To(BeTrue())

	p1 = NewPoint(-1, 2.5, 3.0000001)
	p2 = NewPoint(1, 2.5, 8)
	g.Expect(p1.Equals(p2)).To(BeFalse())

	p1 = NewPoint(-1, 2.5, 3.0000001)
	p2 = NewPoint(-1, 2.50001, 3.0000001)
	g.Expect(p1.Equals(p2)).To(BeFalse())

	g.Expect(p1.Equals(v1)).To(BeFalse())
}

func TestGreaterOrLessThan(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	t1 := NewPoint(0, 0, 0)
	t2 := NewPoint(1, 1, 1)
	g.Expect(t1.LessThan(t2)).To(BeTrue())
	g.Expect(t1.GreaterThan(t2)).To(BeFalse())

	t2 = NewPoint(-1, 1, 1)
	g.Expect(t1.LessThan(t2)).To(BeFalse())

	t2 = NewPoint(1, -1, 1)
	g.Expect(t1.LessThan(t2)).To(BeFalse())

	t2 = NewPoint(1, 1, -1)
	g.Expect(t1.LessThan(t2)).To(BeFalse())

	t2 = NewPoint(-1, -1, -1)
	g.Expect(t1.GreaterThan(t2)).To(BeTrue())

	t2 = NewPoint(-1, 1, -1)
	g.Expect(t1.GreaterThan(t2)).To(BeFalse())

	t2 = NewPoint(-1, -1, 1)
	g.Expect(t1.GreaterThan(t2)).To(BeFalse())

	t1 = NewPoint(0, 0, 0)
	t2 = NewPoint(0, 0, 0)
	g.Expect(t1.LessThan(t2)).To(BeFalse())
	g.Expect(t1.GreaterThan(t2)).To(BeFalse())
}

func TestNegate(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	v := NewVector(1, 2, -3)
	expVector := NewVector(-1, -2, 3)
	g.Expect(v.Negate()).To(Equal(expVector))

	p := NewPoint(-1, 2, 3)
	expPoint := NewPoint(1, -2, -3)
	g.Expect(p.Negate()).To(Equal(expPoint))
}

func TestMagnitude(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	v := NewVector(1, 0, 0)
	g.Expect(v.Magnitude()).To(Equal(1.0))

	v = NewVector(0, 1, 0)
	g.Expect(v.Magnitude()).To(Equal(1.0))

	v = NewVector(0, 0, 1)
	g.Expect(v.Magnitude()).To(Equal(1.0))

	v = NewVector(2, 7, -4)
	g.Expect(v.Magnitude()).To(Equal(8.306623862918075))

	v = NewVector(2, 7, 4)
	g.Expect(v.Magnitude()).To(Equal(8.306623862918075))
}

func TestNormalize(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	v := NewVector(4, 0, 0)
	expVector := NewVector(1, 0, 0)
	g.Expect(v.Normalize()).To(Equal(expVector))

	v = NewVector(1, 2, 3)
	expVector = NewVector(0.2672612419124244, 0.5345224838248488, 0.8017837257372732)
	g.Expect(v.Normalize()).To(Equal(expVector))
}

func TestDotProduct(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	v1 := NewVector(1, 2, 3)
	v2 := NewVector(2, 3, 4)
	g.Expect(v1.DotProduct(v2)).To(Equal(20.0))
}

func TestCrossProduct(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	v1 := NewVector(1, 2, 3)
	v2 := NewVector(2, 3, 4)
	expVector := NewVector(-1, 2, -1)
	g.Expect(v1.CrossProduct(v2)).To(Equal(expVector))

	expVector = NewVector(1, -2, 1)
	g.Expect(v2.CrossProduct(v1)).To(Equal(expVector))
}

func TestReflect(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// reflecting a vector approaching at 45 degrees
	v := NewVector(1, -1, 0)
	normal := NewVector(0, 1, 0)
	expVector := NewVector(1, 1, 0)
	g.Expect(v.Reflect(normal).Equals(expVector)).To(BeTrue())

	// reflecting a vector off a slanted surface
	v = NewVector(0, -1, 0)
	normal = NewVector(math.Sqrt(2)/2, math.Sqrt(2)/2, 0)
	expVector = NewVector(1, 0, 0)
	g.Expect(v.Reflect(normal).Equals(expVector)).To(BeTrue())
}
