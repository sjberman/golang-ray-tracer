package object

import (
	"math"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
	"github.com/sjberman/golang-ray-tracer/pkg/utils"
)

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
func (c *Cube) Intersect(ray *ray.Ray) []*Intersection {
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
func (c *Cube) NormalAt(objectPoint *base.Tuple) *base.Tuple {
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
