package object

import (
	"math"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

// Cylinder is a cylinder object
type Cylinder struct {
	*object
	Minimum, Maximum float64
	Closed           bool
}

// NewCylinder returns a new Cylinder object
func NewCylinder() *Cylinder {
	return &Cylinder{
		object:  newObject(),
		Minimum: math.Inf(-1),
		Maximum: math.Inf(0),
		Closed:  false,
	}
}

// calculates where a ray intersects a cylinder
func (cyl *Cylinder) Intersect(ray *ray.Ray) []*Intersection {
	r := cyl.transformRay(ray)
	// quadratic formula to determine intersection
	a := r.Direction.GetX()*r.Direction.GetX() + r.Direction.GetZ()*r.Direction.GetZ()
	if a <= base.Epsilon {
		// ray is parallel to y axis, check caps
		return Intersections(cyl.intersectCaps(r, []*Intersection{})...)
	}
	b := 2 * (r.Origin.GetX()*r.Direction.GetX() + r.Origin.GetZ()*r.Direction.GetZ())
	c := r.Origin.GetX()*r.Origin.GetX() + r.Origin.GetZ()*r.Origin.GetZ() - 1

	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return []*Intersection{}
	}

	t1 := (-b - math.Sqrt(discriminant)) / (2 * a)
	t2 := (-b + math.Sqrt(discriminant)) / (2 * a)

	if t1 > t2 {
		t := t1
		t1 = t2
		t2 = t
	}

	ints := []*Intersection{}

	y0 := r.Origin.GetY() + t1*r.Direction.GetY()
	if cyl.Minimum < y0 && y0 < cyl.Maximum {
		ints = append(ints, NewIntersection(t1, cyl))
	}
	y1 := r.Origin.GetY() + t2*r.Direction.GetY()
	if cyl.Minimum < y1 && y1 < cyl.Maximum {
		ints = append(ints, NewIntersection(t2, cyl))
	}

	return Intersections(cyl.intersectCaps(r, ints)...)
}

// wrapper for the normalAt interface function, using the common normal function
// with the specific cylinder logic embedded
func (cyl *Cylinder) NormalAt(objectPoint *base.Tuple) *base.Tuple {
	return commonNormalAt(cyl, objectPoint, cylinderNormal)
}

// cylinder-specific calculation of the normal
func cylinderNormal(objectPoint *base.Tuple, o Object) *base.Tuple {
	cyl := o.(*Cylinder)
	// compute the square of the distance from the y axis
	distance := objectPoint.GetX()*objectPoint.GetX() + objectPoint.GetZ()*objectPoint.GetZ()
	if distance < 1 && objectPoint.GetY() >= cyl.Maximum-base.Epsilon {
		return base.NewVector(0, 1, 0)
	} else if distance < 1 && objectPoint.GetY() <= cyl.Minimum+base.Epsilon {
		return base.NewVector(0, -1, 0)
	}
	return base.NewVector(objectPoint.GetX(), 0, objectPoint.GetZ())
}

// checks to see if the intersection at t is within a radius of 1 from the y axis
func (cyl *Cylinder) checkCap(r *ray.Ray, t float64) bool {
	x := r.Origin.GetX() + t*r.Direction.GetX()
	z := r.Origin.GetZ() + t*r.Direction.GetZ()
	return (x*x + z*z) <= 1
}

// checks if a ray intersects the caps of a closed cylinder
func (cyl *Cylinder) intersectCaps(r *ray.Ray, ints []*Intersection) []*Intersection {
	// caps only matter if the cylinder is closed, and might possibly be intersected
	if !cyl.Closed || math.Abs(r.Direction.GetY()) < base.Epsilon {
		return ints
	}

	// check for intersection with lower cap by intersecting ray w/ plane at y = minimum
	t := (cyl.Minimum - r.Origin.GetY()) / r.Direction.GetY()
	if cyl.checkCap(r, t) {
		ints = append(ints, NewIntersection(t, cyl))
	}
	// check for intersection with upper cap by intersecting ray w/ plane at y = maximum
	t = (cyl.Maximum - r.Origin.GetY()) / r.Direction.GetY()
	if cyl.checkCap(r, t) {
		ints = append(ints, NewIntersection(t, cyl))
	}
	return ints
}
