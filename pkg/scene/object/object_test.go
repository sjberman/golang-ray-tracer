package object

import (
	"math"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
)

func testNewObject(g *WithT, o Object) {
	g.Expect(o.GetTransform()).To(Equal(&base.Identity))
	g.Expect(o.GetMaterial()).To(Equal(&DefaultMaterial))

	t := base.Translate(2, 3, 4)
	o.SetTransform(t)
	g.Expect(o.GetTransform()).To(Equal(t))

	m := DefaultMaterial
	m.Ambient = 1
	o.SetMaterial(&m)
	g.Expect(o.GetMaterial()).To(Equal(&m))

	o.SetParent(&Group{})
	g.Expect(o.GetParent()).To(Equal(&Group{}))
}

func TestNewObject(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	o := newObject()
	g.Expect(o.GetTransform()).To(Equal(&base.Identity))
	g.Expect(o.GetMaterial()).To(Equal(&DefaultMaterial))
	g.Expect(o.parent).To(BeNil())

	tuple := base.Translate(2, 3, 4)
	o.SetTransform(tuple)
	g.Expect(o.GetTransform()).To(Equal(tuple))

	m := DefaultMaterial
	m.Ambient = 1
	o.SetMaterial(&m)
	g.Expect(o.GetMaterial()).To(Equal(&m))

	o.SetParent(&Group{})
	g.Expect(o.parent).To(Equal(&Group{}))
}

func TestPatternAt(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// with object transformation
	s := NewSphere()
	s.SetTransform(base.Scale(2, 2, 2))
	p := image.NewMockPattern()
	c := s.PatternAt(base.NewPoint(2, 3, 4), p)
	g.Expect(c).To(Equal(image.NewColor(1, 1.5, 2)))

	// with pattern transformation
	s = NewSphere()
	p = image.NewMockPattern()
	p.SetTransform(base.Scale(2, 2, 2))
	c = s.PatternAt(base.NewPoint(2, 3, 4), p)
	g.Expect(c).To(Equal(image.NewColor(1, 1.5, 2)))

	// with both object and pattern transformation
	s = NewSphere()
	s.SetTransform(base.Scale(2, 2, 2))
	p = image.NewMockPattern()
	p.SetTransform(base.Translate(0.5, 1, 1.5))
	c = s.PatternAt(base.NewPoint(2.5, 3, 3.5), p)
	g.Expect(c).To(Equal(image.NewColor(0.75, 0.5, 0.25)))
}

func TestWorldToObject(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	g1 := NewGroup()
	g1.SetTransform(base.RotateY(math.Pi / 2))
	g2 := NewGroup()
	g2.SetTransform(base.Scale(2, 2, 2))
	g1.Add(g2)

	s := NewSphere()
	s.SetTransform(base.Translate(5, 0, 0))
	g2.Add(s)

	p := s.worldToObject(base.NewPoint(-2, 0, -10))
	g.Expect(p).To(Equal(base.NewPoint(0, 0, -1)))
}

func TestNormalToWorld(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	g1 := NewGroup()
	g1.SetTransform(base.RotateY(math.Pi / 2))
	g2 := NewGroup()
	g2.SetTransform(base.Scale(1, 2, 3))
	g1.Add(g2)

	s := NewSphere()
	s.SetTransform(base.Translate(5, 0, 0))
	g2.Add(s)

	v := s.normalToWorld(base.NewVector(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3))
	g.Expect(v).To(Equal(base.NewVector(0.28571428571428575, 0.4285714285714286, -0.8571428571428572)))
}

func TestNormalAt(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	g1 := NewGroup()
	g1.SetTransform(base.RotateY(math.Pi / 2))
	g2 := NewGroup()
	g2.SetTransform(base.Scale(1, 2, 3))
	g1.Add(g2)

	s := NewSphere()
	s.SetTransform(base.Translate(5, 0, 0))
	g2.Add(s)

	n := s.NormalAt(base.NewPoint(1.7321, 1.1547, -5.5774), nil)
	g.Expect(n).To(Equal(base.NewVector(0.28570368184140726, 0.42854315178114105, -0.8571605294481017)))
}

func TestBoundsSplit(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// perfect cube
	box := &Bounds{Minimum: base.NewPoint(-1, -4, -5), Maximum: base.NewPoint(9, 6, 5)}
	left, right := box.split()
	g.Expect(left.Minimum).To(Equal(base.NewPoint(-1, -4, -5)))
	g.Expect(left.Maximum).To(Equal(base.NewPoint(4, 6, 5)))
	g.Expect(right.Minimum).To(Equal(base.NewPoint(4, -4, -5)))
	g.Expect(right.Maximum).To(Equal(base.NewPoint(9, 6, 5)))

	// x-wide box
	box = &Bounds{Minimum: base.NewPoint(-1, -2, -3), Maximum: base.NewPoint(9, 5.5, 3)}
	left, right = box.split()
	g.Expect(left.Minimum).To(Equal(base.NewPoint(-1, -2, -3)))
	g.Expect(left.Maximum).To(Equal(base.NewPoint(4, 5.5, 3)))
	g.Expect(right.Minimum).To(Equal(base.NewPoint(4, -2, -3)))
	g.Expect(right.Maximum).To(Equal(base.NewPoint(9, 5.5, 3)))

	// y-wide box
	box = &Bounds{Minimum: base.NewPoint(-1, -2, -3), Maximum: base.NewPoint(5, 8, 3)}
	left, right = box.split()
	g.Expect(left.Minimum).To(Equal(base.NewPoint(-1, -2, -3)))
	g.Expect(left.Maximum).To(Equal(base.NewPoint(5, 3, 3)))
	g.Expect(right.Minimum).To(Equal(base.NewPoint(-1, 3, -3)))
	g.Expect(right.Maximum).To(Equal(base.NewPoint(5, 8, 3)))

	// z-wide box
	box = &Bounds{Minimum: base.NewPoint(-1, -2, -3), Maximum: base.NewPoint(5, 3, 7)}
	left, right = box.split()
	g.Expect(left.Minimum).To(Equal(base.NewPoint(-1, -2, -3)))
	g.Expect(left.Maximum).To(Equal(base.NewPoint(5, 3, 2)))
	g.Expect(right.Minimum).To(Equal(base.NewPoint(-1, -2, 2)))
	g.Expect(right.Maximum).To(Equal(base.NewPoint(5, 3, 7)))
}
