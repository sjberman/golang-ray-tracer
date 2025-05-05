package object

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

func TestNewCsg(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	s := NewSphere()
	c := NewCube()
	csg := NewCsg(union, s, c)
	g.Expect(csg.transform).To(Equal(base.Identity))
	g.Expect(csg.operation).To(Equal(union))
	g.Expect(csg.left).To(Equal(s))
	g.Expect(csg.right).To(Equal(c))
	g.Expect(s.parent).To(Equal(csg))
	g.Expect(c.parent).To(Equal(csg))
}

func TestIntersectionAllowed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		op                  string
		lhit, inl, inr, res bool
	}{
		// Union
		{op: union, lhit: true, inl: true, inr: true, res: false},
		{op: union, lhit: true, inl: true, inr: false, res: true},
		{op: union, lhit: true, inl: false, inr: true, res: false},
		{op: union, lhit: true, inl: false, inr: false, res: true},
		{op: union, lhit: false, inl: true, inr: true, res: false},
		{op: union, lhit: false, inl: true, inr: false, res: false},
		{op: union, lhit: false, inl: false, inr: true, res: true},
		{op: union, lhit: false, inl: false, inr: false, res: true},
		// Intersect
		{op: intersection, lhit: true, inl: true, inr: true, res: true},
		{op: intersection, lhit: true, inl: true, inr: false, res: false},
		{op: intersection, lhit: true, inl: false, inr: true, res: true},
		{op: intersection, lhit: true, inl: false, inr: false, res: false},
		{op: intersection, lhit: false, inl: true, inr: true, res: true},
		{op: intersection, lhit: false, inl: true, inr: false, res: true},
		{op: intersection, lhit: false, inl: false, inr: true, res: false},
		{op: intersection, lhit: false, inl: false, inr: false, res: false},
		// Difference
		{op: difference, lhit: true, inl: true, inr: true, res: false},
		{op: difference, lhit: true, inl: true, inr: false, res: true},
		{op: difference, lhit: true, inl: false, inr: true, res: false},
		{op: difference, lhit: true, inl: false, inr: false, res: true},
		{op: difference, lhit: false, inl: true, inr: true, res: true},
		{op: difference, lhit: false, inl: true, inr: false, res: true},
		{op: difference, lhit: false, inl: false, inr: true, res: false},
		{op: difference, lhit: false, inl: false, inr: false, res: false},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			res := intersectionAllowed(test.op, test.lhit, test.inl, test.inr)
			g.Expect(res).To(Equal(test.res))
		})
	}
}

func TestFilterIntersections(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	s := NewSphere()
	c := NewCube()
	ints := Intersections(
		NewIntersection(1, s),
		NewIntersection(2, c),
		NewIntersection(3, s),
		NewIntersection(4, c),
	)
	csg := NewCsg(union, s, c)
	res := csg.filterIntersections(ints)
	g.Expect(len(res)).To(Equal(2))
	g.Expect(res[0]).To(Equal(ints[0]))
	g.Expect(res[1]).To(Equal(ints[3]))

	csg = NewCsg(intersection, s, c)
	res = csg.filterIntersections(ints)
	g.Expect(len(res)).To(Equal(2))
	g.Expect(res[0]).To(Equal(ints[1]))
	g.Expect(res[1]).To(Equal(ints[2]))

	csg = NewCsg(difference, s, c)
	res = csg.filterIntersections(ints)
	g.Expect(len(res)).To(Equal(2))
	g.Expect(res[0]).To(Equal(ints[0]))
	g.Expect(res[1]).To(Equal(ints[1]))
}

func TestCsgIntersect(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// ray misses
	c := NewCsg(union, NewSphere(), NewCube())
	r := ray.NewRay(base.NewPoint(0, 2, -5), base.NewVector(0, 0, 1))
	ints := c.Intersect(r)
	g.Expect(len(ints)).To(BeZero())

	// ray hits
	s1 := NewSphere()
	s2 := NewSphere()
	s2.SetTransform(base.Translate(0, 0, 0.5))
	c = NewCsg(union, s1, s2)
	r = ray.NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
	ints = c.Intersect(r)
	g.Expect(len(ints)).To(Equal(2))
	g.Expect(ints[0].Value).To(Equal(4.0))
	g.Expect(ints[0].Object).To(Equal(s1))
	g.Expect(ints[1].Value).To(Equal(6.5))
	g.Expect(ints[1].Object).To(Equal(s2))
}

func TestCsg_DivideIntoSmallerPieces(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	s1 := NewSphere()
	s1.SetTransform(base.Translate(-1.5, 0, 0))
	s2 := NewSphere()
	s2.SetTransform(base.Translate(1.5, 0, 0))
	g1 := NewGroup()
	g1.Add(s1, s2)

	s3 := NewSphere()
	s3.SetTransform(base.Translate(0, 0, -1.5))
	s4 := NewSphere()
	s4.SetTransform(base.Translate(0, 0, 1.5))
	g2 := NewGroup()
	g2.Add(s3, s4)

	csg := NewCsg(difference, g1, g2)
	csg.Divide(1)

	lgrp, ok := csg.left.(*Group)
	g.Expect(ok).To(BeTrue())
	lgrp1, ok := lgrp.Objects[0].(*Group)
	g.Expect(ok).To(BeTrue())
	lgrp2, ok := lgrp.Objects[1].(*Group)
	g.Expect(ok).To(BeTrue())
	rgrp, ok := csg.right.(*Group)
	g.Expect(ok).To(BeTrue())
	rgrp1, ok := rgrp.Objects[0].(*Group)
	g.Expect(ok).To(BeTrue())
	rgrp2, ok := rgrp.Objects[1].(*Group)
	g.Expect(ok).To(BeTrue())

	g.Expect(lgrp1.Objects[0]).To(Equal(s1))
	g.Expect(lgrp2.Objects[0]).To(Equal(s2))
	g.Expect(rgrp1.Objects[0]).To(Equal(s3))
	g.Expect(rgrp2.Objects[0]).To(Equal(s4))
}
