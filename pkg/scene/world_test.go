package scene

import (
	"math"
	"sync"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/object"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

var (
	testLights  = []*PointLight{NewPointLight(base.NewPoint(-10, 10, -10), image.White)}
	testObjects []object.Object
	setupOnce   sync.Once
)

func worldTestSetup() {
	setupOnce.Do(func() {
		s1 := object.NewSphere()
		s1.GetMaterial().Color = image.NewColor(0.8, 1.0, 0.6)
		s1.GetMaterial().Diffuse = 0.7
		s1.GetMaterial().Specular = 0.2
		s2 := object.NewSphere()
		s2.SetTransform(base.Scale(0.5, 0.5, 0.5))
		testObjects = []object.Object{s1, s2}
	})
}

func TestNewWorld(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)
	worldTestSetup()

	w := NewWorld(testLights, testObjects)
	g.Expect(w.lights).To(Equal(testLights))
	g.Expect(w.objects).To(Equal(testObjects))
}

func TestColorAt(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)
	worldTestSetup()

	// when ray misses
	w := NewWorld(testLights, testObjects)
	r := ray.NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 1, 0))
	color := w.ColorAt(r, remainingReflections)
	g.Expect(color).To(Equal(image.Black))

	// when ray hits
	r = ray.NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
	color = w.ColorAt(r, remainingReflections)
	g.Expect(color).To(Equal(image.NewColor(0.38066119308103435, 0.47582649135129296, 0.28549589481077575)))

	// intersection behind the ray
	outer := w.objects[0].DeepCopy()
	outer.GetMaterial().Ambient = 1
	inner := w.objects[1].DeepCopy()
	inner.GetMaterial().Ambient = 1
	w = NewWorld(testLights, []object.Object{outer, inner})

	r = ray.NewRay(base.NewPoint(0, 0, 0.75), base.NewVector(0, 0, -1))
	color = w.ColorAt(r, remainingReflections)
	g.Expect(color).To(Equal(inner.GetMaterial().Color))

	// with mutally reflective surfaces (assume no infinite recursion)
	lower := object.NewPlane()
	lower.Reflective = 1
	lower.SetTransform(base.Translate(0, -1, 0))

	upper := object.NewPlane()
	upper.Reflective = 1
	upper.SetTransform(base.Translate(0, 1, 0))

	w.objects = []object.Object{lower, upper}
	r = ray.NewRay(base.NewPoint(0, 0, 0), base.NewVector(0, 1, 0))
	g.Eventually(func() *image.Color {
		return w.ColorAt(r, remainingReflections)
	}).ShouldNot(BeNil())
}

func TestShadeHit(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)
	worldTestSetup()

	w := NewWorld(testLights, testObjects)
	floor := object.NewPlane()
	floor.SetTransform(base.Translate(0, -1, 0))
	floor.Transparency = 0.5
	floor.RefractiveIndex = 1.5
	ball := object.NewSphere()
	ball.Color = image.NewColor(1, 0, 0)
	ball.Ambient = 0.5
	ball.SetTransform(base.Translate(0, -3.5, -0.5))
	w.objects = append(w.objects, floor, ball)

	ray := ray.NewRay(base.NewPoint(0, 0, -3), base.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	ints := object.Intersections(object.NewIntersection(math.Sqrt(2), floor))
	hd := prepareComputations(ints[0], ray, ints)
	color := w.shadeHit(hd, 5)
	g.Expect(color).To(Equal(image.NewColor(0.9364253889815014, 0.6864253889815014, 0.6864253889815014)))
}

func TestIsShadowed(t *testing.T) {
	t.Parallel()
	worldTestSetup()

	tests := []struct {
		name     string
		point    *base.Tuple
		inShadow bool
	}{
		{
			name:     "nothing collinear with point and light",
			point:    base.NewPoint(0, 10, 0),
			inShadow: false,
		},
		{
			name:     "object between point and light",
			point:    base.NewPoint(10, -10, 10),
			inShadow: true,
		},
		{
			name:     "object is behind the light",
			point:    base.NewPoint(-20, 20, -20),
			inShadow: false,
		},
		{
			name:     "object is behind the point",
			point:    base.NewPoint(-2, 2, -2),
			inShadow: false,
		},
	}

	w := NewWorld(testLights, testObjects)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			g.Expect(w.isShadowed(w.lights[0], test.point)).To(Equal(test.inShadow))
		})
	}
}

func TestIntersect(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)
	worldTestSetup()

	w := NewWorld(testLights, testObjects)
	r := ray.NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
	ints := w.intersect(r)

	g.Expect(len(ints)).To(Equal(4))
	g.Expect(ints[0].Value).To(Equal(4.0))
	g.Expect(ints[1].Value).To(Equal(4.5))
	g.Expect(ints[2].Value).To(Equal(5.5))
	g.Expect(ints[3].Value).To(Equal(6.0))
}

