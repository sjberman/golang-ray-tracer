package base

import (
	"math"
	"testing"

	. "github.com/onsi/gomega"
)

func TestNewMatrix(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	data := [][]float64{
		{1, 2, 3, 4},
		{5.5, 6.5, 7.5, 8.5},
		{9, 10, 11, 12},
		{13.5, 14.5, 15.5, 16.5},
	}
	m := NewMatrix(data)
	g.Expect(m.data[0][0]).To(Equal(1.0))
	g.Expect(m.data[0][3]).To(Equal(4.0))
	g.Expect(m.data[1][0]).To(Equal(5.5))
	g.Expect(m.data[1][2]).To(Equal(7.5))
	g.Expect(m.data[3][0]).To(Equal(13.5))
	g.Expect(m.data[3][2]).To(Equal(15.5))
}

func TestMatrixEquals(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	data := [][]float64{
		{1, 2, 3, 4},
		{5.5, 6.5, 7.5, 8.5},
		{9, 10, 11, 12},
		{13.5, 14.5, 15.5, 16.5},
	}
	m1 := NewMatrix(data)
	m2 := NewMatrix(data)
	g.Expect(m1.Equals(m2)).To(BeTrue())

	data = [][]float64{
		{1, 2, 3, 4},
	}
	m2 = NewMatrix(data)
	g.Expect(m1.Equals(m2)).To(BeFalse())

	data = [][]float64{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
		{13, 14, 15, 16},
	}
	m2 = NewMatrix(data)
	g.Expect(m1.Equals(m2)).To(BeFalse())
	g.Expect(m1.Equals(nil)).To(BeFalse())
}

func TestMatrixMultiply(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	data := [][]float64{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 8, 7, 6},
		{5, 4, 3, 2},
	}
	m1 := NewMatrix(data)

	data = [][]float64{
		{-2, 1, 2, 3},
		{3, 2, 1, -1},
		{4, 3, 6, 5},
		{1, 2, 7, 8},
	}
	m2 := NewMatrix(data)

	data = [][]float64{
		{20, 22, 50, 48},
		{44, 54, 114, 108},
		{40, 58, 110, 102},
		{16, 26, 46, 42},
	}
	exp := NewMatrix(data)
	g.Expect(m1.Multiply(m2)).To(Equal(exp))
}

func TestMatrixMultiplyTuple(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	data := [][]float64{
		{1, 2, 3, 4},
		{2, 4, 4, 2},
		{8, 6, 4, 1},
		{0, 0, 0, 1},
	}
	m := NewMatrix(data)
	point := NewPoint(1, 2, 3)
	expTuple := NewPoint(18, 24, 33)
	g.Expect(m.MultiplyTuple(point)).To(Equal(expTuple))
}

func TestMatrixMultiplyIdentity(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	data := [][]float64{
		{0, 1, 2, 4},
		{1, 2, 4, 8},
		{2, 4, 8, 16},
		{4, 8, 16, 32},
	}
	m := NewMatrix(data)
	g.Expect(m.Multiply(&Identity)).To(Equal(m))

	tuple := NewTuple(1, 2, 3, 4)
	g.Expect(Identity.MultiplyTuple(tuple)).To(Equal(tuple))
}

func TestTranspose(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	data := [][]float64{
		{0, 9, 3, 0},
		{9, 8, 0, 8},
		{1, 8, 5, 3},
		{0, 0, 5, 8},
	}
	m := NewMatrix(data)

	data = [][]float64{
		{0, 9, 1, 0},
		{9, 8, 8, 0},
		{3, 0, 5, 5},
		{0, 8, 3, 8},
	}
	expMatrix := NewMatrix(data)
	g.Expect(m.Transpose()).To(Equal(expMatrix))

	g.Expect(Identity.Transpose()).To(Equal(&Identity))
}

