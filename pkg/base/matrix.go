package base

import (
	"errors"
)

// Matrix is a matrix of floating point numbers
type Matrix struct {
	size int
	data [][]float64
}

var Identity = &Matrix{
	size: 4,
	data: [][]float64{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	},
}

// NewMatrix returns a new Matrix object
func NewMatrix(data [][]float64) *Matrix {
	return &Matrix{
		size: len(data),
		data: data,
	}
}

// returns an empty data set for new matrices
func newData(size int) [][]float64 {
	data := make([][]float64, size)
	for i := range data {
		data[i] = make([]float64, size)
	}
	return data
}

// Equals returns whether or not two matrices are equal to each other
func (m *Matrix) Equals(m2 *Matrix) bool {
	if (m == nil) != (m2 == nil) {
		return false
	}
	if m.size != m2.size {
		return false
	}
	for i, row := range m.data {
		for j := range row {
			if !equalFloats(m.data[i][j], m2.data[i][j]) {
				return false
			}
		}
	}
	return true
}

// Multiply multiplies two matrices by each other and returns the result.
// Only applies to 4x4 matrices
func (m *Matrix) Multiply(m2 *Matrix) *Matrix {
	res := NewMatrix(newData(m.size))
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			iTuple := listToTuple(m.data[i])
			jTuple := NewTuple(m2.data[0][j], m2.data[1][j], m2.data[2][j], m2.data[3][j])
			res.data[i][j] = iTuple.DotProduct(jTuple)
		}
	}
	return res
}

// MultiplyTuple multiples a matrix by a tuple and returns the resulting tuple.
// Only applies to 4x4 matrices
func (m *Matrix) MultiplyTuple(t *Tuple) *Tuple {
	newVals := make([]float64, 4)
	for i, row := range m.data {
		newVals[i] = listToTuple(row).DotProduct(t)
	}
	return listToTuple(newVals)
}

// Transpose turns a matrix's rows into columns
func (m *Matrix) Transpose() *Matrix {
	res := NewMatrix(newData(m.size))
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			res.data[i][j] = m.data[j][i]
		}
	}
	return res
}

// Inverse returns the inversion of a matrix
func (m *Matrix) Inverse() (*Matrix, error) {
	if !m.invertible() {
		return nil, errors.New("cannot invert matrix with determinant of 0")
	}
	inverse := NewMatrix(newData(m.size))
	for i := 0; i < m.size; i++ {
		for j := 0; j < m.size; j++ {
			c := m.cofactor(i, j)
			inverse.data[j][i] = c / m.determinant()
		}
	}
	return inverse, nil
}

// returns whether or not a matrix can be inverted
func (m *Matrix) invertible() bool {
	return m.determinant() != 0
}

// Returns the determinant of a matrix
func (m *Matrix) determinant() float64 {
	if m.size == 2 {
		return m.data[0][0]*m.data[1][1] - m.data[0][1]*m.data[1][0]
	}
	var det float64
	for i, val := range m.data[0] {
		det += val * m.cofactor(0, i)
	}
	return det
}

// Returns a matrix with row and column removed
func (m *Matrix) subMatrix(row, col int) *Matrix {
	res := NewMatrix(newData(m.size - 1))
	var curRow int
	for i := 0; i < m.size; i++ {
		var curCol int
		if i == row {
			continue
		}
		for j := 0; j < m.size; j++ {
			if j == col {
				continue
			}
			res.data[curRow][curCol] = m.data[i][j]
			curCol++
		}
		curRow++
	}
	return res
}

// Returns the minor (determinant of the submatrix) of a matrix
func (m *Matrix) minor(row, col int) float64 {
	return m.subMatrix(row, col).determinant()
}

// Returns the cofactor of a matrix
func (m *Matrix) cofactor(row, col int) float64 {
	minor := m.minor(row, col)
	if (row+col)%2 == 0 {
		return minor
	}
	return -minor
}

// Converts a list of 4 values to a tuple
func listToTuple(list []float64) *Tuple {
	return NewTuple(list[0], list[1], list[2], list[3])
}
