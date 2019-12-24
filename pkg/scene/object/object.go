package object

import (
	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
	"github.com/sjberman/golang-ray-tracer/pkg/utils"
)

// Object is a generic object in a scene
type Object interface {
	GetMaterial() *Material
	GetTransform() *base.Matrix
	GetParent() *Group
	Bounds() *bounds
	SetTransform(...*base.Matrix)
	SetMaterial(*Material)
	SetParent(*Group)
	PatternAt(*base.Tuple, image.Pattern) *image.Color
	Intersect(*ray.Ray) []*Intersection
	NormalAt(*base.Tuple) *base.Tuple
	worldToObject(*base.Tuple) *base.Tuple
	normalToWorld(*base.Tuple) *base.Tuple
}

// object is the base implementation of an Object
type object struct {
	Material
	transform base.Matrix
	parent    *Group
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

// GetParent gets the Object's parent group
func (o *object) GetParent() *Group {
	return o.parent
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

// SetParent sets the Object's parent group
func (o *object) SetParent(group *Group) {
	o.parent = group
}

// patternAt returns the pattern at a point on the object
func (o *object) PatternAt(worldPoint *base.Tuple, pattern image.Pattern) *image.Color {
	// convert the point from world space to object space
	objectPoint := o.worldToObject(worldPoint)

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
	objectPoint := o.worldToObject(worldPoint)
	objectNormal := objectNormalFunc(objectPoint, o)
	return o.normalToWorld(objectNormal)
}

// converts a point in world space to object space
func (o *object) worldToObject(point *base.Tuple) *base.Tuple {
	if o.parent != nil {
		point = o.parent.worldToObject(point)
	}
	inverse, _ := o.transform.Inverse()
	return inverse.MultiplyTuple(point)
}

// converts a normal in object space to world space
func (o *object) normalToWorld(normal *base.Tuple) *base.Tuple {
	inverse, _ := o.GetTransform().Inverse()
	normal = inverse.Transpose().MultiplyTuple(normal)
	normal.SetW(0)
	normal = normal.Normalize()

	if o.parent != nil {
		normal = o.parent.normalToWorld(normal)
	}
	return normal
}

// bounds represents a bounding box for an object
type bounds struct {
	minimum *base.Tuple
	maximum *base.Tuple
}

// use cube intersect method
func (b *bounds) intersects(r *ray.Ray) bool {
	// find largest minimum t value and smallest maximum t value for each axis
	// (t is intersection point)
	xtMin, xtMax := checkAxis(r.Origin.GetX(), r.Direction.GetX(), b.minimum.GetX(), b.maximum.GetX())
	ytMin, ytMax := checkAxis(r.Origin.GetY(), r.Direction.GetY(), b.minimum.GetY(), b.maximum.GetY())
	if xtMin > ytMax || ytMin > xtMax {
		return false
	}
	ztMin, ztMax := checkAxis(r.Origin.GetZ(), r.Direction.GetZ(), b.minimum.GetZ(), b.maximum.GetZ())

	tMin := utils.Max(xtMin, ytMin, ztMin)
	tMax := utils.Min(xtMax, ytMax, ztMax)
	return tMin <= tMax
}
