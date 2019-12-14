package base

import (
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("matrix tests", func() {
	It("creates matrices", func() {
		data := [][]float64{
			{1, 2, 3, 4},
			{5.5, 6.5, 7.5, 8.5},
			{9, 10, 11, 12},
			{13.5, 14.5, 15.5, 16.5},
		}
		m := NewMatrix(data)
		Expect(m.data[0][0]).To(Equal(1.0))
		Expect(m.data[0][3]).To(Equal(4.0))
		Expect(m.data[1][0]).To(Equal(5.5))
		Expect(m.data[1][2]).To(Equal(7.5))
		Expect(m.data[3][0]).To(Equal(13.5))
		Expect(m.data[3][2]).To(Equal(15.5))
	})

	It("checks the equivalence of matrices", func() {
		data := [][]float64{
			{1, 2, 3, 4},
			{5.5, 6.5, 7.5, 8.5},
			{9, 10, 11, 12},
			{13.5, 14.5, 15.5, 16.5},
		}
		m1 := NewMatrix(data)
		m2 := NewMatrix(data)
		Expect(m1.Equals(m2)).To(BeTrue())

		data = [][]float64{
			{1, 2, 3, 4},
		}
		m2 = NewMatrix(data)
		Expect(m1.Equals(m2)).To(BeFalse())

		data = [][]float64{
			{1, 2, 3, 4},
			{5, 6, 7, 8},
			{9, 10, 11, 12},
			{13, 14, 15, 16},
		}
		m2 = NewMatrix(data)
		Expect(m1.Equals(m2)).To(BeFalse())
		Expect(m1.Equals(nil)).To(BeFalse())
	})

	It("multiplies two matrices together", func() {
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
		Expect(m1.Multiply(m2)).To(Equal(exp))
	})

	It("multiplies a matrix and a tuple together", func() {
		data := [][]float64{
			{1, 2, 3, 4},
			{2, 4, 4, 2},
			{8, 6, 4, 1},
			{0, 0, 0, 1},
		}
		m := NewMatrix(data)
		t := NewPoint(1, 2, 3)
		expTuple := NewPoint(18, 24, 33)
		Expect(m.MultiplyTuple(t)).To(Equal(expTuple))
	})

	It("multiples by the identity matrix to get the same result", func() {
		data := [][]float64{
			{0, 1, 2, 4},
			{1, 2, 4, 8},
			{2, 4, 8, 16},
			{4, 8, 16, 32},
		}
		m := NewMatrix(data)
		Expect(m.Multiply(&Identity)).To(Equal(m))

		t := NewTuple(1, 2, 3, 4)
		Expect(Identity.MultiplyTuple(t)).To(Equal(t))
	})

	It("transposes matrices", func() {
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
		Expect(m.Transpose()).To(Equal(expMatrix))

		Expect(Identity.Transpose()).To(Equal(&Identity))
	})

	It("computes the determinant of a 2x2 matrix", func() {
		data := [][]float64{
			{1, 5},
			{-3, 2},
		}
		m := NewMatrix(data)
		Expect(m.determinant()).To(Equal(17.0))
	})

	It("creates submatrices from larger matrices", func() {
		// reducing 3x3 to 2x2
		data := [][]float64{
			{1, 5, 0},
			{-3, 2, 7},
			{0, 6, -3},
		}
		m := NewMatrix(data)

		data = [][]float64{
			{-3, 2},
			{0, 6},
		}
		expMatrix := NewMatrix(data)
		Expect(m.subMatrix(0, 2)).To(Equal(expMatrix))

		// reducing 4x4 to 3x3
		data = [][]float64{
			{-6, 1, 1, 6},
			{-8, 5, 8, 6},
			{-1, 0, 8, 2},
			{-7, 1, -1, 1},
		}
		m = NewMatrix(data)

		data = [][]float64{
			{-6, 1, 6},
			{-8, 8, 6},
			{-7, -1, 1},
		}
		expMatrix = NewMatrix(data)
		Expect(m.subMatrix(2, 1)).To(Equal(expMatrix))
	})

	It("computes the minor of a matrix", func() {
		data := [][]float64{
			{3, 5, 0},
			{2, -1, -7},
			{6, -1, 5},
		}
		m := NewMatrix(data)
		Expect(m.minor(1, 0)).To(Equal(25.0))
	})

	It("computes the cofactor of a matrix", func() {
		data := [][]float64{
			{3, 5, 0},
			{2, -1, -7},
			{6, -1, 5},
		}
		m := NewMatrix(data)
		Expect(m.minor(0, 0)).To(Equal(-12.0))
		Expect(m.cofactor(0, 0)).To(Equal(-12.0))
		Expect(m.minor(1, 0)).To(Equal(25.0))
		Expect(m.cofactor(1, 0)).To(Equal(-25.0))
	})

	It("computes the determinant of a 3x3 matrix", func() {
		data := [][]float64{
			{1, 2, 6},
			{-5, 8, -4},
			{2, 6, 4},
		}
		m := NewMatrix(data)
		Expect(m.determinant()).To(Equal(-196.0))
	})

	It("computes the determinant of a 4x4 matrix", func() {
		data := [][]float64{
			{-2, -8, 3, 5},
			{-3, 1, 7, 3},
			{1, 2, -9, 6},
			{-6, 7, 7, -9},
		}
		m := NewMatrix(data)
		Expect(m.determinant()).To(Equal(-4071.0))
	})

	It("inverts a matrix", func() {
		// determinant 0, not invertible
		data := [][]float64{
			{-4, 2, -2, 3},
			{9, 6, 2, 6},
			{0, -5, 1, -5},
			{0, 0, 0, 0},
		}
		m := NewMatrix(data)
		_, err := m.Inverse()
		Expect(err).To(HaveOccurred())

		data = [][]float64{
			{-5, 2, 6, -8},
			{1, -5, 1, 8},
			{7, 7, -6, -7},
			{1, -3, 7, 4},
		}
		m = NewMatrix(data)
		invertM, err := m.Inverse()
		Expect(err).ToNot(HaveOccurred())

		data = [][]float64{
			{0.21804511278195488, 0.45112781954887216, 0.24060150375939848, -0.045112781954887216},
			{-0.8082706766917294, -1.4567669172932332, -0.44360902255639095, 0.5206766917293233},
			{-0.07894736842105263, -0.2236842105263158, -0.05263157894736842, 0.19736842105263158},
			{-0.5225563909774437, -0.8139097744360902, -0.3007518796992481, 0.30639097744360905},
		}
		expMatrix := NewMatrix(data)
		Expect(invertM).To(Equal(expMatrix))

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
		inverse, err := m2.Inverse()
		Expect(err).ToNot(HaveOccurred())

		res := m3.Multiply(inverse)
		Expect(res.Equals(m1)).To(BeTrue())
	})

	It("creates a translation matrix", func() {
		t := Translate(5, -3, 2)
		p := NewPoint(-3, 4, 5)
		expPoint := NewPoint(2, 1, 7)
		Expect(t.MultiplyTuple(p)).To(Equal(expPoint))

		inverse, err := t.Inverse()
		Expect(err).ToNot(HaveOccurred())
		expPoint = NewPoint(-8, 7, 3)
		Expect(inverse.MultiplyTuple(p)).To(Equal(expPoint))

		v := NewVector(-3, 4, 5)
		Expect(t.MultiplyTuple(v)).To(Equal(v))
	})

	It("creates a scaling matrix", func() {
		s := Scale(2, 3, 4)
		p := NewPoint(-4, 6, 8)
		expPoint := NewPoint(-8, 18, 32)
		Expect(s.MultiplyTuple(p)).To(Equal(expPoint))

		v := NewVector(-4, 6, 8)
		expVector := NewVector(-8, 18, 32)
		Expect(s.MultiplyTuple(v)).To(Equal(expVector))

		inverse, err := s.Inverse()
		Expect(err).ToNot(HaveOccurred())
		expVector = NewVector(-2, 2, 2)
		Expect(inverse.MultiplyTuple(v)).To(Equal(expVector))

		// reflection
		s = Scale(-1, 1, 1)
		p = NewPoint(2, 3, 4)
		expPoint = NewPoint(-2, 3, 4)
		Expect(s.MultiplyTuple(p)).To(Equal(expPoint))
	})

	It("creates rotation matrices", func() {
		// X axis
		p := NewPoint(0, 1, 0)
		rx := RotateX(math.Pi / 4)
		expPoint := NewPoint(0, math.Sqrt(2)/2, math.Sqrt(2)/2)
		res := rx.MultiplyTuple(p)
		Expect(res.Equals(expPoint)).To(BeTrue())

		rx = RotateX(math.Pi / 2)
		expPoint = NewPoint(0, 0, 1)
		res = rx.MultiplyTuple(p)
		Expect(res.Equals(expPoint)).To(BeTrue())

		// Y axis
		p = NewPoint(0, 0, 1)
		ry := RotateY(math.Pi / 4)
		expPoint = NewPoint(math.Sqrt(2)/2, 0, math.Sqrt(2)/2)
		res = ry.MultiplyTuple(p)
		Expect(res.Equals(expPoint)).To(BeTrue())

		ry = RotateY(math.Pi / 2)
		expPoint = NewPoint(1, 0, 0)
		res = ry.MultiplyTuple(p)
		Expect(res.Equals(expPoint)).To(BeTrue())

		// Z axis
		p = NewPoint(0, 1, 0)
		rz := RotateZ(math.Pi / 4)
		expPoint = NewPoint(-math.Sqrt(2)/2, math.Sqrt(2)/2, 0)
		res = rz.MultiplyTuple(p)
		Expect(res.Equals(expPoint)).To(BeTrue())

		rz = RotateZ(math.Pi / 2)
		expPoint = NewPoint(-1, 0, 0)
		res = rz.MultiplyTuple(p)
		Expect(res.Equals(expPoint)).To(BeTrue())
	})

	It("creates a shearing matrix", func() {
		// moves x in proportion to y
		s := Shear(1, 0, 0, 0, 0, 0)
		p := NewPoint(2, 3, 4)
		expPoint := NewPoint(5, 3, 4)
		Expect(s.MultiplyTuple(p)).To(Equal(expPoint))

		// moves x in proportion to z
		s = Shear(0, 1, 0, 0, 0, 0)
		p = NewPoint(2, 3, 4)
		expPoint = NewPoint(6, 3, 4)
		Expect(s.MultiplyTuple(p)).To(Equal(expPoint))

		// moves y in proportion to x
		s = Shear(0, 0, 1, 0, 0, 0)
		p = NewPoint(2, 3, 4)
		expPoint = NewPoint(2, 5, 4)
		Expect(s.MultiplyTuple(p)).To(Equal(expPoint))

		// moves y in proportion to z
		s = Shear(0, 0, 0, 1, 0, 0)
		p = NewPoint(2, 3, 4)
		expPoint = NewPoint(2, 7, 4)
		Expect(s.MultiplyTuple(p)).To(Equal(expPoint))

		// moves z in proportion to x
		s = Shear(0, 0, 0, 0, 1, 0)
		p = NewPoint(2, 3, 4)
		expPoint = NewPoint(2, 3, 6)
		Expect(s.MultiplyTuple(p)).To(Equal(expPoint))

		// moves z in proportion to y
		s = Shear(0, 0, 0, 0, 0, 1)
		p = NewPoint(2, 3, 4)
		expPoint = NewPoint(2, 3, 7)
		Expect(s.MultiplyTuple(p)).To(Equal(expPoint))
	})

	It("chains transformations", func() {
		// individual transformations applied in sequence
		p := NewPoint(1, 0, 1)
		rx := RotateX(math.Pi / 2)
		s := Scale(5, 5, 5)
		t := Translate(10, 5, 7)

		p2 := rx.MultiplyTuple(p)
		expPoint := NewPoint(1, -1, 0)
		Expect(p2.Equals(expPoint)).To(BeTrue())

		p3 := s.MultiplyTuple(p2)
		expPoint = NewPoint(5, -5, 0)
		Expect(p3.Equals(expPoint)).To(BeTrue())

		p4 := t.MultiplyTuple(p3)
		expPoint = NewPoint(15, 0, 7)
		Expect(p4.Equals(expPoint)).To(BeTrue())

		// chained transformations applied in reverse order
		m := t.Multiply(s.Multiply(rx))
		res := m.MultiplyTuple(p)
		Expect(res.Equals(expPoint)).To(BeTrue())
	})

	It("creates a view transform matrix", func() {
		// default orientation
		from := Origin
		to := NewPoint(0, 0, -1)
		up := NewVector(0, 1, 0)
		m := ViewTransform(from, to, up)
		Expect(m).To(Equal(&Identity))

		// positive z direction (mirror of default)
		to = NewPoint(0, 0, 1)
		m = ViewTransform(from, to, up)
		Expect(m).To(Equal(Scale(-1, 1, -1)))

		// moves the world
		from = NewPoint(0, 0, 8)
		to = Origin
		up = NewVector(0, 1, 0)
		m = ViewTransform(from, to, up)
		Expect(m).To(Equal(Translate(0, 0, -8)))

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
		expMatrix := NewMatrix(data)
		Expect(m).To(Equal(expMatrix))
	})
})
