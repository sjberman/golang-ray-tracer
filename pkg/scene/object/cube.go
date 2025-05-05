package object

import (
	"math"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	grtMath "github.com/sjberman/golang-ray-tracer/pkg/math"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

// Cube is a cube object.
type Cube struct {
	*object
}

// NewCube returns a new Cube object.
func NewCube() *Cube {
	return &Cube{
		object: newObject(),
	}
}

// DeepCopy performs a deep copy of the object to a new object.
func (c *Cube) DeepCopy() Object {
	newObj := NewCube()
	newMaterial := c.Material
	newObj.SetMaterial(&newMaterial)
	newTransform := c.transform
	newObj.SetTransform(&newTransform)

	return newObj
}

// Bounds returns the untransformed bounds of a cube.
func (c *Cube) Bounds() *Bounds {
	return &Bounds{
		Minimum: base.NewPoint(-1, -1, -1),
		Maximum: base.NewPoint(1, 1, 1),
	}
}

// calculates where a ray intersects a cube.
func (c *Cube) Intersect(ray *ray.Ray) []*Intersection {
	r := c.transformRay(ray)
	// find largest minimum t value and smallest maximum t value for each axis
	// (t is intersection point)
	xtMin, xtMax := checkAxis(r.Origin.GetX(), r.Direction.GetX(), -1, 1)
	ytMin, ytMax := checkAxis(r.Origin.GetY(), r.Direction.GetY(), -1, 1)
	if xtMin > ytMax || ytMin > xtMax {
		return []*Intersection{}
	}
	ztMin, ztMax := checkAxis(r.Origin.GetZ(), r.Direction.GetZ(), -1, 1)

	tMin := grtMath.Max(xtMin, ytMin, ztMin)
	tMax := grtMath.Min(xtMax, ytMax, ztMax)

	if tMin > tMax {
		return []*Intersection{}
	}

	return []*Intersection{NewIntersection(tMin, c), NewIntersection(tMax, c)}
}

// wrapper for the normalAt interface function, using the common normal function
// with the specific cube logic embedded.
func (c *Cube) NormalAt(objectPoint *base.Tuple, hit *Intersection) *base.Tuple {
	return commonNormalAt(c, objectPoint, hit, cubeNormal)
}

// cube-specific calculation of the normal.
func cubeNormal(objectPoint *base.Tuple, _ Object, _ *Intersection) *base.Tuple {
	absX := math.Abs(objectPoint.GetX())
	absY := math.Abs(objectPoint.GetY())
	absZ := math.Abs(objectPoint.GetZ())
	maxC := grtMath.Max(absX, absY, absZ)
	switch maxC {
	case absX:
		return base.NewVector(objectPoint.GetX(), 0, 0)
	case absY:
		return base.NewVector(0, objectPoint.GetY(), 0)
	}

	return base.NewVector(0, 0, objectPoint.GetZ())
}

// checkAxis finds the min and max intersection values for the axis.
func checkAxis(origin, direction, min, max float64) (float64, float64) {
	var tMin, tMax float64
	tMinNumerator := min - origin
	tMaxNumerator := max - origin

	if math.Abs(direction) >= base.Epsilon {
		tMin = tMinNumerator / direction
		tMax = tMaxNumerator / direction
	} else {
		// if denominator is effectively zero, multiply by infinity to ensure
		// values have the correct sign (positive or negative)
		tMin = tMinNumerator * math.Inf(1)
		tMax = tMaxNumerator * math.Inf(1)
	}

	if tMin > tMax {
		tMin, tMax = tMax, tMin
	}

	return tMin, tMax
}