func TestReflectedColor(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)
	worldTestSetup()

	// nonreflective material
	w := NewWorld(testLights, testObjects)
	r := ray.NewRay(base.NewPoint(0, 0, 0), base.NewVector(0, 0, 1))
	s := w.objects[1].DeepCopy()
	s.GetMaterial().Ambient = 1
	intersection := object.NewIntersection(1, s)
	hd := prepareComputations(intersection, r, object.Intersections(intersection))
	color := w.reflectedColor(hd, 1)
	g.Expect(color).To(Equal(image.Black))

	// reflective material with max recursive depth
	p := object.NewPlane()
	p.Reflective = 0.5
	p.SetTransform(base.Translate(0, -1, 0))
	w.objects = append(w.objects, p)
	r = ray.NewRay(base.NewPoint(0, 0, -3), base.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	intersection = object.NewIntersection(math.Sqrt(2), p)
	hd = prepareComputations(intersection, r, object.Intersections(intersection))
	color = w.reflectedColor(hd, 0)
	g.Expect(color).To(Equal(image.Black))

	// reflective material
	color = w.reflectedColor(hd, 1)
	g.Expect(color).To(Equal(image.NewColor(0.19033404421300906, 0.23791755526626135, 0.14275053315975678)))

	// shadeHit with reflective material
	color = w.shadeHit(hd, remainingReflections)
	g.Expect(color).To(Equal(image.NewColor(0.8767594331945104, 0.9243429442477628, 0.8291759221412582)))
}

func TestRefractedColor(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)
	worldTestSetup()

	// opaque
	w := NewWorld(testLights, testObjects)
	r := ray.NewRay(base.NewPoint(0, 0, 0), base.NewVector(0, 0, 1))
	s := w.objects[0].DeepCopy()
	ints := object.Intersections(object.NewIntersection(4, s), object.NewIntersection(6, s))
	hd := prepareComputations(ints[0], r, ints)
	color := w.refractedColor(hd, 5)
	g.Expect(color).To(Equal(image.Black))

	// max recursive depth
	s.GetMaterial().Transparency = 1
	s.GetMaterial().RefractiveIndex = 1.5
	r = ray.NewRay(base.NewPoint(0, 0, 0), base.NewVector(0, 0, 1))
	ints = object.Intersections(object.NewIntersection(4, s), object.NewIntersection(6, s))
	hd = prepareComputations(ints[0], r, ints)
	color = w.refractedColor(hd, 0)
	g.Expect(color).To(Equal(image.Black))

	// total internal reflection
	r = ray.NewRay(base.NewPoint(0, 0, math.Sqrt(2)/2), base.NewVector(0, 1, 0))
	ints = object.Intersections(object.NewIntersection(-math.Sqrt(2)/2, s), object.NewIntersection(math.Sqrt(2)/2, s))
	// inside of sphere, so look at second intersection
	hd = prepareComputations(ints[1], r, ints)
	color = w.refractedColor(hd, 5)
	g.Expect(color).To(Equal(image.Black))
}

func TestRefractedColor_RefractedRay(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)
	worldTestSetup()

	s1 := testObjects[0].DeepCopy()
	s1.GetMaterial().Ambient = 1
	s1.GetMaterial().Pattern = image.NewMockPattern()
	s2 := testObjects[1].DeepCopy()
	s2.GetMaterial().Transparency = 1
	s2.GetMaterial().RefractiveIndex = 1.5

	w := NewWorld(testLights, []object.Object{s1, s2})

	r := ray.NewRay(base.NewPoint(0, 0, 0.1), base.NewVector(0, 1, 0))
	ints := object.Intersections(
		object.NewIntersection(-0.9899, s1), object.NewIntersection(-0.4899, s2),
		object.NewIntersection(0.4899, s2), object.NewIntersection(0.9899, s1),
	)
	hd := prepareComputations(ints[2], r, ints)
	color := w.refractedColor(hd, 5)
	g.Expect(color).To(Equal(image.NewColor(0, 0.9988845862650526, 0.04721846378372032)))
}

