package scene

import (
	"math"
	"sync"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
)

// The total number of recursive reflection traces allowed
const remainingReflections = 4

// World represents the collection of all objects in a scene
type World struct {
	light   *PointLight
	objects []Object
}

// NewWorld returns a new World object
func NewWorld(light *PointLight, objects []Object) *World {
	return &World{
		light:   light,
		objects: objects,
	}
}

// ColorAt returns the color of a specific ray intersection in the world
func (w *World) ColorAt(r *Ray, remaining int) *image.Color {
	intersections := w.intersect(r)
	hit := Hit(intersections)
	if hit == nil {
		return image.Black
	}
	hd := prepareComputations(hit, r, intersections)
	return w.shadeHit(hd, remaining)
}

// shadeHit returns the color at the intersection encapsulated by hitData
func (w *World) shadeHit(hd *hitData, remaining int) *image.Color {
	shadowed := w.isShadowed(hd.overPoint)
	surface := Lighting(
		w.light, hd.object, hd.object.GetMaterial(), hd.point, hd.eyev, hd.normalv, shadowed)
	reflected := w.reflectedColor(hd, remaining)
	refracted := w.refractedColor(hd, remaining)

	material := hd.object.GetMaterial()
	if material.reflective > 0 && material.transparency > 0 {
		reflectance := schlick(hd)
		reflect := reflected.Multiply(reflectance)
		refract := refracted.Multiply(1 - reflectance)
		return surface.Add(reflect).Add(refract)
	}
	return surface.Add(reflected).Add(refracted)
}

// isShadowed returns if a point is in a shadow
func (w *World) isShadowed(point *base.Tuple) bool {
	v, _ := w.light.position.Subtract(point)
	distance := v.Magnitude()
	direction := v.Normalize()

	ray := NewRay(point, direction)
	ints := w.intersect(ray)
	hit := Hit(ints)
	if hit != nil && hit.GetValue() < distance {
		return true
	}
	return false
}

// intersect returns all the intersections between a ray and the objects in the world
func (w *World) intersect(r *Ray) []*Intersection {
	ints := make([]*Intersection, 0, 2*len(w.objects))
	for _, o := range w.objects {
		ints = append(ints, o.intersect(r)...)
	}
	return sortIntersections(ints)
}

// reflectedColor returns the color from a reflected ray
func (w *World) reflectedColor(hd *hitData, remaining int) *image.Color {
	if remaining < 1 || hd.object.GetMaterial().reflective == 0 {
		return image.Black
	}
	remaining--
	reflectRay := NewRay(hd.overPoint, hd.reflectv)
	color := w.ColorAt(reflectRay, remaining)
	return color.Multiply(hd.object.GetMaterial().reflective)
}

// refractedColor returns the color from a refracted ray
func (w *World) refractedColor(hd *hitData, remaining int) *image.Color {
	if remaining < 1 || hd.object.GetMaterial().transparency == 0 {
		return image.Black
	}
	remaining--
	// find the ratio of first index of refraction to the second (inversion of Snell's Law)
	nRatio := hd.n1 / hd.n2
	// cos(theta_i) is the same as the dot product of the two vectors
	cosI := hd.eyev.DotProduct(hd.normalv)
	// find sin(theta_t)^2 via trig identity
	sin2t := math.Pow(nRatio, 2) * (1 - math.Pow(cosI, 2))
	if sin2t > 1 {
		// total internal reflection
		return image.Black
	}
	// find cos(theta_t) via trig identity
	cosT := math.Sqrt(1 - sin2t)
	// compute direction of refracted ray
	direction, _ := hd.normalv.Multiply((nRatio*cosI - cosT)).Subtract(hd.eyev.Multiply(nRatio))
	refractRay := NewRay(hd.underPoint, direction)
	color := w.ColorAt(refractRay, remaining)
	return color.Multiply(hd.object.GetMaterial().transparency)
}

// hitData contains information about a hit intersection
type hitData struct {
	value      float64
	object     Object
	point      *base.Tuple
	overPoint  *base.Tuple
	underPoint *base.Tuple
	eyev       *base.Tuple
	normalv    *base.Tuple
	reflectv   *base.Tuple
	n1, n2     float64 // refractive index for source/dest of ray
	inside     bool
}

// Uses an intersection and ray to build up the hit data
func prepareComputations(
	intersection *Intersection,
	ray *Ray,
	allIntersections []*Intersection,
) *hitData {
	hd := &hitData{
		value:  intersection.GetValue(),
		object: intersection.GetObject(),
		eyev:   ray.GetDirection().Negate(),
	}
	hd.point = ray.Position(hd.value)
	hd.normalv = hd.object.normalAt(hd.point)

	if hd.normalv.DotProduct(hd.eyev) < 0 {
		// Hit occurs inside the object (normal points away from eye)
		hd.inside = true
		hd.normalv = hd.normalv.Negate()
	}
	hd.reflectv = ray.GetDirection().Reflect(hd.normalv)
	// have a point just above normal point to account for accidental shadow calculation when
	// a ray hits the object it's leaving
	hd.overPoint, _ = hd.point.Add(hd.normalv.Multiply(base.Epsilon * 2))

	// have a point just below normal point for the origination of refracted rays
	hd.underPoint, _ = hd.point.Subtract(hd.normalv.Multiply(base.Epsilon * 2))

	// containers is a list of objects that we've intersected, but haven't yet exited
	containers := []Object{}
	for _, iSection := range allIntersections {
		// if intersection is the hit, n1 is the refractive index of the last object in containers list
		if iSection == intersection {
			if len(containers) == 0 {
				hd.n1 = 1
			} else {
				hd.n1 = containers[len(containers)-1].GetMaterial().refractiveIndex
			}
		}

		// if intersection's object is already in containers list, then this intersection must
		// be exiting the object; otherwise, the intersection is entering the object
		intObject := iSection.GetObject()
		if contains(containers, intObject) {
			containers = remove(containers, intObject)
		} else {
			containers = append(containers, intObject)
		}

		// if intersection is the hit, n2 is the refractive index of the last object in containers list
		if iSection == intersection {
			if len(containers) == 0 {
				hd.n2 = 1
			} else {
				hd.n2 = containers[len(containers)-1].GetMaterial().refractiveIndex
			}
			break
		}
	}
	return hd
}

// returns the reflectance, which represents what fraction of light is reflected
// (named for Christophe Schlick)
func schlick(hd *hitData) float64 {
	cos := hd.eyev.DotProduct(hd.normalv)

	// total internal reflection only occurs if n1 > n2
	if hd.n1 > hd.n2 {
		n := hd.n1 / hd.n2
		sin2t := math.Pow(n, 2) * (1 - math.Pow(cos, 2))
		if sin2t > 1 {
			return 1
		}
		cos = math.Sqrt(1 - sin2t)
	}
	r0 := math.Pow((hd.n1-hd.n2)/(hd.n1+hd.n2), 2)
	return r0 + (1-r0)*math.Pow(1-cos, 5)
}

func Render(c *Camera, w *World) *image.Canvas {
	canvas := image.NewCanvas(c.hsize, c.vsize)

	for y := 0; y < c.vsize-1; y++ {
		var wg sync.WaitGroup
		wg.Add(c.hsize - 1)
		for x := 0; x < c.hsize-1; x++ {
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

// returns if an Object slice contains an object
func contains(s []Object, o Object) bool {
	for _, e := range s {
		if o == e {
			return true
		}
	}
	return false
}

// removes an Object from a slice of Objects
func remove(s []Object, o Object) []Object {
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
