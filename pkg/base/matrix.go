package base

import "math"

const four = 4

// Matrix is a matrix of floating point numbers.
type Matrix struct {
	size int
	data [][]float64
}

// Identity is the identity matrix.
var Identity = Matrix{
	size: four,
	data: [][]float64{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	},
}

// NewMatrix returns a new Matrix object.
func NewMatrix(data [][]float64) *Matrix {
	return &Matrix{
		size: len(data),
		data: data,
	}
}

// returns an empty data set for new matrices.
func newData(size int) [][]float64 {
	data := make([][]float64, size)
	for i := range data {
		data[i] = make([]float64, size)
	}

	return data
}

// Equals returns whether or not two matrices are equal to each other.
func (m *Matrix) Equals(m2 *Matrix) bool {
	if (m == nil) != (m2 == nil) {
		return false
	}
	if m.size != m2.size {
		return false
	}
	for i, row := range m.data {
		for j := range row {
			if !EqualFloats(m.data[i][j], m2.data[i][j]) {
				return false
			}
		}
	}

	return true
}

// Multiply multiplies two matrices by each other and returns the result.
// Only applies to 4x4 matrices.
func (m *Matrix) Multiply(m2 *Matrix) *Matrix {
	res := NewMatrix(newData(m.size))
	for i := 0; i < m.size; i++ {
		for j := 0; j < m.size; j++ {
			iTuple := listToTuple(m.data[i])
			jTuple := NewTuple(m2.data[0][j], m2.data[1][j], m2.data[2][j], m2.data[3][j])
			res.data[i][j] = iTuple.DotProduct(jTuple)
		}
	}

	return res
}

// MultiplyTuple multiples a matrix by a tuple and returns the resulting tuple.
// Only applies to 4x4 matrices.
func (m *Matrix) MultiplyTuple(t *Tuple) *Tuple {
	newVals := make([]float64, m.size)
	for i, row := range m.data {
		newVals[i] = listToTuple(row).DotProduct(t)
	}

	return listToTuple(newVals)
}

// Transpose turns a matrix's rows into columns.
func (m *Matrix) Transpose() *Matrix {
	res := NewMatrix(newData(m.size))
	for i := 0; i < m.size; i++ {
		for j := 0; j < m.size; j++ {
			res.data[i][j] = m.data[j][i]
		}
	}

	return res
}

// Inverse returns the inversion of a 4x4 matrix.
func (m *Matrix) Inverse() *Matrix {
	s0 := m.data[0][0]*m.data[1][1] - m.data[1][0]*m.data[0][1]
	s1 := m.data[0][0]*m.data[1][2] - m.data[1][0]*m.data[0][2]
	s2 := m.data[0][0]*m.data[1][3] - m.data[1][0]*m.data[0][3]
	s3 := m.data[0][1]*m.data[1][2] - m.data[1][1]*m.data[0][2]
	s4 := m.data[0][1]*m.data[1][3] - m.data[1][1]*m.data[0][3]
	s5 := m.data[0][2]*m.data[1][3] - m.data[1][2]*m.data[0][3]

	c5 := m.data[2][2]*m.data[3][3] - m.data[3][2]*m.data[2][3]
	c4 := m.data[2][1]*m.data[3][3] - m.data[3][1]*m.data[2][3]
	c3 := m.data[2][1]*m.data[3][2] - m.data[3][1]*m.data[2][2]
	c2 := m.data[2][0]*m.data[3][3] - m.data[3][0]*m.data[2][3]
	c1 := m.data[2][0]*m.data[3][2] - m.data[3][0]*m.data[2][2]
	c0 := m.data[2][0]*m.data[3][1] - m.data[3][0]*m.data[2][1]

	// Should check for 0 determinant
	invdet := 1.0 / (s0*c5 - s1*c4 + s2*c3 + s3*c2 - s4*c1 + s5*c0)

	data := newData(four)
	data[0][0] = (m.data[1][1]*c5 - m.data[1][2]*c4 + m.data[1][3]*c3) * invdet
	data[0][1] = (-m.data[0][1]*c5 + m.data[0][2]*c4 - m.data[0][3]*c3) * invdet
	data[0][2] = (m.data[3][1]*s5 - m.data[3][2]*s4 + m.data[3][3]*s3) * invdet
	data[0][3] = (-m.data[2][1]*s5 + m.data[2][2]*s4 - m.data[2][3]*s3) * invdet

	data[1][0] = (-m.data[1][0]*c5 + m.data[1][2]*c2 - m.data[1][3]*c1) * invdet
	data[1][1] = (m.data[0][0]*c5 - m.data[0][2]*c2 + m.data[0][3]*c1) * invdet
	data[1][2] = (-m.data[3][0]*s5 + m.data[3][2]*s2 - m.data[3][3]*s1) * invdet
	data[1][3] = (m.data[2][0]*s5 - m.data[2][2]*s2 + m.data[2][3]*s1) * invdet

	data[2][0] = (m.data[1][0]*c4 - m.data[1][1]*c2 + m.data[1][3]*c0) * invdet
	data[2][1] = (-m.data[0][0]*c4 + m.data[0][1]*c2 - m.data[0][3]*c0) * invdet
	data[2][2] = (m.data[3][0]*s4 - m.data[3][1]*s2 + m.data[3][3]*s0) * invdet
	data[2][3] = (-m.data[2][0]*s4 + m.data[2][1]*s2 - m.data[2][3]*s0) * invdet

	data[3][0] = (-m.data[1][0]*c3 + m.data[1][1]*c1 - m.data[1][2]*c0) * invdet
	data[3][1] = (m.data[0][0]*c3 - m.data[0][1]*c1 + m.data[0][2]*c0) * invdet
	data[3][2] = (-m.data[3][0]*s3 + m.data[3][1]*s1 - m.data[3][2]*s0) * invdet
	data[3][3] = (m.data[2][0]*s3 - m.data[2][1]*s1 + m.data[2][2]*s0) * invdet

	return NewMatrix(data)
}

