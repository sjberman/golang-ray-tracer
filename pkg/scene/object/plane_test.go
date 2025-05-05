package object

import (
	"math"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

func TestNewPlane(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	p := NewPlane()
	testNewObject(g, p)
	g.Expect(p.Bounds()).To(Equal(&Bounds{
		Minimum: base.NewPoint(math.Inf(-1), 0, math.Inf(-1)),
		Maximum: base.NewPoint(math.Inf(1), 0, math.Inf(1)),
	}))
}

func TesPlaneIntersect(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// ray parallel to plane
	p := NewPlane()
	r := ray.NewRay(base.NewPoint(0, 10, 0), base.NewVector(0, 0, 1))
	ints := p.Intersect(r)
	g.Expect(len(ints)).To(BeZero())

	// coplanar ray
	r = ray.NewRay(base.NewPoint(0, 0, 0), base.NewVector(0, 0, 1))
	ints = p.Intersect(r)
	g.Expect(len(ints)).To(BeZero())

	// intersects plane from above
	r = ray.NewRay(base.NewPoint(0, 1, 0), base.NewVector(0, -1, 0))
	ints = p.Intersect(r)
	g.Expect(len(ints)).To(Equal(1))
	g.Expect(ints[0].Value).To(Equal(1.0))
	g.Expect(ints[0].Object).To(Equal(p))

	// intersects plane from below
	r = ray.NewRay(base.NewPoint(0, -1, 0), base.NewVector(0, 1, 0))
	ints = p.Intersect(r)
	g.Expect(len(ints)).To(Equal(1))
	g.Expect(ints[0].Value).To(Equal(1.0))
	g.Expect(ints[0].Object).To(Equal(p))
}

func TestPlaneNormalAt(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	p := NewPlane()
	constVector := base.NewVector(0, 1, 0)

	n := p.NormalAt(base.NewPoint(0, 0, 0), nil)
	g.Expect(n).To(Equal(constVector))
	n = p.NormalAt(base.NewPoint(10, 0, -10), nil)
	g.Expect(n).To(Equal(constVector))
	n = p.NormalAt(base.NewPoint(-5, 0, 150), nil)
	g.Expect(n).To(Equal(constVector))
}
