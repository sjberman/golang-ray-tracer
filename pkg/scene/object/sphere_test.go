package object

import (
	"math"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

func TestNewSphere(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	s := NewSphere()
	testNewObject(g, s)
	g.Expect(s.Bounds()).To(Equal(&Bounds{
		Minimum: base.NewPoint(-1, -1, -1),
		Maximum: base.NewPoint(1, 1, 1),
	}))

	s = GlassSphere()
	g.Expect(s.Transparency).To(Equal(1.0))
	g.Expect(s.RefractiveIndex).To(Equal(1.5))
}

func TestSphereIntersect(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	r := ray.NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
	s := NewSphere()

	ints := s.Intersect(r)
	g.Expect(len(ints)).To(Equal(2))
	g.Expect(ints[0].Value).To(Equal(4.0))
	g.Expect(ints[1].Value).To(Equal(6.0))
	g.Expect(ints[0].Object).To(Equal(s))
	g.Expect(ints[1].Object).To(Equal(s))

	// tangent
	r = ray.NewRay(base.NewPoint(0, 1, -5), base.NewVector(0, 0, 1))
	s = NewSphere()

	ints = s.Intersect(r)
	g.Expect(len(ints)).To(Equal(2))
	g.Expect(ints[0].Value).To(Equal(5.0))
	g.Expect(ints[1].Value).To(Equal(5.0))

	// too high, no  intersection
	r = ray.NewRay(base.NewPoint(0, 2, -5), base.NewVector(0, 0, 1))
	s = NewSphere()

	ints = s.Intersect(r)
	g.Expect(len(ints)).To(Equal(0))

	// ray starts within the sphere
	r = ray.NewRay(base.Origin, base.NewVector(0, 0, 1))
	s = NewSphere()

	ints = s.Intersect(r)
	g.Expect(len(ints)).To(Equal(2))
	g.Expect(ints[0].Value).To(Equal(-1.0))
	g.Expect(ints[1].Value).To(Equal(1.0))

	// ray starts past the sphere
	r = ray.NewRay(base.NewPoint(0, 0, 5), base.NewVector(0, 0, 1))
	s = NewSphere()

	ints = s.Intersect(r)
	g.Expect(len(ints)).To(Equal(2))
	g.Expect(ints[0].Value).To(Equal(-6.0))
	g.Expect(ints[1].Value).To(Equal(-4.0))

	//  Intersect a scaled sphere
	r = ray.NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
	s = NewSphere()
	s.SetTransform(base.Scale(2, 2, 2))
	ints = s.Intersect(r)
	g.Expect(len(ints)).To(Equal(2))
	g.Expect(ints[0].Value).To(Equal(3.0))
	g.Expect(ints[1].Value).To(Equal(7.0))

	r = ray.NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
	s = NewSphere()
	s.SetTransform(base.Translate(5, 0, 0))
	ints = s.Intersect(r)
	g.Expect(len(ints)).To(BeZero())
}

func TestSphereNormalAt(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// x axis
	s := NewSphere()
	n := s.NormalAt(base.NewPoint(1, 0, 0), nil)
	g.Expect(n).To(Equal(base.NewVector(1, 0, 0)))

	// y axis
	n = s.NormalAt(base.NewPoint(0, 1, 0), nil)
	g.Expect(n).To(Equal(base.NewVector(0, 1, 0)))

	// z axis
	n = s.NormalAt(base.NewPoint(0, 0, 1), nil)
	g.Expect(n).To(Equal(base.NewVector(0, 0, 1)))

	// non axis
	n = s.NormalAt(base.NewPoint(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3), nil)
	expVector := base.NewVector(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3)
	g.Expect(n.GetX()).To(BeNumerically("~", expVector.GetX()))
	g.Expect(n.GetY()).To(BeNumerically("~", expVector.GetY()))
	g.Expect(n.GetZ()).To(BeNumerically("~", expVector.GetZ()))

	// surface normal is a normalized vector
	g.Expect(n).To(Equal(n.Normalize()))

	// translated sphere
	s.SetTransform(base.Translate(0, 1, 0))
	n = s.NormalAt(base.NewPoint(0, 1.70711, -0.70711), nil)
	g.Expect(n).To(Equal(base.NewVector(0, 0.7071067811865475, -0.7071067811865476)))

	// scaled/rotated sphere
	m := base.Scale(1, 0.5, 1).Multiply(base.RotateZ(math.Pi / 5))
	s.SetTransform(m)
	n = s.NormalAt(base.NewPoint(0, math.Sqrt(2)/2, -math.Sqrt(2)/2), nil)
	g.Expect(n).To(Equal(base.NewVector(0, 0.9701425001453319, -0.24253562503633286)))
}
