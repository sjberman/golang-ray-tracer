package object

import (
	"math"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

func TestObject(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Object Suite")
}

func testNewObject(o Object) {
	Expect(o.GetTransform()).To(Equal(&base.Identity))
	Expect(o.GetMaterial()).To(Equal(&DefaultMaterial))

	t := base.Translate(2, 3, 4)
	o.SetTransform(t)
	Expect(o.GetTransform()).To(Equal(t))

	m := DefaultMaterial
	m.Ambient = 1
	o.SetMaterial(&m)
	Expect(o.GetMaterial()).To(Equal(&m))

	o.SetParent(&Group{})
	Expect(o.GetParent()).To(Equal(&Group{}))
}

type data struct {
	ray     *ray.Ray
	expVal1 float64
	expVal2 float64
}

var _ = Describe("object tests", func() {
	It("creates objects", func() {
		o := newObject()
		Expect(o.GetTransform()).To(Equal(&base.Identity))
		Expect(o.GetMaterial()).To(Equal(&DefaultMaterial))
		Expect(o.parent).To(BeNil())

		t := base.Translate(2, 3, 4)
		o.SetTransform(t)
		Expect(o.GetTransform()).To(Equal(t))

		m := DefaultMaterial
		m.Ambient = 1
		o.SetMaterial(&m)
		Expect(o.GetMaterial()).To(Equal(&m))

		o.SetParent(&Group{})
		Expect(o.parent).To(Equal(&Group{}))
	})

	It("returns the pattern at a point", func() {
		// with object transformation
		s := NewSphere()
		s.SetTransform(base.Scale(2, 2, 2))
		p := image.NewMockPattern()
		c := s.PatternAt(base.NewPoint(2, 3, 4), p)
		Expect(c).To(Equal(image.NewColor(1, 1.5, 2)))

		// with pattern transformation
		s = NewSphere()
		p = image.NewMockPattern()
		p.SetTransform(base.Scale(2, 2, 2))
		c = s.PatternAt(base.NewPoint(2, 3, 4), p)
		Expect(c).To(Equal(image.NewColor(1, 1.5, 2)))

		// with both object and pattern transformation
		s = NewSphere()
		s.SetTransform(base.Scale(2, 2, 2))
		p = image.NewMockPattern()
		p.SetTransform(base.Translate(0.5, 1, 1.5))
		c = s.PatternAt(base.NewPoint(2.5, 3, 3.5), p)
		Expect(c).To(Equal(image.NewColor(0.75, 0.5, 0.25)))
	})

	It("converts a point in world space to object space", func() {
		g1 := NewGroup()
		g1.SetTransform(base.RotateY(math.Pi / 2))
		g2 := NewGroup()
		g2.SetTransform(base.Scale(2, 2, 2))
		g1.Add(g2)

		s := NewSphere()
		s.SetTransform(base.Translate(5, 0, 0))
		g2.Add(s)

		p := s.worldToObject(base.NewPoint(-2, 0, -10))
		Expect(p).To(Equal(base.NewPoint(0, 0, -1)))
	})

	It("converts a normal from object space to world space", func() {
		g1 := NewGroup()
		g1.SetTransform(base.RotateY(math.Pi / 2))
		g2 := NewGroup()
		g2.SetTransform(base.Scale(1, 2, 3))
		g1.Add(g2)

		s := NewSphere()
		s.SetTransform(base.Translate(5, 0, 0))
		g2.Add(s)

		v := s.normalToWorld(base.NewVector(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3))
		Expect(v).To(Equal(base.NewVector(0.28571428571428575, 0.4285714285714286, -0.8571428571428572)))
	})

	It("calculates the normal of a child object", func() {
		g1 := NewGroup()
		g1.SetTransform(base.RotateY(math.Pi / 2))
		g2 := NewGroup()
		g2.SetTransform(base.Scale(1, 2, 3))
		g1.Add(g2)

		s := NewSphere()
		s.SetTransform(base.Translate(5, 0, 0))
		g2.Add(s)

		n := s.NormalAt(base.NewPoint(1.7321, 1.1547, -5.5774), nil)
		Expect(n).To(Equal(base.NewVector(0.28570368184140726, 0.42854315178114105, -0.8571605294481017)))
	})
})

var _ = Describe("bounds tests", func() {
	It("splits a bounding box", func() {
		// perfect cube
		box := &Bounds{Minimum: base.NewPoint(-1, -4, -5), Maximum: base.NewPoint(9, 6, 5)}
		left, right := box.split()
		Expect(left.Minimum).To(Equal(base.NewPoint(-1, -4, -5)))
		Expect(left.Maximum).To(Equal(base.NewPoint(4, 6, 5)))
		Expect(right.Minimum).To(Equal(base.NewPoint(4, -4, -5)))
		Expect(right.Maximum).To(Equal(base.NewPoint(9, 6, 5)))

		// x-wide box
		box = &Bounds{Minimum: base.NewPoint(-1, -2, -3), Maximum: base.NewPoint(9, 5.5, 3)}
		left, right = box.split()
		Expect(left.Minimum).To(Equal(base.NewPoint(-1, -2, -3)))
		Expect(left.Maximum).To(Equal(base.NewPoint(4, 5.5, 3)))
		Expect(right.Minimum).To(Equal(base.NewPoint(4, -2, -3)))
		Expect(right.Maximum).To(Equal(base.NewPoint(9, 5.5, 3)))

		// y-wide box
		box = &Bounds{Minimum: base.NewPoint(-1, -2, -3), Maximum: base.NewPoint(5, 8, 3)}
		left, right = box.split()
		Expect(left.Minimum).To(Equal(base.NewPoint(-1, -2, -3)))
		Expect(left.Maximum).To(Equal(base.NewPoint(5, 3, 3)))
		Expect(right.Minimum).To(Equal(base.NewPoint(-1, 3, -3)))
		Expect(right.Maximum).To(Equal(base.NewPoint(5, 8, 3)))

		// z-wide box
		box = &Bounds{Minimum: base.NewPoint(-1, -2, -3), Maximum: base.NewPoint(5, 3, 7)}
		left, right = box.split()
		Expect(left.Minimum).To(Equal(base.NewPoint(-1, -2, -3)))
		Expect(left.Maximum).To(Equal(base.NewPoint(5, 3, 2)))
		Expect(right.Minimum).To(Equal(base.NewPoint(-1, -2, 2)))
		Expect(right.Maximum).To(Equal(base.NewPoint(5, 3, 7)))
	})
})
