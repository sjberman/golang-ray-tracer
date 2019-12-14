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
		color := w.ColorAt(ray)
		Expect(color).To(Equal(image.Black))

		// when ray hits
		ray = NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
		color = w.ColorAt(ray)
		Expect(color).To(Equal(image.NewColor(0.38066119308103435, 0.47582649135129296, 0.28549589481077575)))

		// intersection behind the ray
		outer := w.objects[0]
		outer.GetMaterial().SetAmbient(1)
		inner := w.objects[1]
		inner.GetMaterial().SetAmbient(1)

		ray = NewRay(base.NewPoint(0, 0, 0.75), base.NewVector(0, 0, -1))
		color = w.ColorAt(ray)
		Expect(color).To(Equal(inner.GetMaterial().color))
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

	It("creates useful hit data", func() {
		ray := NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
		intersection := NewIntersection(4, NewSphere())
		hd := prepareComputations(intersection, ray)
		Expect(hd.value).To(Equal(intersection.GetValue()))
		Expect(hd.object).To(Equal(intersection.GetObject()))
		Expect(hd.point).To(Equal(base.NewPoint(0, 0, -1)))
		Expect(hd.eyev).To(Equal(base.NewVector(0, 0, -1)))
		Expect(hd.normalv).To(Equal(base.NewVector(0, 0, -1)))
		Expect(hd.inside).To(BeFalse())

		// hit should offset the point
		s := NewSphere()
		s.SetTransform(base.Translate(0, 0, 1))
		intersection = NewIntersection(5, s)
		hd = prepareComputations(intersection, ray)
		Expect(hd.overPoint.GetZ()).To(BeNumerically("<", -base.Epsilon/2))
		Expect(hd.point.GetZ()).To(BeNumerically(">", hd.overPoint.GetZ()))

		// intersection occurs on inside
		ray = NewRay(base.Origin, base.NewVector(0, 0, 1))
		intersection = NewIntersection(1, NewSphere())
		hd = prepareComputations(intersection, ray)
		Expect(hd.point).To(Equal(base.NewPoint(0, 0, 1)))
		Expect(hd.eyev).To(Equal(base.NewVector(0, 0, -1)))
		Expect(hd.normalv).To(Equal(base.NewVector(0, 0, -1)))
		Expect(hd.inside).To(BeTrue())
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
