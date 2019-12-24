package object

import (
	"math"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
	"github.com/sjberman/golang-ray-tracer/pkg/utils"
)

// Group represents a group of objects
type Group struct {
	*object
	objects []Object // list of objects in the group
}

// NewGroup returns a new Group object
func NewGroup() *Group {
	return &Group{
		object: newObject(),
	}
}

// Bounds returns the bounding box for the group of objects
func (g *Group) Bounds() *bounds {
	minX, minY, minZ := math.Inf(1), math.Inf(1), math.Inf(1)
	maxX, maxY, maxZ := math.Inf(-1), math.Inf(-1), math.Inf(-1)

	for _, o := range g.objects {
		objBounds := o.Bounds()
		// bounding box points for the object (in parent space)
		points := []*base.Tuple{
			o.GetTransform().MultiplyTuple(objBounds.minimum),
			o.GetTransform().MultiplyTuple(
				base.NewPoint(objBounds.minimum.GetX(), objBounds.minimum.GetY(), objBounds.maximum.GetZ())),
			o.GetTransform().MultiplyTuple(
				base.NewPoint(objBounds.minimum.GetX(), objBounds.maximum.GetY(), objBounds.minimum.GetZ())),
			o.GetTransform().MultiplyTuple(
				base.NewPoint(objBounds.minimum.GetX(), objBounds.maximum.GetY(), objBounds.maximum.GetZ())),
			o.GetTransform().MultiplyTuple(
				base.NewPoint(objBounds.maximum.GetX(), objBounds.minimum.GetY(), objBounds.minimum.GetZ())),
			o.GetTransform().MultiplyTuple(
				base.NewPoint(objBounds.maximum.GetX(), objBounds.minimum.GetY(), objBounds.maximum.GetZ())),
			o.GetTransform().MultiplyTuple(
				base.NewPoint(objBounds.maximum.GetX(), objBounds.maximum.GetY(), objBounds.minimum.GetZ())),
			o.GetTransform().MultiplyTuple(objBounds.maximum),
		}
		for _, p := range points {
			minX = utils.Min(minX, p.GetX())
			minY = utils.Min(minY, p.GetY())
			minZ = utils.Min(minZ, p.GetZ())

			maxX = utils.Max(maxX, p.GetX())
			maxY = utils.Max(maxY, p.GetY())
			maxZ = utils.Max(maxZ, p.GetZ())
		}
	}
	return &bounds{minimum: base.NewPoint(minX, minY, minZ), maximum: base.NewPoint(maxX, maxY, maxZ)}
}

// Add adds children to the group's collection, and sets the group as each child's parent
func (g *Group) Add(objs ...Object) {
	g.objects = append(g.objects, objs...)
	for _, o := range objs {
		o.SetParent(g)
	}
}

// calculates where a ray intersects the objects in a group
func (g *Group) Intersect(ray *ray.Ray) []*Intersection {
	r := g.transformRay(ray)
	if !g.Bounds().intersects(r) {
		return []*Intersection{}
	}

	var ints []*Intersection
	for _, o := range g.objects {
		ints = append(ints, o.Intersect(r)...)
	}
	return sortIntersections(ints)
}

// unused (interface satisfier)
func (g *Group) NormalAt(objectPoint *base.Tuple) *base.Tuple { return nil }
