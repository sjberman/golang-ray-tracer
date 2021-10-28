package object

import (
	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
	"github.com/sjberman/golang-ray-tracer/pkg/utils"
)

// Object is a generic object in a scene.
type Object interface {
	GetMaterial() *Material
	GetTransform() *base.Matrix
	GetParent() Object
	Bounds() *Bounds
	SetTransform(...*base.Matrix)
	SetMaterial(*Material)
	SetParent(Object)
	PatternAt(*base.Tuple, image.Pattern) *image.Color
	Intersect(*ray.Ray) []*Intersection
	NormalAt(*base.Tuple, *Intersection) *base.Tuple
	Divide(int)
	worldToObject(*base.Tuple) *base.Tuple
	normalToWorld(*base.Tuple) *base.Tuple
	DeepCopy() Object
}

// object is the base implementation of an Object.
type object struct {
	Material
	transform base.Matrix
	parent    Object
}

// newObject returns a new object.
func newObject() *object {
	return &object{
		transform: base.Identity,
		Material:  DefaultMaterial,
	}
}

// GetTransform gets the Object's transform matrix.
func (o *object) GetTransform() *base.Matrix {
	return &o.transform
}

// GetMaterial gets the Object's material.
func (o *object) GetMaterial() *Material {
	return &o.Material
}

// GetParent gets the Object's parent Object.
func (o *object) GetParent() Object {
	return o.parent
}

// SetTransform sets the Object's transform to the supplied matrix.
func (o *object) SetTransform(matrix ...*base.Matrix) {
	t := base.Identity
	for _, m := range matrix {
		t = *t.Multiply(m)
	}
	o.transform = t
}

// SetMaterial sets the Object's material.
func (o *object) SetMaterial(material *Material) {
	o.Material = *material
}

// SetParent sets the Object's parent Object.
func (o *object) SetParent(obj Object) {
	o.parent = obj
}

// patternAt returns the pattern at a point on the object.
func (o *object) PatternAt(worldPoint *base.Tuple, pattern image.Pattern) *image.Color {
	// convert the point from world space to object space
	objectPoint := o.worldToObject(worldPoint)

	// convert point to pattern space
	patternInverse := pattern.GetTransform().Inverse()
	patternPoint := patternInverse.MultiplyTuple(objectPoint)

	return pattern.PatternAt(patternPoint)
}

// transform the ray to the inverse of the object's transform;
// this is the same as transforming the object.
func (o *object) transformRay(r *ray.Ray) *ray.Ray {
	objInverse := o.GetTransform().Inverse()

	return r.Transform(objInverse)
}

// common normal function with the Object's specific calculation function passed in.
func commonNormalAt(
	o Object,
	worldPoint *base.Tuple,
	hit *Intersection,
	objectNormalFunc func(*base.Tuple, Object, *Intersection) *base.Tuple,
) *base.Tuple {
	objectPoint := o.worldToObject(worldPoint)
	objectNormal := objectNormalFunc(objectPoint, o, hit)

	return o.normalToWorld(objectNormal)
}

// converts a point in world space to object space.
func (o *object) worldToObject(point *base.Tuple) *base.Tuple {
	if o.parent != nil {
		point = o.parent.worldToObject(point)
	}
	inverse := o.transform.Inverse()

	return inverse.MultiplyTuple(point)
}

// converts a normal in object space to world space.
func (o *object) normalToWorld(normal *base.Tuple) *base.Tuple {
	inverse := o.GetTransform().Inverse()
	normal = inverse.Transpose().MultiplyTuple(normal)
	normal.SetW(0)
	normal = normal.Normalize()

	if o.parent != nil {
		normal = o.parent.normalToWorld(normal)
	}

	return normal
}

// returns whether or not A contains B.
func includes(a Object, b Object) bool {
	if grp, ok := a.(*Group); ok {
		for _, o := range grp.Objects {
			if includes(o, b) {
				return true
			}
		}

		return false
	}
	if csg, ok := a.(*Csg); ok {
		return includes(csg.left, b) || includes(csg.right, b)
	}

	return a == b
}

// Remove removes an object.Object from a slice of Objects.
func Remove(s []Object, o Object) []Object {
	for i, obj := range s {
		if obj == o {
			copy(s[i:], s[i+1:])
			s[len(s)-1] = nil
			s = s[:len(s)-1]

			return s
		}
	}

	return s
}

// unused (interface satisfier).
func (o *object) Divide(_ int) {}

// Bounds represents a bounding box for an object.
type Bounds struct {
	Minimum *base.Tuple
	Maximum *base.Tuple
}

// DeepCopy performs a deep copy of the object to a new object.
func (b *Bounds) DeepCopy() *Bounds {
	min := *b.Minimum
	max := *b.Maximum

	return &Bounds{
		Minimum: &min,
		Maximum: &max,
	}
}

// use cube intersect method.
func (b *Bounds) intersects(r *ray.Ray) bool {
	// find largest minimum t value and smallest maximum t value for each axis
	// (t is intersection point)
	xtMin, xtMax := checkAxis(r.Origin.GetX(), r.Direction.GetX(), b.Minimum.GetX(), b.Maximum.GetX())
	ytMin, ytMax := checkAxis(r.Origin.GetY(), r.Direction.GetY(), b.Minimum.GetY(), b.Maximum.GetY())
	if xtMin > ytMax || ytMin > xtMax {
		return false
	}
	ztMin, ztMax := checkAxis(r.Origin.GetZ(), r.Direction.GetZ(), b.Minimum.GetZ(), b.Maximum.GetZ())

	tMin := utils.Max(xtMin, ytMin, ztMin)
	tMax := utils.Min(xtMax, ytMax, ztMax)

	return tMin <= tMax
}

// returns two non-overlapping bounding boxes.
func (b *Bounds) split() (*Bounds, *Bounds) {
	// get the box's largest dimension
	dx := b.Maximum.GetX() - b.Minimum.GetX()
	dy := b.Maximum.GetY() - b.Minimum.GetY()
	dz := b.Maximum.GetZ() - b.Minimum.GetZ()

	greatest := utils.Max(dx, dy, dz)

	// variables to help construct the points on the dividing plane
	x0, y0, z0 := b.Minimum.GetX(), b.Minimum.GetY(), b.Minimum.GetZ()
	x1, y1, z1 := b.Maximum.GetX(), b.Maximum.GetY(), b.Maximum.GetZ()

	// adjust the points so that they lie on the dividing plane
	switch greatest {
	case dx:
		x1 = x0 + dx/2.0
		x0 = x1
	case dy:
		y1 = y0 + dy/2.0
		y0 = y1
	case dz:
		z1 = z0 + dz/2.0
		z0 = z1
	}
	midMin := base.NewPoint(x0, y0, z0)
	midMax := base.NewPoint(x1, y1, z1)

	// construct and return the two halves of the bounding box
	left := &Bounds{Minimum: b.Minimum, Maximum: midMax}
	right := &Bounds{Minimum: midMin, Maximum: b.Maximum}

	return left, right
}
