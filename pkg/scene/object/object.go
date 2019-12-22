package object

import (
	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

// Object is a generic object in a scene
type Object interface {
	GetMaterial() *Material
	GetTransform() *base.Matrix
	SetTransform(...*base.Matrix)
	SetMaterial(*Material)
	PatternAt(*base.Tuple, image.Pattern) *image.Color
	Intersect(*ray.Ray) []*Intersection
	NormalAt(*base.Tuple) *base.Tuple
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
func (o *object) PatternAt(worldPoint *base.Tuple, pattern image.Pattern) *image.Color {
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
func (o *object) transformRay(r *ray.Ray) *ray.Ray {
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
