package base

import (
	"errors"
	"math"
)

const (
	vector = 0
	point  = 1
)

// Tuple is a 3D coordinate (either vector or point)
type Tuple struct {
	xAxis float64
	yAxis float64
	zAxis float64
	w     float64
}

// Origin represents the origin point
var Origin = &Tuple{
	xAxis: 0,
	yAxis: 0,
	zAxis: 0,
	w:     1,
}

// NewTuple returns a generic tuple
func NewTuple(x, y, z, w float64) *Tuple {
	return &Tuple{
		xAxis: roundValue(x),
		yAxis: roundValue(y),
		zAxis: roundValue(z),
		w:     w,
	}
}

// NewVector returns a tuple object of type Vector
func NewVector(x, y, z float64) *Tuple {
	return NewTuple(x, y, z, vector)
}

// NewPoint returns a tuple object of type Point
func NewPoint(x, y, z float64) *Tuple {
	return NewTuple(x, y, z, point)
}

// IsVector returns whether or not a tuple is a vector
func (t *Tuple) IsVector() bool {
	return t.w == vector
}

// IsPoint returns whether or not a tuple is a point
func (t *Tuple) IsPoint() bool {
	return t.w == point
}

// GetX() returns the x coordinate of the tuple
func (t *Tuple) GetX() float64 {
	return t.xAxis
}

// GetY() returns the y coordinate of the tuple
func (t *Tuple) GetY() float64 {
	return t.yAxis
}

// GetZ() returns the z coordinate of the tuple
func (t *Tuple) GetZ() float64 {
	return t.zAxis
}

// SetW() sets the w value
func (t *Tuple) SetW(val float64) {
	t.w = val
}

// Add adds two tuples together and returns the result
func (t *Tuple) Add(t2 *Tuple) (*Tuple, error) {
	if t.IsPoint() && t2.IsPoint() {
		return nil, errors.New("cannot add two points together")
	}
	return NewTuple(t.xAxis+t2.xAxis, t.yAxis+t2.yAxis, t.zAxis+t2.zAxis, t.w+t2.w), nil
}

// Subtract returns the difference between two tuples
func (t *Tuple) Subtract(t2 *Tuple) (*Tuple, error) {
	if t.IsVector() && t2.IsPoint() {
		return nil, errors.New("cannot subtract a point from a vector")
	}
	return NewTuple(t.xAxis-t2.xAxis, t.yAxis-t2.yAxis, t.zAxis-t2.zAxis, t.w-t2.w), nil
}

// Multiply returns a tuple multiplied by a value
func (t *Tuple) Multiply(val float64) *Tuple {
	return NewTuple(t.xAxis*val, t.yAxis*val, t.zAxis*val, t.w)
}

// Divide returns a tuple divided by a value
func (t *Tuple) Divide(val float64) *Tuple {
	return NewTuple(t.xAxis/val, t.yAxis/val, t.zAxis/val, t.w)
}

// Epsilon is the +/- value for floating point algebra to be considered equal
var Epsilon = 0.00001

// EqualFloats uses approximation to determine if two floats are equivalent
func EqualFloats(one, two float64) bool {
	return math.Abs(one-two) <= Epsilon
}

// Equals returns whether or not two tuples are equal to each other
func (t *Tuple) Equals(t2 *Tuple) bool {
	if !EqualFloats(t.xAxis, t2.xAxis) {
		return false
	}
	if !EqualFloats(t.yAxis, t2.yAxis) {
		return false
	}
	if !EqualFloats(t.zAxis, t2.zAxis) {
		return false
	}
	return t.w == t2.w
}

// Negate returns the calling tuple with its fields negated
func (t *Tuple) Negate() *Tuple {
	return NewTuple(-t.xAxis, -t.yAxis, -t.zAxis, t.w)
}

// Magnitude returns the length of a vector (using Euclidean distance formula)
func (t *Tuple) Magnitude() float64 {
	a := t.xAxis * t.xAxis
	b := t.yAxis * t.yAxis
	c := t.zAxis * t.zAxis
	d := t.w * t.w
	return math.Sqrt(a + b + c + d)
}

// Normalize converts a vector into a unit vector (magnitude of 1)
func (t *Tuple) Normalize() *Tuple {
	magnitude := t.Magnitude()
	if magnitude != 0 {
		return t.Divide(magnitude)
	}
	return t
}

// DotProduct returns the dot product of two tuples
func (t *Tuple) DotProduct(t2 *Tuple) float64 {
	a := t.xAxis * t2.xAxis
	b := t.yAxis * t2.yAxis
	c := t.zAxis * t2.zAxis
	d := t.w * t2.w
	return a + b + c + d
}

// CrossProduct returns the cross product of two vectors
func (t *Tuple) CrossProduct(t2 *Tuple) *Tuple {
	newX := (t.yAxis * t2.zAxis) - (t.zAxis * t2.yAxis)
	newY := (t.zAxis * t2.xAxis) - (t.xAxis * t2.zAxis)
	newZ := (t.xAxis * t2.yAxis) - (t.yAxis * t2.xAxis)
	return NewVector(newX, newY, newZ)
}

// Reflect returns the reflection vector around a normal vector
func (t *Tuple) Reflect(normal *Tuple) *Tuple {
	reflection, _ := t.Subtract(normal.Multiply(2).Multiply(t.DotProduct(normal)))
	return reflection
}

// Converts a list of 4 values to a tuple
func listToTuple(list []float64) *Tuple {
	return NewTuple(list[0], list[1], list[2], list[3])
}

// if value is within epsilon of its floor or ceiling, round it
func roundValue(val float64) float64 {
	ret := val
	if val-math.Floor(val) <= Epsilon {
		ret = math.Floor(val)
	} else if math.Ceil(val)-val <= Epsilon {
		ret = math.Ceil(val)
	}
	return ret
}
