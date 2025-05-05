package object

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

func TestNewTriangle(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	p1 := base.NewPoint(0, 1, 0)
	p2 := base.NewPoint(-1, 0, 0)
	p3 := base.NewPoint(1, 0, 0)

	triangle := NewTriangle(p1, p2, p3)
	testNewObject(g, triangle)

	g.Expect(triangle.P1).To(Equal(p1))
	g.Expect(triangle.P2).To(Equal(p2))
	g.Expect(triangle.P3).To(Equal(p3))
	g.Expect(triangle.edge1).To(Equal(base.NewVector(-1, -1, 0)))
	g.Expect(triangle.edge2).To(Equal(base.NewVector(1, -1, 0)))
	g.Expect(triangle.normal).To(Equal(base.NewVector(0, 0, -1)))

	g.Expect(triangle.Bounds()).To(Equal(&Bounds{
		Minimum: base.NewPoint(-1, 0, 0),
		Maximum: base.NewPoint(1, 1, 0),
	}))
}

func TestTriangleIntersect(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// ray is parallel (misses)
	triangle := NewTriangle(base.NewPoint(0, 1, 0), base.NewPoint(-1, 0, 0), base.NewPoint(1, 0, 0))
	r := ray.NewRay(base.NewPoint(0, -1, -2), base.NewVector(0, 1, 0))
	ints := triangle.Intersect(r)
	g.Expect(len(ints)).To(BeZero())

	// ray misses the P1-P3 edge
	r = ray.NewRay(base.NewPoint(1, 1, -2), base.NewVector(0, 0, 1))
	ints = triangle.Intersect(r)
	g.Expect(len(ints)).To(BeZero())

	// ray misses the P1-P2 edge
	r = ray.NewRay(base.NewPoint(-1, 1, -2), base.NewVector(0, 0, 1))
	ints = triangle.Intersect(r)
	g.Expect(len(ints)).To(BeZero())

	// ray misses the P2-P3 edge
	r = ray.NewRay(base.NewPoint(0, -1, -2), base.NewVector(0, 0, 1))
	ints = triangle.Intersect(r)
	g.Expect(len(ints)).To(BeZero())

	// hit
	r = ray.NewRay(base.NewPoint(0, 0.5, -2), base.NewVector(0, 0, 1))
	ints = triangle.Intersect(r)
	g.Expect(len(ints)).To(Equal(1))
	g.Expect(ints[0].Value).To(Equal(2.0))
}

func TestTriangleNormalAt(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	triangle := NewTriangle(base.NewPoint(0, 1, 0), base.NewPoint(-1, 0, 0), base.NewPoint(1, 0, 0))
	n1 := triangle.NormalAt(base.NewPoint(0, 0.5, 0), nil)
	n2 := triangle.NormalAt(base.NewPoint(-0.5, 0.75, 0), nil)
	n3 := triangle.NormalAt(base.NewPoint(0.5, 0.25, 0), nil)

	g.Expect(n1).To(Equal(triangle.normal))
	g.Expect(n2).To(Equal(triangle.normal))
	g.Expect(n3).To(Equal(triangle.normal))
}

func TestSmoothTriangle(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	p1 := base.NewPoint(0, 1, 0)
	p2 := base.NewPoint(-1, 0, 0)
	p3 := base.NewPoint(1, 0, 0)
	n1 := base.NewVector(0, 1, 0)
	n2 := base.NewVector(-1, 0, 0)
	n3 := base.NewVector(1, 0, 0)
	tri := NewSmoothTriangle(p1, p2, p3, n1, n2, n3)

	g.Expect(tri.P1).To(Equal(p1))
	g.Expect(tri.P2).To(Equal(p2))
	g.Expect(tri.P3).To(Equal(p3))
	g.Expect(tri.N1).To(Equal(n1))
	g.Expect(tri.N2).To(Equal(n2))
	g.Expect(tri.N3).To(Equal(n3))

	intersection := NewIntersection(1, tri)
	intersection.u = 0.45
	intersection.v = 0.25

	n := tri.NormalAt(base.NewPoint(0, 0, 0), intersection)
	g.Expect(n).To(Equal(base.NewVector(-0.5547001962252291, 0.8320502943378437, 0)))
}
