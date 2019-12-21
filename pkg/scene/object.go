package scene

import (
	"math"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
	"github.com/sjberman/golang-ray-tracer/pkg/utils"
)

// Object is a generic object in a scene
type Object interface {
	GetMaterial() *Material
	GetTransform() *base.Matrix
	SetTransform(...*base.Matrix)
	SetMaterial(*Material)
	patternAt(*base.Tuple, image.Pattern) *image.Color
	intersect(*Ray) []*Intersection
	normalAt(*base.Tuple) *base.Tuple
}

// object is the base implementation of an Object
type object struct {
	Material
	transform base.Matrix
}

// newObject returns a new object
func newObject() *object {
	return &object{
		transform: base.Identity,
		Material:  DefaultMaterial,
	}
}

// GetTransform gets the Object's transform matrix
func (o *object) GetTransform() *base.Matrix {
	return &o.transform
}

// GetMaterial gets the Object's material
func (o *object) GetMaterial() *Material {
	return &o.Material
}

// SetTransform sets the Object's transform to the supplied matrix
func (o *object) SetTransform(matrix ...*base.Matrix) {
	t := base.Identity
	for _, m := range matrix {
		t = *t.Multiply(m)
	}
	o.transform = t
}

// SetMaterial sets the Object's material
func (o *object) SetMaterial(material *Material) {
	o.Material = *material
}

// patternAt returns the pattern at a point on the object
func (o *object) patternAt(worldPoint *base.Tuple, pattern image.Pattern) *image.Color {
	// convert the point from world space to object space
	objInverse, _ := o.GetTransform().Inverse()
	objectPoint := objInverse.MultiplyTuple(worldPoint)

	// convert point to pattern space
	patternInverse, _ := pattern.GetTransform().Inverse()
	patternPoint := patternInverse.MultiplyTuple(objectPoint)

	return pattern.PatternAt(patternPoint)
}

// transform the ray to the inverse of the object's transform;
// this is the same as transforming the object
func (o *object) transformRay(r *Ray) *Ray {
	objInverse, _ := o.GetTransform().Inverse()
	return r.Transform(objInverse)
}

// common normal function with the Object's specific calculation function passed in
func commonNormalAt(
	o Object,
	worldPoint *base.Tuple,
	objectNormalFunc func(*base.Tuple, Object) *base.Tuple,
) *base.Tuple {
	// convert the point from world space to object space
	// (object is likely not at the world origin)
	inverse, _ := o.GetTransform().Inverse()
	objectPoint := inverse.MultiplyTuple(worldPoint)

	objectNormal := objectNormalFunc(objectPoint, o)
	// convert normal back to world space
	worldNormal := inverse.Transpose().MultiplyTuple(objectNormal)
	// ensure this is a vector
	worldNormal.SetW(0)
	return worldNormal.Normalize()
}

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

// calculates where a ray intersects a sphere
func (s *Sphere) intersect(ray *Ray) []*Intersection {
	r := s.transformRay(ray)
	// sphere is centered at world origin
	sphereToRay, _ := r.Origin.Subtract(base.Origin)

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

	return intersections(NewIntersection(t1, s), NewIntersection(t2, s))
}

// wrapper for the normalAt interface function, using the common normal function
// with the specific sphere logic embedded
func (s *Sphere) normalAt(objectPoint *base.Tuple) *base.Tuple {
	return commonNormalAt(s, objectPoint, sphereNormal)
}

// sphere-specific calculation of the normal
func sphereNormal(objectPoint *base.Tuple, o Object) *base.Tuple {
	normal, _ := objectPoint.Subtract(base.Origin)
	return normal
}

// GlassSphere creates a glass sphere object
func GlassSphere() *Sphere {
	s := NewSphere()
	s.Transparency = 1
	s.RefractiveIndex = 1.5
	return s
}

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

