package object

import (
	"math"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

// Group represents a group of objects.
type Group struct {
	*object
	Objects []Object // list of objects in the group
	bounds  *Bounds
}

// NewGroup returns a new Group object.
func NewGroup() *Group {
	return &Group{
		object: newObject(),
	}
}

// DeepCopy performs a deep copy of the object to a new object.
func (g *Group) DeepCopy() Object {
	newObj := NewGroup()
	for _, o := range g.Objects {
		cpy := o.DeepCopy()
		newObj.Objects = append(newObj.Objects, cpy)
	}

	newMaterial := g.Material
	newObj.SetMaterial(&newMaterial)
	newTransform := g.transform
	newObj.SetTransform(&newTransform)

	if g.bounds != nil {
		newObj.bounds = g.bounds.DeepCopy()
	}

	return newObj
}

// Bounds returns the bounding box for the group of objects.
func (g *Group) Bounds() *Bounds {
	return g.bounds
}

// SetMaterial sets the Group's material.
func (g *Group) SetMaterial(material *Material) {
	for _, o := range g.Objects {
		o.SetMaterial(material)
	}
	g.Material = *material
}

// calculates the bounds based on supplied objects (for group and csg).
func calculateBounds(objects []Object) *Bounds {
	minX, minY, minZ := math.Inf(1), math.Inf(1), math.Inf(1)
	maxX, maxY, maxZ := math.Inf(-1), math.Inf(-1), math.Inf(-1)

	for _, o := range objects {
		objBounds := o.Bounds()
		if objBounds == nil {
			continue
		}
		// bounding box points for the object (in parent space)
		points := []*base.Tuple{
			o.GetTransform().MultiplyTuple(objBounds.Minimum),
			o.GetTransform().MultiplyTuple(
				base.NewPoint(objBounds.Minimum.GetX(), objBounds.Minimum.GetY(), objBounds.Maximum.GetZ())),
			o.GetTransform().MultiplyTuple(
				base.NewPoint(objBounds.Minimum.GetX(), objBounds.Maximum.GetY(), objBounds.Minimum.GetZ())),
			o.GetTransform().MultiplyTuple(
				base.NewPoint(objBounds.Minimum.GetX(), objBounds.Maximum.GetY(), objBounds.Maximum.GetZ())),
			o.GetTransform().MultiplyTuple(
				base.NewPoint(objBounds.Maximum.GetX(), objBounds.Minimum.GetY(), objBounds.Minimum.GetZ())),
			o.GetTransform().MultiplyTuple(
				base.NewPoint(objBounds.Maximum.GetX(), objBounds.Minimum.GetY(), objBounds.Maximum.GetZ())),
			o.GetTransform().MultiplyTuple(
				base.NewPoint(objBounds.Maximum.GetX(), objBounds.Maximum.GetY(), objBounds.Minimum.GetZ())),
			o.GetTransform().MultiplyTuple(objBounds.Maximum),
		}
		for _, p := range points {
			minX = math.Min(minX, p.GetX())
			minY = math.Min(minY, p.GetY())
			minZ = math.Min(minZ, p.GetZ())

			maxX = math.Max(maxX, p.GetX())
			maxY = math.Max(maxY, p.GetY())
			maxZ = math.Max(maxZ, p.GetZ())
		}
	}

	return &Bounds{Minimum: base.NewPoint(minX, minY, minZ), Maximum: base.NewPoint(maxX, maxY, maxZ)}
}

// Add adds children to the group's collection, and sets the group as each child's parent.
func (g *Group) Add(objs ...Object) {
	g.Objects = append(g.Objects, objs...)
	for _, o := range objs {
		o.SetParent(g)
	}
	g.bounds = calculateBounds(g.Objects)
	if g.parent != nil {
		if grp, ok := g.parent.(*Group); ok {
			grp.bounds = calculateBounds(grp.Objects)
		}
	}
}

// divide a group into smaller pieces (boxes).
func (g *Group) Divide(threshold int) {
	if threshold <= len(g.Objects) {
		left, right := g.partitionChildren()
		if len(left) > 0 {
			g.makeSubgroup(left)
		}
		if len(right) > 0 {
			g.makeSubgroup(right)
		}
	}
	for _, o := range g.Objects {
		o.Divide(threshold)
	}
}

// returns two lists of the children that fit into the corresponding halves of the
// group's bounding box.
func (g *Group) partitionChildren() ([]Object, []Object) {
	leftObjs, rightObjs := []Object{}, []Object{}
	leftBox, rightBox := g.Bounds().split()
	var toRemove []Object
	for _, o := range g.Objects {
		bounds := o.Bounds()
		if bounds == nil {
			continue
		}
		bMin := o.GetTransform().MultiplyTuple(bounds.Minimum)
		bMax := o.GetTransform().MultiplyTuple(bounds.Maximum)
		if bMin.GreaterThan(bMax) {
			bMin, bMax = bMax, bMin
		}
		// check which box the object fits in, otherwise don't put it in one
		if (bMin.GreaterThan(leftBox.Minimum) || bMin.Equals(leftBox.Minimum)) &&
			(bMax.LessThan(leftBox.Maximum) || bMax.Equals(leftBox.Maximum)) {
			leftObjs = append(leftObjs, o)
			toRemove = append(toRemove, o)
		} else if (bMin.GreaterThan(rightBox.Minimum) || bMin.Equals(rightBox.Minimum)) &&
			(bMax.LessThan(rightBox.Maximum) || bMax.Equals(rightBox.Maximum)) {
			rightObjs = append(rightObjs, o)
			toRemove = append(toRemove, o)
		}
	}
	// removed unboxed objects from the group
	for _, o := range toRemove {
		g.Objects = Remove(g.Objects, o)
	}

	return leftObjs, rightObjs
}

func (g *Group) makeSubgroup(objs []Object) {
	subgroup := NewGroup()
	subgroup.Add(objs...)
	g.Add(subgroup)
}

// calculates where a ray intersects the objects in a group.
func (g *Group) Intersect(ray *ray.Ray) []*Intersection {
	r := g.transformRay(ray)
	if g.Bounds() == nil || !g.Bounds().intersects(r) {
		return []*Intersection{}
	}

	var ints []*Intersection
	for _, o := range g.Objects {
		ints = append(ints, o.Intersect(r)...)
	}

	return sortIntersections(ints)
}

// unused (interface satisfier).
func (*Group) NormalAt(_ *base.Tuple, _ *Intersection) *base.Tuple { return nil }