func TestPrepareComputations(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	r := ray.NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
	intersection := object.NewIntersection(4, object.NewSphere())
	hd := prepareComputations(intersection, r, object.Intersections(intersection))
	g.Expect(hd.value).To(Equal(intersection.Value))
	g.Expect(hd.object).To(Equal(intersection.Object))
	g.Expect(hd.point).To(Equal(base.NewPoint(0, 0, -1)))
	g.Expect(hd.eyev).To(Equal(base.NewVector(0, 0, -1)))
	g.Expect(hd.normalv).To(Equal(base.NewVector(0, 0, -1)))
	g.Expect(hd.reflectv).To(Equal(base.NewVector(0, 0, -1)))
	g.Expect(hd.inside).To(BeFalse())

	// hit should offset the point
	s := object.NewSphere()
	s.SetTransform(base.Translate(0, 0, 1))
	intersection = object.NewIntersection(5, s)
	hd = prepareComputations(intersection, r, object.Intersections(intersection))
	g.Expect(hd.overPoint.GetZ()).To(BeNumerically("<", -base.Epsilon/2))
	g.Expect(hd.point.GetZ()).To(BeNumerically(">", hd.overPoint.GetZ()))

	// intersection occurs on inside
	r = ray.NewRay(base.Origin, base.NewVector(0, 0, 1))
	intersection = object.NewIntersection(1, object.NewSphere())
	hd = prepareComputations(intersection, r, object.Intersections(intersection))
	g.Expect(hd.point).To(Equal(base.NewPoint(0, 0, 1)))
	g.Expect(hd.eyev).To(Equal(base.NewVector(0, 0, -1)))
	g.Expect(hd.normalv).To(Equal(base.NewVector(0, 0, -1)))
	g.Expect(hd.inside).To(BeTrue())

	// reflect off plane
	r = ray.NewRay(base.NewPoint(0, 1, -1), base.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	intersection = object.NewIntersection(math.Sqrt(2), object.NewPlane())
	hd = prepareComputations(intersection, r, object.Intersections(intersection))
	g.Expect(hd.reflectv).To(Equal(base.NewVector(0, math.Sqrt(2)/2, math.Sqrt(2)/2)))

	// under point is offset below the surface
	r = ray.NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
	s = object.GlassSphere()
	s.SetTransform(base.Translate(0, 0, 1))
	intersection = object.NewIntersection(5, s)
	hd = prepareComputations(intersection, r, object.Intersections(intersection))
	g.Expect(hd.underPoint.GetZ()).To(BeNumerically(">", base.Epsilon/2))
	g.Expect(hd.point.GetZ()).To(BeNumerically("<", hd.underPoint.GetZ()))
}

func TestPrepareComputations_RefractiveIndicesMultipleIntersections(t *testing.T) {
	t.Parallel()

	s1 := object.GlassSphere()
	s1.SetTransform(base.Scale(2, 2, 2))

	s2 := object.GlassSphere()
	s2.SetTransform(base.Translate(0, 0, -0.25))
	s2.RefractiveIndex = 2

	s3 := object.GlassSphere()
	s3.SetTransform(base.Translate(0, 0, 0.25))
	s3.RefractiveIndex = 2.5

	r := ray.NewRay(base.NewPoint(0, 0, -4), base.NewVector(0, 0, 1))
	ints := object.Intersections(
		object.NewIntersection(2, s1), object.NewIntersection(2.75, s2), object.NewIntersection(3.25, s3),
		object.NewIntersection(4.75, s2), object.NewIntersection(5.25, s3), object.NewIntersection(6, s1),
	)

	tests := []struct {
		intersection *object.Intersection
		expectedN1   float64
		expectedN2   float64
	}{
		{
			intersection: ints[0],
			expectedN1:   1.0,
			expectedN2:   1.5,
		},
		{
			intersection: ints[1],
			expectedN1:   1.5,
			expectedN2:   2.0,
		},
		{
			intersection: ints[2],
			expectedN1:   2.0,
			expectedN2:   2.5,
		},
		{
			intersection: ints[3],
			expectedN1:   2.5,
			expectedN2:   2.5,
		},
		{
			intersection: ints[4],
			expectedN1:   2.5,
			expectedN2:   1.5,
		},
		{
			intersection: ints[5],
			expectedN1:   1.5,
			expectedN2:   1.0,
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			hd := prepareComputations(test.intersection, r, ints)
			g.Expect(hd.n1).To(Equal(test.expectedN1))
			g.Expect(hd.n2).To(Equal(test.expectedN2))
		})
	}
}

func TestSchlick(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// under total internal reflection
	s := object.GlassSphere()
	r := ray.NewRay(base.NewPoint(0, 0, math.Sqrt(2)/2), base.NewVector(0, 1, 0))
	ints := object.Intersections(object.NewIntersection(-math.Sqrt(2)/2, s), object.NewIntersection(math.Sqrt(2)/2, s))
	hd := prepareComputations(ints[1], r, ints)

	reflectance := schlick(hd)
	g.Expect(reflectance).To(Equal(1.0))

	// perpendicular viewing angle
	r = ray.NewRay(base.NewPoint(0, 0, 0), base.NewVector(0, 1, 0))
	ints = object.Intersections(object.NewIntersection(-1, s), object.NewIntersection(1, s))
	hd = prepareComputations(ints[1], r, ints)
	reflectance = schlick(hd)
	g.Expect(base.EqualFloats(reflectance, 0.04)).To(BeTrue())

	// small angle and n2 > n1 (looking into the distance == high reflection)
	r = ray.NewRay(base.NewPoint(0, 0.99, -2), base.NewVector(0, 0, 1))
	ints = object.Intersections(object.NewIntersection(1.8589, s))
	hd = prepareComputations(ints[0], r, ints)
	reflectance = schlick(hd)
	g.Expect(reflectance).To(Equal(0.4887308101221217))
}

func TestRender(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)
	worldTestSetup()

	w := NewWorld(testLights, testObjects)
	c := NewCamera(11, 11, math.Pi/2)
	from := base.NewPoint(0, 0, -5)
	to := base.Origin
	up := base.NewVector(0, 1, 0)
	c.SetTransform(base.ViewTransform(from, to, up))

	canvas := Render(c, w)
	expColor := image.NewColor(0.38066119308103435, 0.47582649135129296, 0.28549589481077575)
	g.Expect(canvas.PixelAt(5, 5)).To(Equal(expColor))
}