// calculates where a ray intersects a plane
func (p *Plane) intersect(ray *Ray) []*Intersection {
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
func (p *Plane) normalAt(objectPoint *base.Tuple) *base.Tuple {
	return commonNormalAt(p, objectPoint, planeNormal)
}

// plane-specific calculation of the normal
func planeNormal(objectPoint *base.Tuple, o Object) *base.Tuple {
	return base.NewVector(0, 1, 0)
}

// Cube is a cube object
type Cube struct {
	*object
}

// NewCube returns a new Cube object
func NewCube() *Cube {
	return &Cube{
		object: newObject(),
	}
}

// calculates where a ray intersects a cube
func (c *Cube) intersect(ray *Ray) []*Intersection {
	r := c.transformRay(ray)
	// find largest minimum t value and smallest maximum t value for each axis
	// (t is intersection point)
	xtMin, xtMax := checkAxis(r.Origin.GetX(), r.Direction.GetX())
	ytMin, ytMax := checkAxis(r.Origin.GetY(), r.Direction.GetY())
	if xtMin > ytMax || ytMin > xtMax {
		return []*Intersection{}
	}
	ztMin, ztMax := checkAxis(r.Origin.GetZ(), r.Direction.GetZ())

	tMin := utils.Max(xtMin, ytMin, ztMin)
	tMax := utils.Min(xtMax, ytMax, ztMax)

	if tMin > tMax {
		return []*Intersection{}
	}
	return []*Intersection{NewIntersection(tMin, c), NewIntersection(tMax, c)}
}

// wrapper for the normalAt interface function, using the common normal function
// with the specific cube logic embedded
func (c *Cube) normalAt(objectPoint *base.Tuple) *base.Tuple {
	return commonNormalAt(c, objectPoint, cubeNormal)
}

// cube-specific calculation of the normal
func cubeNormal(objectPoint *base.Tuple, o Object) *base.Tuple {
	absX := math.Abs(objectPoint.GetX())
	absY := math.Abs(objectPoint.GetY())
	absZ := math.Abs(objectPoint.GetZ())
	maxC := utils.Max(absX, absY, absZ)
	if maxC == absX {
		return base.NewVector(objectPoint.GetX(), 0, 0)
	} else if maxC == absY {
		return base.NewVector(0, objectPoint.GetY(), 0)
	}
	return base.NewVector(0, 0, objectPoint.GetZ())
}

// checkAxis finds the min and max intersection values for the axis
func checkAxis(origin, direction float64) (float64, float64) {
	var tMin, tMax float64
	tMinNumerator := -1 - origin
	tMaxNumerator := 1 - origin

	if math.Abs(direction) >= base.Epsilon {
		tMin = tMinNumerator / direction
		tMax = tMaxNumerator / direction
	} else {
		// if denominator is effectively zero, multiply by infinity to ensure
		// values have the correct sign (positive or negative)
		tMin = tMinNumerator * math.Inf(0)
		tMax = tMaxNumerator * math.Inf(0)
	}

	if tMin > tMax {
		t := tMin
		tMin = tMax
		tMax = t
	}

	return tMin, tMax
}

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
func (cyl *Cylinder) intersect(ray *Ray) []*Intersection {
	r := cyl.transformRay(ray)
	// quadratic formula to determine intersection
	a := r.Direction.GetX()*r.Direction.GetX() + r.Direction.GetZ()*r.Direction.GetZ()
	if a <= base.Epsilon {
		// ray is parallel to y axis, check caps
		return intersections(cyl.intersectCaps(r, []*Intersection{})...)
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

	return intersections(cyl.intersectCaps(r, ints)...)
}

// wrapper for the normalAt interface function, using the common normal function
// with the specific cube logic embedded
func (cyl *Cylinder) normalAt(objectPoint *base.Tuple) *base.Tuple {
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
func checkCap(r *Ray, t float64) bool {
	x := r.Origin.GetX() + t*r.Direction.GetX()
	z := r.Origin.GetZ() + t*r.Direction.GetZ()
	return (x*x + z*z) <= 1
}

// checks if a ray intersects the caps of a closed cylinder
func (cyl *Cylinder) intersectCaps(r *Ray, ints []*Intersection) []*Intersection {
	// caps only matter if the cylinder is closed, and might possibly be intersected
	if !cyl.Closed || math.Abs(r.Direction.GetY()) < base.Epsilon {
		return ints
	}

	// check for intersection with lower cap by intersecting ray w/ plane at y = minimum
	t := (cyl.Minimum - r.Origin.GetY()) / r.Direction.GetY()
	if checkCap(r, t) {
		ints = append(ints, NewIntersection(t, cyl))
	}
	// check for intersection with upper cap by intersecting ray w/ plane at y = maximum
	t = (cyl.Maximum - r.Origin.GetY()) / r.Direction.GetY()
	if checkCap(r, t) {
		ints = append(ints, NewIntersection(t, cyl))
	}
	return ints
}
