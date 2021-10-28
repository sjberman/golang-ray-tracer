package object

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

var _ = Describe("csg tests", func() {
	It("creates csgs", func() {
		s := NewSphere()
		c := NewCube()
		csg := NewCsg(union, s, c)
		Expect(csg.transform).To(Equal(base.Identity))
		Expect(csg.operation).To(Equal(union))
		Expect(csg.left).To(Equal(s))
		Expect(csg.right).To(Equal(c))
		Expect(s.parent).To(Equal(csg))
		Expect(c.parent).To(Equal(csg))
	})

	It("determines the rule for a CSG operation", func() {
		type data struct {
			op                  string
			lhit, inl, inr, res bool
		}
		testCases := []data{
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
		for _, t := range testCases {
			res := intersectionAllowed(t.op, t.lhit, t.inl, t.inr)
			Expect(res).To(Equal(t.res))
		}
	})

	It("filters a list of intersections", func() {
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
		Expect(len(res)).To(Equal(2))
		Expect(res[0]).To(Equal(ints[0]))
		Expect(res[1]).To(Equal(ints[3]))

		csg = NewCsg(intersection, s, c)
		res = csg.filterIntersections(ints)
		Expect(len(res)).To(Equal(2))
		Expect(res[0]).To(Equal(ints[1]))
		Expect(res[1]).To(Equal(ints[2]))

		csg = NewCsg(difference, s, c)
		res = csg.filterIntersections(ints)
		Expect(len(res)).To(Equal(2))
		Expect(res[0]).To(Equal(ints[0]))
		Expect(res[1]).To(Equal(ints[1]))
	})

	It("calculates intersections", func() {
		// ray misses
		c := NewCsg(union, NewSphere(), NewCube())
		r := ray.NewRay(base.NewPoint(0, 2, -5), base.NewVector(0, 0, 1))
		ints := c.Intersect(r)
		Expect(len(ints)).To(BeZero())

		// ray hits
		s1 := NewSphere()
		s2 := NewSphere()
		s2.SetTransform(base.Translate(0, 0, 0.5))
		c = NewCsg(union, s1, s2)
		r = ray.NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
		ints = c.Intersect(r)
		Expect(len(ints)).To(Equal(2))
		Expect(ints[0].Value).To(Equal(4.0))
		Expect(ints[0].Object).To(Equal(s1))
		Expect(ints[1].Value).To(Equal(6.5))
		Expect(ints[1].Object).To(Equal(s2))
	})

	It("divides into smaller pieces", func() {
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
		Expect(ok).To(BeTrue())
		lgrp1, ok := lgrp.Objects[0].(*Group)
		Expect(ok).To(BeTrue())
		lgrp2, ok := lgrp.Objects[1].(*Group)
		Expect(ok).To(BeTrue())
		rgrp, ok := csg.right.(*Group)
		Expect(ok).To(BeTrue())
		rgrp1, ok := rgrp.Objects[0].(*Group)
		Expect(ok).To(BeTrue())
		rgrp2, ok := rgrp.Objects[1].(*Group)
		Expect(ok).To(BeTrue())

		Expect(lgrp1.Objects[0]).To(Equal(s1))
		Expect(lgrp2.Objects[0]).To(Equal(s2))
		Expect(rgrp1.Objects[0]).To(Equal(s3))
		Expect(rgrp2.Objects[0]).To(Equal(s4))
	})
})
