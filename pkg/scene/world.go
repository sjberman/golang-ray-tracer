package scene

import (
	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
)

// World represents the collection of all objects in a scene
type World struct {
	light   *PointLight
	objects []*Sphere
}

// NewWorld returns a new World object
func NewWorld(light *PointLight, objects []*Sphere) *World {
	return &World{
		light:   light,
		objects: objects,
	}
}

// ColorAt returns the color of a specific ray intersection in the world
func (w *World) ColorAt(r *Ray) *image.Color {
	intersections := w.intersect(r)
	hit := Hit(intersections)
	if hit == nil {
		return &image.Black
	}
	hd := prepareComputations(hit, r)
	return Lighting(w.light, hd.object.GetMaterial(), hd.point, hd.eyev, hd.normalv)
}

// intersect returns all the intersections between a ray and the objects in the world
func (w *World) intersect(r *Ray) []*Intersection {
	ints := make([]*Intersection, 0, 2*len(w.objects))
	for _, o := range w.objects {
		ints = append(ints, r.Intersect(o)...)
	}
	return sortIntersections(ints)
}

// hitData contains information about a hit intersection
type hitData struct {
	value   float64
	object  *Sphere
	point   *base.Tuple
	eyev    *base.Tuple
	normalv *base.Tuple
	inside  bool
}

// Uses an intersection and ray to build up the hit data
func prepareComputations(intersection *Intersection, ray *Ray) *hitData {
	hd := &hitData{
		value:  intersection.GetValue(),
		object: intersection.GetObject(),
		eyev:   ray.GetDirection().Negate(),
	}
	hd.point = ray.Position(hd.value)
	hd.normalv = hd.object.NormalAt(hd.point)

	if hd.normalv.DotProduct(hd.eyev) < 0 {
		// Hit occurs inside the shape (normal points away from eye)
		hd.inside = true
		hd.normalv = hd.normalv.Negate()
	}
	return hd
}

func Render(c *Camera, w *World) *image.Canvas {
	canvas := image.NewCanvas(c.hsize, c.vsize)

	for y := 0; y < c.vsize-1; y++ {
		for x := 0; x < c.hsize-1; x++ {
			ray := c.RayForPixel(x, y)
			color := w.ColorAt(ray)
			canvas.WritePixel(x, y, color)
		}
	}
	return canvas
}

// // ShadeHit returns the color at the intersection encapsulated by hitData
// func (w *World) ShadeHit(hd *hitData) *image.Color {
// 	return Lighting(w.light, hd.object.GetMaterial(), hd.point, hd.eyev, hd.normalv)
// }
