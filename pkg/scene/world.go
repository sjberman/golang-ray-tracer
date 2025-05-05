package scene

import (
	"math"
	"slices"
	"sync"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/object"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

// The total number of recursive reflection traces allowed.
const remainingReflections = 4

// World represents the collection of all objects in a scene.
type World struct {
	lights  []*PointLight
	objects []object.Object
}

// NewWorld returns a new World object.
func NewWorld(lights []*PointLight, objects []object.Object) *World {
	return &World{
		lights:  lights,
		objects: objects,
	}
}

// ColorAt returns the color of a specific ray intersection in the world.
func (w *World) ColorAt(r *ray.Ray, remaining int) *image.Color {
	intersections := w.intersect(r)
	hit := object.Hit(intersections)
	if hit == nil {
		return image.Black
	}
	hd := prepareComputations(hit, r, intersections)

	return w.shadeHit(hd, remaining)
}

// shadeHit returns the color at the intersection encapsulated by hitData.
func (w *World) shadeHit(hd *hitData, remaining int) *image.Color {
	surface := image.Black
	for _, light := range w.lights {
		shadowed := w.isShadowed(light, hd.overPoint)
		surface = surface.Add(lighting(
			light, hd.object, hd.object.GetMaterial(), hd.point, hd.eyev, hd.normalv, shadowed))
	}
	reflected := w.reflectedColor(hd, remaining)
	refracted := w.refractedColor(hd, remaining)

	material := hd.object.GetMaterial()
	if material.Reflective > 0 && material.Transparency > 0 {
		reflectance := schlick(hd)
		reflect := reflected.Multiply(reflectance)
		refract := refracted.Multiply(1 - reflectance)

		return surface.Add(reflect).Add(refract)
	}

	return surface.Add(reflected).Add(refracted)
}

// isShadowed returns if a point is in a shadow.
func (w *World) isShadowed(light *PointLight, point *base.Tuple) bool {
	v := light.position.Subtract(point)
	distance := v.Magnitude()
	direction := v.Normalize()

	ray := ray.NewRay(point, direction)
	ints := w.intersect(ray)
	hit := object.Hit(ints)
	if hit != nil && hit.Value < distance && hit.Object.GetMaterial().Shadow {
		return true
	}

	return false
}

// intersect returns all the intersections between a ray and the objects in the world.
func (w *World) intersect(r *ray.Ray) []*object.Intersection {
	ints := make([]*object.Intersection, 0, 2*len(w.objects))
	for _, o := range w.objects {
		ints = append(ints, o.Intersect(r)...)
	}

	return object.Intersections(ints...)
}

// reflectedColor returns the color from a reflected ray.
func (w *World) reflectedColor(hd *hitData, remaining int) *image.Color {
	if remaining < 1 || hd.object.GetMaterial().Reflective == 0 {
		return image.Black
	}

	remaining--
	reflectRay := ray.NewRay(hd.overPoint, hd.reflectv)
	color := w.ColorAt(reflectRay, remaining)

	return color.Multiply(hd.object.GetMaterial().Reflective)
}

// refractedColor returns the color from a refracted ray.
func (w *World) refractedColor(hd *hitData, remaining int) *image.Color {
	if remaining < 1 || hd.object.GetMaterial().Transparency == 0 {
		return image.Black
	}

	remaining--
	// find the ratio of first index of refraction to the second (inversion of Snell's Law)
	nRatio := hd.n1 / hd.n2
	// cos(theta_i) is the same as the dot product of the two vectors
	cosI := hd.eyev.DotProduct(hd.normalv)
	// find sin(theta_t)^2 via trig identity
	sin2t := nRatio * nRatio * (1 - cosI*cosI)
	if sin2t > 1 {
		// total internal reflection
		return image.Black
	}
	// find cos(theta_t) via trig identity
	cosT := math.Sqrt(1 - sin2t)
	// compute direction of refracted ray
	direction := hd.normalv.Multiply((nRatio*cosI - cosT)).Subtract(hd.eyev.Multiply(nRatio))
	refractRay := ray.NewRay(hd.underPoint, direction)
	color := w.ColorAt(refractRay, remaining)

	return color.Multiply(hd.object.GetMaterial().Transparency)
}

// hitData contains information about a hit intersection.
type hitData struct {
	value      float64
	object     object.Object
	point      *base.Tuple
	overPoint  *base.Tuple
	underPoint *base.Tuple
	eyev       *base.Tuple
	normalv    *base.Tuple
	reflectv   *base.Tuple
	n1, n2     float64 // refractive index for source/dest of ray
	inside     bool
}

// Uses an intersection and ray to build up the hit data.
func prepareComputations(
	intersection *object.Intersection,
	ray *ray.Ray,
	allIntersections []*object.Intersection,
) *hitData {
	hd := &hitData{
		value:  intersection.Value,
		object: intersection.Object,
		eyev:   ray.Direction.Negate(),
	}
	hd.point = ray.Position(hd.value)
	hd.normalv = hd.object.NormalAt(hd.point, intersection)

	if hd.normalv.DotProduct(hd.eyev) < 0 {
		// Hit occurs inside the object (normal points away from eye)
		hd.inside = true
		hd.normalv = hd.normalv.Negate()
	}
	hd.reflectv = ray.Direction.Reflect(hd.normalv)
	// have a point just above normal point to account for accidental shadow calculation when
	// a ray hits the object it's leaving
	hd.overPoint = hd.point.Add(hd.normalv.Multiply(base.Epsilon * 2))

	// have a point just below normal point for the origination of refracted rays
	hd.underPoint = hd.point.Subtract(hd.normalv.Multiply(base.Epsilon * 2))

	// containers is a list of objects that we've intersected, but haven't yet exited
	containers := []object.Object{}
	for _, iSection := range allIntersections {
		// if intersection is the hit, n1 is the refractive index of the last object in containers list
		if iSection == intersection {
			if len(containers) == 0 {
				hd.n1 = 1
			} else {
				hd.n1 = containers[len(containers)-1].GetMaterial().RefractiveIndex
			}
		}

		// if intersection's object is already in containers list, then this intersection must
		// be exiting the object; otherwise, the intersection is entering the object
		if slices.Contains(containers, iSection.Object) {
			containers = object.Remove(containers, iSection.Object)
		} else {
			containers = append(containers, iSection.Object)
		}

		// if intersection is the hit, n2 is the refractive index of the last object in containers list
		if iSection == intersection {
			if len(containers) == 0 {
				hd.n2 = 1
			} else {
				hd.n2 = containers[len(containers)-1].GetMaterial().RefractiveIndex
			}

			break
		}
	}

	return hd
}

// returns the reflectance, which represents what fraction of light is reflected
// (named for Christophe Schlick).
func schlick(hd *hitData) float64 {
	cos := hd.eyev.DotProduct(hd.normalv)

	// total internal reflection only occurs if n1 > n2
	if hd.n1 > hd.n2 {
		n := hd.n1 / hd.n2
		sin2t := n * n * (1 - cos*cos)
		if sin2t > 1 {
			return 1
		}
		cos = math.Sqrt(1 - sin2t)
	}
	r0 := (hd.n1 - hd.n2) / (hd.n1 + hd.n2) * ((hd.n1 - hd.n2) / (hd.n1 + hd.n2))

	return r0 + (1-r0)*math.Pow(1-cos, 5)
}

func Render(c *Camera, w *World) *image.Canvas {
	canvas := image.NewCanvas(c.hsize, c.vsize)

	for y := range c.vsize - 1 {
		var wg sync.WaitGroup
		wg.Add(c.hsize - 1)

		for x := range c.hsize - 1 {
			go func(x, y int) {
				defer wg.Done()

				ray := c.RayForPixel(x, y)
				color := w.ColorAt(ray, remainingReflections)
				canvas.WritePixel(x, y, color)
			}(x, y)
		}
		wg.Wait()
	}

	return canvas
}