func TestInverse(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	data := [][]float64{
		{-5, 2, 6, -8},
		{1, -5, 1, 8},
		{7, 7, -6, -7},
		{1, -3, 7, 4},
	}
	m := NewMatrix(data)
	invertM := m.Inverse()

	data = [][]float64{
		{0.21804511278195488, 0.45112781954887216, 0.24060150375939848, -0.045112781954887216},
		{-0.8082706766917293, -1.456766917293233, -0.44360902255639095, 0.5206766917293233},
		{-0.07894736842105263, -0.22368421052631576, -0.05263157894736842, 0.19736842105263158},
		{-0.5225563909774436, -0.8139097744360901, -0.3007518796992481, 0.306390977443609},
	}
	expMatrix := NewMatrix(data)
	g.Expect(invertM).To(Equal(expMatrix))

	// A * B = C, C * inverse(B) = A
	data2 := [][]float64{
		{3, -9, 7, 3},
		{3, -8, 2, -9},
		{-4, 4, 4, 1},
		{-6, 5, -1, 1},
	}
	m1 := NewMatrix(data2)

	data2 = [][]float64{
		{8, 2, 2, 2},
		{3, -1, 7, 0},
		{7, 0, 5, 4},
		{6, -2, 0, 5},
	}
	m2 := NewMatrix(data2)

	m3 := m1.Multiply(m2)
	inverse := m2.Inverse()

	res := m3.Multiply(inverse)
	g.Expect(res.Equals(m1)).To(BeTrue())
}

func TestTranslate(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	translate := Translate(5, -3, 2)
	p := NewPoint(-3, 4, 5)
	expPoint := NewPoint(2, 1, 7)
	g.Expect(translate.MultiplyTuple(p)).To(Equal(expPoint))

	inverse := translate.Inverse()
	expPoint = NewPoint(-8, 7, 3)
	g.Expect(inverse.MultiplyTuple(p)).To(Equal(expPoint))

	v := NewVector(-3, 4, 5)
	g.Expect(translate.MultiplyTuple(v)).To(Equal(v))
}

func TestScale(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	s := Scale(2, 3, 4)
	p := NewPoint(-4, 6, 8)
	expPoint := NewPoint(-8, 18, 32)
	g.Expect(s.MultiplyTuple(p)).To(Equal(expPoint))

	v := NewVector(-4, 6, 8)
	expVector := NewVector(-8, 18, 32)
	g.Expect(s.MultiplyTuple(v)).To(Equal(expVector))

	inverse := s.Inverse()
	expVector = NewVector(-2, 2, 2)
	g.Expect(inverse.MultiplyTuple(v)).To(Equal(expVector))

	// reflection
	s = Scale(-1, 1, 1)
	p = NewPoint(2, 3, 4)
	expPoint = NewPoint(-2, 3, 4)
	g.Expect(s.MultiplyTuple(p)).To(Equal(expPoint))
}

func TestRotation(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// X axis
	p := NewPoint(0, 1, 0)
	rx := RotateX(math.Pi / 4)
	expPoint := NewPoint(0, math.Sqrt(2)/2, math.Sqrt(2)/2)
	res := rx.MultiplyTuple(p)
	g.Expect(res.Equals(expPoint)).To(BeTrue())

	rx = RotateX(math.Pi / 2)
	expPoint = NewPoint(0, 0, 1)
	res = rx.MultiplyTuple(p)
	g.Expect(res.Equals(expPoint)).To(BeTrue())

	// Y axis
	p = NewPoint(0, 0, 1)
	ry := RotateY(math.Pi / 4)
	expPoint = NewPoint(math.Sqrt(2)/2, 0, math.Sqrt(2)/2)
	res = ry.MultiplyTuple(p)
	g.Expect(res.Equals(expPoint)).To(BeTrue())

	ry = RotateY(math.Pi / 2)
	expPoint = NewPoint(1, 0, 0)
	res = ry.MultiplyTuple(p)
	g.Expect(res.Equals(expPoint)).To(BeTrue())

	// Z axis
	p = NewPoint(0, 1, 0)
	rz := RotateZ(math.Pi / 4)
	expPoint = NewPoint(-math.Sqrt(2)/2, math.Sqrt(2)/2, 0)
	res = rz.MultiplyTuple(p)
	g.Expect(res.Equals(expPoint)).To(BeTrue())

	rz = RotateZ(math.Pi / 2)
	expPoint = NewPoint(-1, 0, 0)
	res = rz.MultiplyTuple(p)
	g.Expect(res.Equals(expPoint)).To(BeTrue())
}