// Translate returns a translation matrix.
func Translate(x, y, z float64) *Matrix {
	data := [][]float64{
		{1, 0, 0, x},
		{0, 1, 0, y},
		{0, 0, 1, z},
		{0, 0, 0, 1},
	}

	return NewMatrix(data)
}

// Scale returns a scaling matrix.
func Scale(x, y, z float64) *Matrix {
	data := [][]float64{
		{x, 0, 0, 0},
		{0, y, 0, 0},
		{0, 0, z, 0},
		{0, 0, 0, 1},
	}

	return NewMatrix(data)
}

// RotateX returns an x-axis rotation matrix.
func RotateX(radians float64) *Matrix {
	data := [][]float64{
		{1, 0, 0, 0},
		{0, math.Cos(radians), -math.Sin(radians), 0},
		{0, math.Sin(radians), math.Cos(radians), 0},
		{0, 0, 0, 1},
	}

	return NewMatrix(data)
}

// RotateY returns a y-axis rotation matrix.
func RotateY(radians float64) *Matrix {
	data := [][]float64{
		{math.Cos(radians), 0, math.Sin(radians), 0},
		{0, 1, 0, 0},
		{-math.Sin(radians), 0, math.Cos(radians), 0},
		{0, 0, 0, 1},
	}

	return NewMatrix(data)
}

// RotateZ returns a z-axis rotation matrix.
func RotateZ(radians float64) *Matrix {
	data := [][]float64{
		{math.Cos(radians), -math.Sin(radians), 0, 0},
		{math.Sin(radians), math.Cos(radians), 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}

	return NewMatrix(data)
}

// Shear returns a shearing (or skewing) matrix.
func Shear(xy, xz, yx, yz, zx, zy float64) *Matrix {
	// xy means "x moved in proportion to y" and so on for the rest
	data := [][]float64{
		{1, xy, xz, 0},
		{yx, 1, yz, 0},
		{zx, zy, 1, 0},
		{0, 0, 0, 1},
	}

	return NewMatrix(data)
}

// ViewTransform returns a transformation matrix to move the world.
// from is the camera, to is where we look, and up is
// the vector indicating which direction is up.
func ViewTransform(from, to, up *Tuple) *Matrix {
	forward := to.Subtract(from).Normalize()
	upNormal := up.Normalize()
	left := forward.CrossProduct(upNormal)
	trueUp := left.CrossProduct(forward)

	data := [][]float64{
		{left.xAxis, left.yAxis, left.zAxis, 0},
		{trueUp.xAxis, trueUp.yAxis, trueUp.zAxis, 0},
		{-forward.xAxis, -forward.yAxis, -forward.zAxis, 0},
		{0, 0, 0, 1},
	}
	orientation := NewMatrix(data)

	return orientation.Multiply(Translate(-from.xAxis, -from.yAxis, -from.zAxis))
}
