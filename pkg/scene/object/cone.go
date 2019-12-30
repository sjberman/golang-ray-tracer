package object

import (
	"math"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

// Cone is a cone object
type Cone struct {
	*object
	*Cylinder
	Minimum, Maximum float64
	Closed           bool
}

// NewCylinder returns a new Cone object
func NewCone() *Cone {
	return &Cone{
		object:  newObject(),
		Minimum: math.Inf(-1),
		Maximum: math.Inf(0),
		Closed:  false,
	}
}

// Bounds returns the untransformed bounds of a cone
func (cone *Cone) Bounds() *Bounds {
	return &Bounds{
		Minimum: base.NewPoint(-1, cone.Minimum, -1),
		Maximum: base.NewPoint(1, cone.Maximum, 1),
	}
}

// calculates where a ray intersects a cone
func (cone *Cone) Intersect(ray *ray.Ray) []*Intersection {
	r := cone.transformRay(ray)
	// quadratic formula to determine intersection
	dx, ox := r.Direction.GetX(), r.Origin.GetX()
	dy, oy := r.Direction.GetY(), r.Origin.GetY()
	dz, oz := r.Direction.GetZ(), r.Origin.GetZ()
	a := dx*dx - dy*dy + dz*dz
	b := 2*ox*dx - 2*oy*dy + 2*oz*dz
	c := ox*ox - oy*oy + oz*oz

	ints := []*Intersection{}
	if math.Abs(a) <= base.Epsilon && math.Abs(b) > base.Epsilon {
		// hits the tip
		ints = append(ints, NewIntersection(-c/(2*b), cone))
	}
	if math.Abs(a) <= base.Epsilon && math.Abs(b) <= base.Epsilon {
		// ray misses
		return []*Intersection{}
	}

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

	y0 := r.Origin.GetY() + t1*r.Direction.GetY()
	if cone.Minimum < y0 && y0 < cone.Maximum {
		ints = append(ints, NewIntersection(t1, cone))
	}
	y1 := r.Origin.GetY() + t2*r.Direction.GetY()
	if cone.Minimum < y1 && y1 < cone.Maximum {
		ints = append(ints, NewIntersection(t2, cone))
	}

	return Intersections(cone.intersectCaps(r, ints)...)
}

// wrapper for the normalAt interface function, using the common normal function
// with the specific cone logic embedded
func (cone *Cone) NormalAt(objectPoint *base.Tuple, hit *Intersection) *base.Tuple {
	return commonNormalAt(cone, objectPoint, hit, coneNormal)
}

// cone-specific calculation of the normal
func coneNormal(objectPoint *base.Tuple, o Object, _ *Intersection) *base.Tuple {
	cone := o.(*Cone)
	x := objectPoint.GetX()
	y := objectPoint.GetY()
	z := objectPoint.GetZ()
	// compute the square of the distance from the y axis
	distance := x*x + z*z
	if distance < 1 && y >= cone.Maximum-base.Epsilon {
		return base.NewVector(0, 1, 0)
	} else if distance < 1 && y <= cone.Minimum+base.Epsilon {
		return base.NewVector(0, -1, 0)
	}
	yNormal := math.Sqrt(distance)
	if y > 0 {
		yNormal *= -1
	}
	return base.NewVector(x, yNormal, z)
}

// checks to see if the intersection at t is within a radius of y*y from the y axis
func (cone *Cone) checkCap(r *ray.Ray, t, y float64) bool {
	x := r.Origin.GetX() + t*r.Direction.GetX()
	z := r.Origin.GetZ() + t*r.Direction.GetZ()
	return (x*x + z*z) <= y*y
}

// checks if a ray intersects the caps of a closed cone
func (cone *Cone) intersectCaps(r *ray.Ray, ints []*Intersection) []*Intersection {
	// caps only matter if the cylinder is closed, and might possibly be intersected
	if !cone.Closed || math.Abs(r.Direction.GetY()) < base.Epsilon {
		return ints
	}

	// check for intersection with lower cap by intersecting ray w/ plane at y = minimum
	t := (cone.Minimum - r.Origin.GetY()) / r.Direction.GetY()
	if cone.checkCap(r, t, cone.Minimum) {
		ints = append(ints, NewIntersection(t, cone))
	}
	// check for intersection with upper cap by intersecting ray w/ plane at y = maximum
	t = (cone.Maximum - r.Origin.GetY()) / r.Direction.GetY()
	if cone.checkCap(r, t, cone.Maximum) {
		ints = append(ints, NewIntersection(t, cone))
	}
	return ints
}
