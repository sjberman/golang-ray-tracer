package scene

import (
	"math"

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

// Intersect returns the positions in which a ray intersects a sphere
func (r *Ray) Intersect(s *Sphere) []*Intersection {
	// transform the ray to the inverse of the spheres transform;
	// this is the same as transforming the sphere
	sphereInverse, _ := s.transform.Inverse()
	newRay := r.Transform(sphereInverse)
	// sphere is centered at world origin
	sphereToRay, _ := newRay.origin.Subtract(base.Origin)

	// quadratic formula to determine intersection
	a := newRay.direction.DotProduct(newRay.direction)
	b := 2 * newRay.direction.DotProduct(sphereToRay)
	c := sphereToRay.DotProduct(sphereToRay) - 1

	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return []*Intersection{}
	}
	t1 := (-b - math.Sqrt(discriminant)) / (2 * a)
	t2 := (-b + math.Sqrt(discriminant)) / (2 * a)

	return intersections(NewIntersection(t1, s), NewIntersection(t2, s))
}

// Transform applies the transformation matrix to the ray
func (r *Ray) Transform(matrix *base.Matrix) *Ray {
	return NewRay(matrix.MultiplyTuple(r.origin), matrix.MultiplyTuple(r.direction))
}
