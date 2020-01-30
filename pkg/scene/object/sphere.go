package object

import (
	"math"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

// Sphere is a sphere object
type Sphere struct {
	*object
}

// NewSphere returns a new Sphere object
func NewSphere() *Sphere {
	return &Sphere{
		object: newObject(),
	}
}

// Bounds returns the untransformed bounds of a sphere
func (s *Sphere) Bounds() *Bounds {
	return &Bounds{
		Minimum: base.NewPoint(-1, -1, -1),
		Maximum: base.NewPoint(1, 1, 1),
	}
}

// calculates where a ray intersects a sphere
func (s *Sphere) Intersect(ray *ray.Ray) []*Intersection {
	r := s.transformRay(ray)
	// sphere is centered at world origin
	sphereToRay := r.Origin.Subtract(base.Origin)

	// quadratic formula to determine intersection
	a := r.Direction.DotProduct(r.Direction)
	b := 2 * r.Direction.DotProduct(sphereToRay)
	c := sphereToRay.DotProduct(sphereToRay) - 1

	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return []*Intersection{}
	}
	t1 := (-b - math.Sqrt(discriminant)) / (2 * a)
	t2 := (-b + math.Sqrt(discriminant)) / (2 * a)

	return Intersections(NewIntersection(t1, s), NewIntersection(t2, s))
}

// wrapper for the normalAt interface function, using the common normal function
// with the specific sphere logic embedded
func (s *Sphere) NormalAt(objectPoint *base.Tuple, hit *Intersection) *base.Tuple {
	return commonNormalAt(s, objectPoint, hit, sphereNormal)
}

// sphere-specific calculation of the normal
func sphereNormal(objectPoint *base.Tuple, _ Object, _ *Intersection) *base.Tuple {
	return objectPoint.Subtract(base.Origin)
}

// GlassSphere creates a glass sphere object
func GlassSphere() *Sphere {
	s := NewSphere()
	s.Diffuse = 0.1
	s.Transparency = 1
	s.Reflective = 0.9
	s.RefractiveIndex = 1.5
	s.Specular = 1
	s.Shininess = 300
	return s
}
