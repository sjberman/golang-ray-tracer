package object

import (
	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

const (
	union        = "union"
	intersection = "intersection"
	difference   = "difference"
)

// Csg is a Constructive Solid Geometry object.
type Csg struct {
	*object
	operation   string
	left, right Object
	bounds      *Bounds
}

// NewCsg returns a new Csg object.
func NewCsg(op string, l, r Object) *Csg {
	csg := &Csg{
		object:    newObject(),
		operation: op,
		left:      l,
		right:     r,
	}
	l.SetParent(csg)
	r.SetParent(csg)
	csg.bounds = calculateBounds([]Object{csg.left, csg.right})

	return csg
}

// DeepCopy performs a deep copy of the object to a new object.
func (csg *Csg) DeepCopy() Object {
	left := csg.left.DeepCopy()
	right := csg.right.DeepCopy()
	newObj := NewCsg(csg.operation, left, right)

	newMaterial := csg.Material
	newObj.SetMaterial(&newMaterial)
	newTransform := csg.transform
	newObj.SetTransform(&newTransform)

	if csg.bounds != nil {
		newObj.bounds = csg.bounds.DeepCopy()
	}

	return newObj
}

// Bounds returns the bounding box for the csg of objects.
func (csg *Csg) Bounds() *Bounds {
	return csg.bounds
}

// determines if an intersection is allowed based on the operation.
func intersectionAllowed(op string, lhit, inl, inr bool) bool {
	switch op {
	case union:
		return (lhit && !inr) || (!lhit && !inl)
	case intersection:
		return (lhit && inr) || (!lhit && inl)
	case difference:
		return (lhit && !inr) || (!lhit && inl)
	}

	return false
}

// checks if each intersection is allowed and filters.
func (csg *Csg) filterIntersections(ints []*Intersection) []*Intersection {
	var inl, inr bool
	result := []*Intersection{}

	for _, intersection := range ints {
		lhit := includes(csg.left, intersection.Object)
		if intersectionAllowed(csg.operation, lhit, inl, inr) {
			result = append(result, intersection)
		}

		// depending on which object was hit, toggle either inl or inr
		if lhit {
			inl = !inl
		} else {
			inr = !inr
		}
	}

	return result
}

// calculates where a ray intersects objects in a csg.
func (csg *Csg) Intersect(ray *ray.Ray) []*Intersection {
	r := csg.transformRay(ray)
	ints := csg.left.Intersect(r)
	ints = append(ints, csg.right.Intersect(r)...)

	return csg.filterIntersections(sortIntersections(ints))
}

// divide a csg into smaller pieces (boxes).
func (csg *Csg) Divide(threshold int) {
	csg.left.Divide(threshold)
	csg.right.Divide(threshold)
}

// unused (interface satisfier).
func (*Csg) NormalAt(_ *base.Tuple, _ *Intersection) *base.Tuple { return nil }
