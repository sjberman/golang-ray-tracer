package scene

import (
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
)

var _ = Describe("world tests", func() {
	var light *PointLight
	var objects []Object

	BeforeEach(func() {
		light = NewPointLight(base.NewPoint(-10, 10, -10), image.White)
		s1 := NewSphere()
		s1.GetMaterial().SetColor(image.NewColor(0.8, 1.0, 0.6))
		s1.GetMaterial().SetDiffuse(0.7)
		s1.GetMaterial().SetSpecular(0.2)
		s2 := NewSphere()
		s2.SetTransform(base.Scale(0.5, 0.5, 0.5))
		objects = []Object{s1, s2}
	})

	It("creates worlds", func() {
		w := NewWorld(light, objects)
		Expect(w.light).To(Equal(light))
		Expect(w.objects).To(Equal(objects))
	})

	It("computes intersections in a world", func() {
		w := NewWorld(light, objects)
		ray := NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
		ints := w.intersect(ray)

		Expect(len(ints)).To(Equal(4))
		Expect(ints[0].value).To(Equal(4.0))
		Expect(ints[1].value).To(Equal(4.5))
		Expect(ints[2].value).To(Equal(5.5))
		Expect(ints[3].value).To(Equal(6.0))
	})

	It("computes the color at a ray intersection", func() {
		// when ray misses
		w := NewWorld(light, objects)
		ray := NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 1, 0))
		color := w.ColorAt(ray, remainingReflections)
		Expect(color).To(Equal(image.Black))

		// when ray hits
		ray = NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
		color = w.ColorAt(ray, remainingReflections)
		Expect(color).To(Equal(image.NewColor(0.38066119308103435, 0.47582649135129296, 0.28549589481077575)))

		// intersection behind the ray
		outer := w.objects[0]
		outer.GetMaterial().SetAmbient(1)
		inner := w.objects[1]
		inner.GetMaterial().SetAmbient(1)

		ray = NewRay(base.NewPoint(0, 0, 0.75), base.NewVector(0, 0, -1))
		color = w.ColorAt(ray, remainingReflections)
		Expect(color).To(Equal(inner.GetMaterial().color))

		// with mutally reflective surfaces (assume no infinite recursion)
		lower := NewPlane()
		lower.SetReflective(1)
		lower.SetTransform(base.Translate(0, -1, 0))

		upper := NewPlane()
		upper.SetReflective(1)
		upper.SetTransform(base.Translate(0, 1, 0))

		w.objects = []Object{lower, upper}
		ray = NewRay(base.NewPoint(0, 0, 0), base.NewVector(0, 1, 0))
		Eventually(func() *image.Color {
			return w.ColorAt(ray, remainingReflections)
		}).ShouldNot(BeNil())
	})

	It("determines if a point is shadowed", func() {
		// nothing collinear with point and light
		w := NewWorld(light, objects)
		p := base.NewPoint(0, 10, 0)
		Expect(w.isShadowed(p)).To(BeFalse())

		// object between point and light
		p = base.NewPoint(10, -10, 10)
		Expect(w.isShadowed(p)).To(BeTrue())

		// object is behind the light
		p = base.NewPoint(-20, 20, -20)
		Expect(w.isShadowed(p)).To(BeFalse())

		// object is behind the point
		p = base.NewPoint(-2, 2, -2)
		Expect(w.isShadowed(p)).To(BeFalse())
	})

	It("returns the reflected color", func() {
		// nonreflective material
		w := NewWorld(light, objects)
		ray := NewRay(base.NewPoint(0, 0, 0), base.NewVector(0, 0, 1))
		s := w.objects[1]
		s.GetMaterial().SetAmbient(1)
		intersection := NewIntersection(1, s)
		hd := prepareComputations(intersection, ray, intersections(intersection))
		color := w.reflectedColor(hd, 1)
		Expect(color).To(Equal(image.Black))

		// reflective material with max recursive depth
		p := NewPlane()
		p.SetReflective(0.5)
		p.SetTransform(base.Translate(0, -1, 0))
		w.objects = append(w.objects, p)
		ray = NewRay(base.NewPoint(0, 0, -3), base.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
		intersection = NewIntersection(math.Sqrt(2), p)
		hd = prepareComputations(intersection, ray, intersections(intersection))
		color = w.reflectedColor(hd, 0)
		Expect(color).To(Equal(image.Black))

		// reflective material
		color = w.reflectedColor(hd, 1)
		Expect(color).To(Equal(image.NewColor(0.19033404421300906, 0.23791755526626135, 0.14275053315975678)))

		// shadeHit with reflective material
		color = w.shadeHit(hd, remainingReflections)
		Expect(color).To(Equal(image.NewColor(0.8767594331945104, 0.9243429442477628, 0.8291759221412582)))
	})

	Context("refraction", func() {
		Specify("opaque surface", func() {
			w := NewWorld(light, objects)
			ray := NewRay(base.NewPoint(0, 0, 0), base.NewVector(0, 0, 1))
			s := w.objects[0]
			ints := intersections(NewIntersection(4, s), NewIntersection(6, s))
			hd := prepareComputations(ints[0], ray, ints)
			color := w.refractedColor(hd, 5)
			Expect(color).To(Equal(image.Black))
		})

		Specify("refracted color with max recursive depth", func() {
			w := NewWorld(light, objects)
			s := w.objects[0]
			s.GetMaterial().SetTransparency(1)
			s.GetMaterial().SetRefractiveIndex(1.5)
			ray := NewRay(base.NewPoint(0, 0, 0), base.NewVector(0, 0, 1))
			ints := intersections(NewIntersection(4, s), NewIntersection(6, s))
			hd := prepareComputations(ints[0], ray, ints)
			color := w.refractedColor(hd, 0)
			Expect(color).To(Equal(image.Black))

		})

		Specify("refracted color under total internal reflection", func() {
			w := NewWorld(light, objects)
			s := w.objects[0]
			ray := NewRay(base.NewPoint(0, 0, math.Sqrt(2)/2), base.NewVector(0, 1, 0))
			ints := intersections(NewIntersection(-math.Sqrt(2)/2, s), NewIntersection(math.Sqrt(2)/2, s))
			// inside of sphere, so look at second intersection
			hd := prepareComputations(ints[1], ray, ints)
			color := w.refractedColor(hd, 5)
			Expect(color).To(Equal(image.Black))

		})

		Specify("refracted color with a refracted ray", func() {
			w := NewWorld(light, objects)
			s1 := w.objects[0]
			s1.GetMaterial().SetAmbient(1)
			s1.GetMaterial().SetPattern(NewMockPattern())
			s2 := w.objects[1]
			s2.GetMaterial().SetTransparency(1)
			s2.GetMaterial().SetRefractiveIndex(1.5)
			ray := NewRay(base.NewPoint(0, 0, 0.1), base.NewVector(0, 1, 0))
			ints := intersections(
				NewIntersection(-0.9899, s1), NewIntersection(-0.4899, s2),
				NewIntersection(0.4899, s2), NewIntersection(0.9899, s1))
			hd := prepareComputations(ints[2], ray, ints)
			color := w.refractedColor(hd, 5)
			Expect(color).To(Equal(image.NewColor(0, 0.9988845862650526, 0.04721846378372032)))

		})
	})

	It("calculates the shadeHit with reflective, transparent material", func() {
		w := NewWorld(light, objects)
		floor := NewPlane()
		floor.SetTransform(base.Translate(0, -1, 0))
		floor.SetTransparency(0.5)
		floor.SetRefractiveIndex(1.5)
		ball := NewSphere()
		ball.SetColor(image.NewColor(1, 0, 0))
		ball.SetAmbient(0.5)
		ball.SetTransform(base.Translate(0, -3.5, -0.5))
		w.objects = append(w.objects, floor, ball)

		ray := NewRay(base.NewPoint(0, 0, -3), base.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
		ints := intersections(NewIntersection(math.Sqrt(2), floor))
		hd := prepareComputations(ints[0], ray, ints)
		color := w.shadeHit(hd, 5)
		Expect(color).To(Equal(image.NewColor(0.9364253889815014, 0.6864253889815014, 0.6864253889815014)))
	})

	It("creates useful hit data", func() {
		ray := NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
		intersection := NewIntersection(4, NewSphere())
		hd := prepareComputations(intersection, ray, intersections(intersection))
		Expect(hd.value).To(Equal(intersection.GetValue()))
		Expect(hd.object).To(Equal(intersection.GetObject()))
		Expect(hd.point).To(Equal(base.NewPoint(0, 0, -1)))
		Expect(hd.eyev).To(Equal(base.NewVector(0, 0, -1)))
		Expect(hd.normalv).To(Equal(base.NewVector(0, 0, -1)))
		Expect(hd.reflectv).To(Equal(base.NewVector(0, 0, -1)))
		Expect(hd.inside).To(BeFalse())

		// hit should offset the point
		s := NewSphere()
		s.SetTransform(base.Translate(0, 0, 1))
		intersection = NewIntersection(5, s)
		hd = prepareComputations(intersection, ray, intersections(intersection))
		Expect(hd.overPoint.GetZ()).To(BeNumerically("<", -base.Epsilon/2))
		Expect(hd.point.GetZ()).To(BeNumerically(">", hd.overPoint.GetZ()))

		// intersection occurs on inside
		ray = NewRay(base.Origin, base.NewVector(0, 0, 1))
		intersection = NewIntersection(1, NewSphere())
		hd = prepareComputations(intersection, ray, intersections(intersection))
		Expect(hd.point).To(Equal(base.NewPoint(0, 0, 1)))
		Expect(hd.eyev).To(Equal(base.NewVector(0, 0, -1)))
		Expect(hd.normalv).To(Equal(base.NewVector(0, 0, -1)))
		Expect(hd.inside).To(BeTrue())

		// reflect off plane
		ray = NewRay(base.NewPoint(0, 1, -1), base.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
		intersection = NewIntersection(math.Sqrt(2), NewPlane())
		hd = prepareComputations(intersection, ray, intersections(intersection))
		Expect(hd.reflectv).To(Equal(base.NewVector(0, math.Sqrt(2)/2, math.Sqrt(2)/2)))

		// under point is offset below the surface
		ray = NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
		s = GlassSphere()
		s.SetTransform(base.Translate(0, 0, 1))
		intersection = NewIntersection(5, s)
		hd = prepareComputations(intersection, ray, intersections(intersection))
		Expect(hd.underPoint.GetZ()).To(BeNumerically(">", base.Epsilon/2))
		Expect(hd.point.GetZ()).To(BeNumerically("<", hd.underPoint.GetZ()))
	})

	It("calculates refractive indices of multiple intersections", func() {
		s1 := GlassSphere()
		s1.SetTransform(base.Scale(2, 2, 2))

		s2 := GlassSphere()
		s2.SetTransform(base.Translate(0, 0, -0.25))
		s2.SetRefractiveIndex(2)

		s3 := GlassSphere()
		s3.SetTransform(base.Translate(0, 0, 0.25))
		s3.SetRefractiveIndex(2.5)

		ray := NewRay(base.NewPoint(0, 0, -4), base.NewVector(0, 0, 1))
		ints := intersections(
			NewIntersection(2, s1), NewIntersection(2.75, s2), NewIntersection(3.25, s3),
			NewIntersection(4.75, s2), NewIntersection(5.25, s3), NewIntersection(6, s1))
		type data struct {
			intersection *Intersection
			expectedN1   float64
			expectedN2   float64
		}
		testCases := []data{
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

		for _, t := range testCases {
			hd := prepareComputations(t.intersection, ray, ints)
			Expect(hd.n1).To(Equal(t.expectedN1))
			Expect(hd.n2).To(Equal(t.expectedN2))
		}
	})

	It("computes the Schlick approximation", func() {
		// under total internal reflection
		s := GlassSphere()
		ray := NewRay(base.NewPoint(0, 0, math.Sqrt(2)/2), base.NewVector(0, 1, 0))
		ints := intersections(NewIntersection(-math.Sqrt(2)/2, s), NewIntersection(math.Sqrt(2)/2, s))
		hd := prepareComputations(ints[1], ray, ints)

		reflectance := schlick(hd)
		Expect(reflectance).To(Equal(1.0))

		// perpendicular viewing angle
		ray = NewRay(base.NewPoint(0, 0, 0), base.NewVector(0, 1, 0))
		ints = intersections(NewIntersection(-1, s), NewIntersection(1, s))
		hd = prepareComputations(ints[1], ray, ints)
		reflectance = schlick(hd)
		Expect(base.EqualFloats(reflectance, 0.04)).To(BeTrue())

		// small angle and n2 > n1 (looking into the distance == high reflection)
		ray = NewRay(base.NewPoint(0, 0.99, -2), base.NewVector(0, 0, 1))
		ints = intersections(NewIntersection(1.8589, s))
		hd = prepareComputations(ints[0], ray, ints)
		reflectance = schlick(hd)
		Expect(reflectance).To(Equal(0.4887308101221217))
	})

	It("renders the world", func() {
		w := NewWorld(light, objects)
		c := NewCamera(11, 11, math.Pi/2)
		from := base.NewPoint(0, 0, -5)
		to := base.Origin
		up := base.NewVector(0, 1, 0)
		c.SetTransform(base.ViewTransform(from, to, up))

		canvas := Render(c, w)
		expColor := image.NewColor(0.38066119308103435, 0.47582649135129296, 0.28549589481077575)
		Expect(canvas.PixelAt(5, 5)).To(Equal(expColor))
	})
})
