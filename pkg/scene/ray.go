package scene

import (
	"github.com/sjberman/golang-ray-tracer/pkg/base"
)

// Ray is a light ray with an origin and direction
type Ray struct {
	origin    *base.Tuple
	direction *base.Tuple
}

// NewRay returns a new Ray object
func NewRay(point *base.Tuple, vector *base.Tuple) *Ray {
	return &Ray{
		origin:    point,
		direction: vector,
	}
}

// GetDirection returns the ray's direction
func (r *Ray) GetDirection() *base.Tuple {
	return r.direction
}

// Position returns the point at a given distance along the ray
func (r *Ray) Position(distance float64) *base.Tuple {
	sum, _ := r.origin.Add(r.direction.Multiply(distance))
	return sum
}

// Transform applies the transformation matrix to the ray
func (r *Ray) Transform(matrix *base.Matrix) *Ray {
	return NewRay(matrix.MultiplyTuple(r.origin), matrix.MultiplyTuple(r.direction))
}
