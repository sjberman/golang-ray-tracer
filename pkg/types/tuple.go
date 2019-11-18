package types

import (
	"errors"
	"math"
)

// tupleType is an integer representing a vector or point
type tupleType int

const (
	vector tupleType = iota
	point
)

// Tuple is a 3D coordinate (either vector or point)
type Tuple struct {
	xAxis     float64
	yAxis     float64
	zAxis     float64
	tupleType tupleType
}

// CreateVector returns a tuple object of type Vector
func CreateVector(x, y, z float64) *Tuple {
	return &Tuple{
		xAxis:     x,
		yAxis:     y,
		zAxis:     z,
		tupleType: vector,
	}
}

// CreatePoint returns a tuple object of type Point
func CreatePoint(x, y, z float64) *Tuple {
	return &Tuple{
		xAxis:     x,
		yAxis:     y,
		zAxis:     z,
		tupleType: point,
	}
}

// IsVector returns whether or not a tuple is a vector
func (t *Tuple) IsVector() bool {
	return t.tupleType == vector
}

// IsPoint returns whether or not a tuple is a point
func (t *Tuple) IsPoint() bool {
	return t.tupleType == point
}

// Add adds two tuples together and returns the result
func (t *Tuple) Add(t2 *Tuple) (*Tuple, error) {
	if t.IsPoint() && t2.IsPoint() {
		return nil, errors.New("cannot add two points together")
	}
	return &Tuple{
		xAxis:     t.xAxis + t2.xAxis,
		yAxis:     t.yAxis + t2.yAxis,
		zAxis:     t.zAxis + t2.zAxis,
		tupleType: t.tupleType + t2.tupleType,
	}, nil
}

// Subtract returns the difference between two tuples
func (t *Tuple) Subtract(t2 *Tuple) (*Tuple, error) {
	if t.IsVector() && t2.IsPoint() {
		return nil, errors.New("cannot subtract a point from a vector")
	}
	return &Tuple{
		xAxis:     t.xAxis - t2.xAxis,
		yAxis:     t.yAxis - t2.yAxis,
		zAxis:     t.zAxis - t2.zAxis,
		tupleType: t.tupleType - t2.tupleType,
	}, nil
}

// Multiply returns a tuple multiplied by a value
func (t *Tuple) Multiply(val float64) *Tuple {
	return &Tuple{
		xAxis:     t.xAxis * val,
		yAxis:     t.yAxis * val,
		zAxis:     t.zAxis * val,
		tupleType: t.tupleType,
	}
}

// Divide returns a tuple divided by a value
func (t *Tuple) Divide(val float64) *Tuple {
	return &Tuple{
		xAxis:     t.xAxis / val,
		yAxis:     t.yAxis / val,
		zAxis:     t.zAxis / val,
		tupleType: t.tupleType,
	}
}

var epsilon = math.Nextafter(1, 2) - 1

func equalFloats(one, two float64) bool {
	return math.Abs(one-two) <= epsilon
}

// Equals returns whether or not two tuples are equal to each other
func (t *Tuple) Equals(t2 *Tuple) bool {
	if !equalFloats(t.xAxis, t2.xAxis) {
		return false
	}
	if !equalFloats(t.yAxis, t2.yAxis) {
		return false
	}
	if !equalFloats(t.zAxis, t2.zAxis) {
		return false
	}
	return t.tupleType == t2.tupleType
}

// Negate returns the calling tuple with its fields negated
func (t *Tuple) Negate() *Tuple {
	return &Tuple{
		xAxis:     -t.xAxis,
		yAxis:     -t.yAxis,
		zAxis:     -t.zAxis,
		tupleType: t.tupleType,
	}
}

// Magnitude returns the length of a vector (using Euclidean distance formula)
func (t *Tuple) Magnitude() float64 {
	// TODO: error check to disallow point
	a := t.xAxis * t.xAxis
	b := t.yAxis * t.yAxis
	c := t.zAxis * t.zAxis
	return math.Sqrt(a + b + c)
}

// Normalize converts a vector into a unit vector (magnitude of 1)
func (t *Tuple) Normalize() *Tuple {
	// TODO: error check to disallow point
	return t.Divide(t.Magnitude())
}

// DotProduct returns the dot product of two vectors
func (t *Tuple) DotProduct(t2 *Tuple) float64 {
	// TODO: error check to disallow point
	a := t.xAxis * t2.xAxis
	b := t.yAxis * t2.yAxis
	c := t.zAxis * t2.zAxis
	return a + b + c
}

// CrossProduct returns the cross product of two vectors
func (t *Tuple) CrossProduct(t2 *Tuple) *Tuple {
	// TODO: error check to disallow point
	newX := (t.yAxis * t2.zAxis) - (t.zAxis * t2.yAxis)
	newY := (t.zAxis * t2.xAxis) - (t.xAxis * t2.zAxis)
	newZ := (t.xAxis * t2.yAxis) - (t.yAxis * t2.xAxis)
	return CreateVector(newX, newY, newZ)
}
