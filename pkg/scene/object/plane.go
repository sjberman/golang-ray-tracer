package object

import (
	"math"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

// Plane is a plane object
type Plane struct {
	*object
}

// NewPlane returns a new Plane object
func NewPlane() *Plane {
	return &Plane{
		object: newObject(),
	}
}

// Bounds returns the untransformed bounds of a plane
func (p *Plane) Bounds() *Bounds {
	return &Bounds{
		Minimum: base.NewPoint(math.Inf(-1), 0, math.Inf(-1)),
		Maximum: base.NewPoint(math.Inf(1), 0, math.Inf(1)),
	}
}

// calculates where a ray intersects a plane
func (p *Plane) Intersect(ray *ray.Ray) []*Intersection {
	r := p.transformRay(ray)
	// parallel to plane (y == 0)
	if math.Abs(r.Direction.GetY()) < base.Epsilon {
		return []*Intersection{}
	}

	t := -r.Origin.GetY() / r.Direction.GetY()
	return []*Intersection{NewIntersection(t, p)}
}

// wrapper for the normalAt interface function, using the common normal function
// with the specific plane logic embedded
func (p *Plane) NormalAt(objectPoint *base.Tuple, hit *Intersection) *base.Tuple {
	return commonNormalAt(p, objectPoint, hit, planeNormal)
}

// plane-specific calculation of the normal
func planeNormal(objectPoint *base.Tuple, _ Object, _ *Intersection) *base.Tuple {
	return base.NewVector(0, 1, 0)
}
