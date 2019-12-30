package ray

import "github.com/sjberman/golang-ray-tracer/pkg/base"

// Ray is a light ray with an origin and direction
type Ray struct {
	Origin    *base.Tuple
	Direction *base.Tuple
}

// NewRay returns a new Ray object
func NewRay(point *base.Tuple, vector *base.Tuple) *Ray {
	return &Ray{
		Origin:    point,
		Direction: vector,
	}
}

// Position returns the point at a given distance along the ray
func (r *Ray) Position(distance float64) *base.Tuple {
	return r.Origin.Add(r.Direction.Multiply(distance))
}

// Transform applies the transformation matrix to the ray
func (r *Ray) Transform(matrix *base.Matrix) *Ray {
	return NewRay(matrix.MultiplyTuple(r.Origin), matrix.MultiplyTuple(r.Direction))
}