func TestShear(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// moves x in proportion to y
	s := Shear(1, 0, 0, 0, 0, 0)
	p := NewPoint(2, 3, 4)
	expPoint := NewPoint(5, 3, 4)
	g.Expect(s.MultiplyTuple(p)).To(Equal(expPoint))

	// moves x in proportion to z
	s = Shear(0, 1, 0, 0, 0, 0)
	p = NewPoint(2, 3, 4)
	expPoint = NewPoint(6, 3, 4)
	g.Expect(s.MultiplyTuple(p)).To(Equal(expPoint))

	// moves y in proportion to x
	s = Shear(0, 0, 1, 0, 0, 0)
	p = NewPoint(2, 3, 4)
	expPoint = NewPoint(2, 5, 4)
	g.Expect(s.MultiplyTuple(p)).To(Equal(expPoint))

	// moves y in proportion to z
	s = Shear(0, 0, 0, 1, 0, 0)
	p = NewPoint(2, 3, 4)
	expPoint = NewPoint(2, 7, 4)
	g.Expect(s.MultiplyTuple(p)).To(Equal(expPoint))

	// moves z in proportion to x
	s = Shear(0, 0, 0, 0, 1, 0)
	p = NewPoint(2, 3, 4)
	expPoint = NewPoint(2, 3, 6)
	g.Expect(s.MultiplyTuple(p)).To(Equal(expPoint))

	// moves z in proportion to y
	s = Shear(0, 0, 0, 0, 0, 1)
	p = NewPoint(2, 3, 4)
	expPoint = NewPoint(2, 3, 7)
	g.Expect(s.MultiplyTuple(p)).To(Equal(expPoint))
}

func TestChainTransformations(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// individual transformations applied in sequence
	p := NewPoint(1, 0, 1)
	rx := RotateX(math.Pi / 2)
	s := Scale(5, 5, 5)
	translate := Translate(10, 5, 7)

	p2 := rx.MultiplyTuple(p)
	expPoint := NewPoint(1, -1, 0)
	g.Expect(p2.Equals(expPoint)).To(BeTrue())

	p3 := s.MultiplyTuple(p2)
	expPoint = NewPoint(5, -5, 0)
	g.Expect(p3.Equals(expPoint)).To(BeTrue())

	p4 := translate.MultiplyTuple(p3)
	expPoint = NewPoint(15, 0, 7)
	g.Expect(p4.Equals(expPoint)).To(BeTrue())

	// chained transformations applied in reverse order
	m := translate.Multiply(s.Multiply(rx))
	res := m.MultiplyTuple(p)
	g.Expect(res.Equals(expPoint)).To(BeTrue())
}

func TestViewTransform(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// default orientation
	from := Origin
	to := NewPoint(0, 0, -1)
	up := NewVector(0, 1, 0)
	m := ViewTransform(from, to, up)
	g.Expect(m).To(Equal(&Identity))

	// positive z direction (mirror of default)
	to = NewPoint(0, 0, 1)
	m = ViewTransform(from, to, up)
	g.Expect(m).To(Equal(Scale(-1, 1, -1)))

	// moves the world
	from = NewPoint(0, 0, 8)
	to = Origin
	up = NewVector(0, 1, 0)
	m = ViewTransform(from, to, up)
	g.Expect(m).To(Equal(Translate(0, 0, -8)))

	// arbitrary view
	from = NewPoint(1, 3, 2)
	to = NewPoint(4, -2, 8)
	up = NewVector(1, 1, 0)
	m = ViewTransform(from, to, up)
	data := [][]float64{
		{-0.5070925528371099, 0.5070925528371099, 0.6761234037828132, -2.366431913239846},
		{0.7677159338596801, 0.6060915267313263, 0.12121830534626524, -2.8284271247461894},
		{-0.35856858280031806, 0.5976143046671968, -0.7171371656006361, 0},
		{0, 0, 0, 1},
	}

	for i := range data {
		for j := range data[i] {
			g.Expect(m.data[i][j]).To(BeNumerically("~", data[i][j]))
		}
	}
}
