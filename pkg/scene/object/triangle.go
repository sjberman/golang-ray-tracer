package object

import (
	"math"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
	"github.com/sjberman/golang-ray-tracer/pkg/utils"
)

// Triangle is a triangle object
type Triangle struct {
	*object
	P1, P2, P3   *base.Tuple
	edge1, edge2 *base.Tuple
	normal       *base.Tuple
}

// SmoothTriangle is a smooth triangle object
type SmoothTriangle struct {
	*Triangle
	N1, N2, N3 *base.Tuple
}

// NewTriangle returns a new Triangle object
func NewTriangle(p1, p2, p3 *base.Tuple) *Triangle {
	e1 := p2.Subtract(p1)
	e2 := p3.Subtract(p1)
	normal := e2.CrossProduct(e1).Normalize()

	return &Triangle{
		object: newObject(),
		P1:     p1,
		P2:     p2,
		P3:     p3,
		edge1:  e1,
		edge2:  e2,
		normal: normal,
	}
}

// NewSmoothTriangle returns a new SmoothTriangle object
func NewSmoothTriangle(p1, p2, p3, n1, n2, n3 *base.Tuple) *SmoothTriangle {
	return &SmoothTriangle{
		Triangle: NewTriangle(p1, p2, p3),
		N1:       n1,
		N2:       n2,
		N3:       n3,
	}
}

// DeepCopy performs a deep copy of the object to a new object
func (t *Triangle) DeepCopy() Object {
	newObj := NewTriangle(t.P1, t.P2, t.P3)
	newMaterial := t.Material
	newObj.SetMaterial(&newMaterial)
	newTransform := t.transform
	newObj.SetTransform(&newTransform)
	return newObj
}

// DeepCopy performs a deep copy of the object to a new object
func (t *SmoothTriangle) DeepCopy() Object {
	newObj := NewSmoothTriangle(t.P1, t.P2, t.P3, t.N1, t.N2, t.N3)
	newMaterial := t.Material
	newObj.SetMaterial(&newMaterial)
	newTransform := t.transform
	newObj.SetTransform(&newTransform)
	return newObj
}

// Bounds returns the untransformed bounds of a triangle
func (t *Triangle) Bounds() *Bounds {
	xMin := utils.Min(t.P1.GetX(), t.P2.GetX(), t.P3.GetX())
	yMin := utils.Min(t.P1.GetY(), t.P2.GetY(), t.P3.GetY())
	zMin := utils.Min(t.P1.GetZ(), t.P2.GetZ(), t.P3.GetZ())

	xMax := utils.Max(t.P1.GetX(), t.P2.GetX(), t.P3.GetX())
	yMax := utils.Max(t.P1.GetY(), t.P2.GetY(), t.P3.GetY())
	zMax := utils.Max(t.P1.GetZ(), t.P2.GetZ(), t.P3.GetZ())

	return &Bounds{
		Minimum: base.NewPoint(xMin, yMin, zMin),
		Maximum: base.NewPoint(xMax, yMax, zMax),
	}
}

// calculates where a ray intersects a triangle
func (t *Triangle) Intersect(ray *ray.Ray) []*Intersection {
	r := t.transformRay(ray)
	// Muller-Trumbore algorithm
	dirCrossE2 := r.Direction.CrossProduct(t.edge2)
	determinant := t.edge1.DotProduct(dirCrossE2)
	if math.Abs(determinant) < base.Epsilon {
		// ray is parallel
		return []*Intersection{}
	}

	f := 1.0 / determinant
	p1ToOrigin := r.Origin.Subtract(t.P1)
	u := f * p1ToOrigin.DotProduct(dirCrossE2)
	if u < 0 || u > 1 {
		// ray misses P1-P3 edge
		return []*Intersection{}
	}
	originCrossE1 := p1ToOrigin.CrossProduct(t.edge1)
	v := f * ray.Direction.DotProduct(originCrossE1)
	if v < 0 || (u+v) > 1 {
		// ray misses P1-P2 or P2-P3 edges
		return []*Intersection{}
	}

	val := f * t.edge2.DotProduct(originCrossE1)
	intersection := NewIntersection(val, t)
	intersection.u = u
	intersection.v = v
	return []*Intersection{intersection}
}

// wrapper for the normalAt interface function, using the common normal function
// with the specific triangle logic embedded
func (t *Triangle) NormalAt(objectPoint *base.Tuple, hit *Intersection) *base.Tuple {
	return commonNormalAt(t, objectPoint, hit, triangleNormal)
}

// triangle-specific calculation of the normal
func triangleNormal(objectPoint *base.Tuple, o Object, _ *Intersection) *base.Tuple {
	return o.(*Triangle).normal
}

// calculates where a ray intersects a smooth triangle
func (t *SmoothTriangle) Intersect(ray *ray.Ray) []*Intersection {
	intersection := t.Triangle.Intersect(ray)
	if len(intersection) > 0 {
		// replace Triangle object with SmoothTriangle object
		intersection[0].Object = t
	}
	return intersection
}

// wrapper for the normalAt interface function, using the common normal function
// with the specific smooth triangle logic embedded
func (t *SmoothTriangle) NormalAt(objectPoint *base.Tuple, hit *Intersection) *base.Tuple {
	return commonNormalAt(t, objectPoint, hit, smoothTriangleNormal)
}

// smoothTriangle-specific calculation of the normal
func smoothTriangleNormal(objectPoint *base.Tuple, o Object, hit *Intersection) *base.Tuple {
	t := o.(*SmoothTriangle)
	return t.N2.Multiply(hit.u).Add(t.N3.Multiply(hit.v).Add(t.N1.Multiply((1 - hit.u - hit.v))))
}
