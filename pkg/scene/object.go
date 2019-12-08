package scene

import (
	"math"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
)

// Object is a generic object in a scene
type Object interface {
	GetMaterial() *Material
	GetTransform() *base.Matrix
	SetTransform(matrix *base.Matrix)
	SetMaterial(material *Material)
	intersect(*Ray) []*Intersection
	normalAt(worldPoint *base.Tuple) *base.Tuple
}

// object is the base implementation of the Object interface
type object struct {
	transform        base.Matrix
	material         Material
	intersectFunc    func(*Ray, *object) []*Intersection
	objectNormalFunc func(point *base.Tuple) *base.Tuple
}

// newObject returns a new Object
func newObject(
	intersectFunc func(*Ray, *object) []*Intersection,
	objectNormalFunc func(*base.Tuple) *base.Tuple,
) *object {
	return &object{
		transform:        base.Identity,
		material:         defaultMaterial,
		intersectFunc:    intersectFunc,
		objectNormalFunc: objectNormalFunc,
	}
}

// GetTransform gets the Object's transform matrix
func (o *object) GetTransform() *base.Matrix {
	return &o.transform
}

// GetMaterial gets the Object's material
func (o *object) GetMaterial() *Material {
	return &o.material
}

// SetTransform sets the Object's transform to the supplied matrix
func (o *object) SetTransform(matrix *base.Matrix) {
	o.transform = *matrix
}

// SetMaterial sets the Object's material
func (o *object) SetMaterial(material *Material) {
	o.material = *material
}

// calculates where a ray intersects the object
func (o *object) intersect(r *Ray) []*Intersection {
	// transform the ray to the inverse of the object's transform;
	// this is the same as transforming the object
	objInverse, _ := o.GetTransform().Inverse()
	newRay := r.Transform(objInverse)

	return o.intersectFunc(newRay, o)
}

// normalAt returns the surface normal at a position on the object
func (o *object) normalAt(worldPoint *base.Tuple) *base.Tuple {
	// convert the point from world space to object space
	// (object is likely not at the world origin)
	inverse, _ := o.transform.Inverse()
	objectPoint := inverse.MultiplyTuple(worldPoint)

	objectNormal := o.objectNormalFunc(objectPoint)
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
		object: newObject(sphereIntersect, sphereObjectNormal),
	}
}

// calculates where a ray intersects a sphere
func sphereIntersect(r *Ray, o *object) []*Intersection {
	// sphere is centered at world origin
	sphereToRay, _ := r.origin.Subtract(base.Origin)

	// quadratic formula to determine intersection
	a := r.direction.DotProduct(r.direction)
	b := 2 * r.direction.DotProduct(sphereToRay)
	c := sphereToRay.DotProduct(sphereToRay) - 1

	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return []*Intersection{}
	}
	t1 := (-b - math.Sqrt(discriminant)) / (2 * a)
	t2 := (-b + math.Sqrt(discriminant)) / (2 * a)

	return intersections(NewIntersection(t1, o), NewIntersection(t2, o))
}

// computes the object normal for a sphere
func sphereObjectNormal(objectPoint *base.Tuple) *base.Tuple {
	normal, _ := objectPoint.Subtract(base.Origin)
	return normal
}
