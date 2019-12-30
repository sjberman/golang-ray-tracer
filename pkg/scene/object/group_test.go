package object

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

var _ = Describe("group tests", func() {
	It("creates groups", func() {
		g := NewGroup()
		Expect(g.transform).To(Equal(base.Identity))
		Expect(len(g.Objects)).To(BeZero())
	})

	It("adds a child", func() {
		g := NewGroup()
		s := NewSphere()
		g.Add(s)
		Expect(g.Objects).To(ContainElement(s))
		Expect(s.parent).To(Equal(g))
	})

	It("calculates intersections", func() {
		// empty group
		g := NewGroup()
		r := ray.NewRay(base.NewPoint(0, 0, 0), base.NewVector(0, 0, 1))
		ints := g.Intersect(r)
		Expect(len(ints)).To(BeZero())

		// non-empty group
		s1 := NewSphere()
		s2 := NewSphere()
		s2.SetTransform(base.Translate(0, 0, -3))
		s3 := NewSphere()
		s3.SetTransform(base.Translate(5, 0, 0))
		g.Add(s1, s2, s3)
		r = ray.NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
		ints = g.Intersect(r)
		Expect(len(ints)).To(Equal(4))
		Expect(ints[0].Object).To(Equal(s2))
		Expect(ints[1].Object).To(Equal(s2))
		Expect(ints[2].Object).To(Equal(s1))
		Expect(ints[3].Object).To(Equal(s1))

		// transformed group
		g = NewGroup()
		g.SetTransform(base.Scale(2, 2, 2))
		s := NewSphere()
		s.SetTransform(base.Translate(5, 0, 0))
		g.Add(s)
		r = ray.NewRay(base.NewPoint(10, 0, -10), base.NewVector(0, 0, 1))
		ints = g.Intersect(r)
		Expect(len(ints)).To(Equal(2))
	})

	It("computes a bounding box", func() {
		g := NewGroup()
		s := NewSphere()
		s.SetTransform(base.Translate(2, 5, -3), base.Scale(2, 2, 2))
		c := NewCylinder()
		c.Minimum = -2
		c.Maximum = 2
		c.SetTransform(base.Translate(-4, -1, 4), base.Scale(0.5, 1, 0.5))
		g.Add(s, c)

		box := g.Bounds()
		Expect(box.Minimum).To(Equal(base.NewPoint(-4.5, -3, -5)))
		Expect(box.Maximum).To(Equal(base.NewPoint(4, 7, 4.5)))
	})

	It("divides into smaller pieces", func() {
		s1 := NewSphere()
		s1.SetTransform(base.Translate(-2, -2, 0))
		s2 := NewSphere()
		s2.SetTransform(base.Translate(-2, 2, 0))
		s3 := NewSphere()
		s3.SetTransform(base.Scale(4, 4, 4))

		g := NewGroup()
		g.Add(s1, s2, s3)
		g.Divide(1)
		Expect(g.Objects[0]).To(Equal(s3))
		subgroup := g.Objects[1].(*Group)
		Expect(len(subgroup.Objects)).To(Equal(2))
		g1 := subgroup.Objects[0].(*Group)
		g2 := subgroup.Objects[1].(*Group)
		Expect(g1.Objects[0]).To(Equal(s1))
		Expect(g2.Objects[0]).To(Equal(s2))

		// too few children in top group
		s1.SetTransform(base.Translate(-2, 0, 0))
		s2.SetTransform(base.Translate(2, 1, 0))
		s3.SetTransform(base.Translate(2, -1, 0))
		s4 := NewSphere()

		subgroup = NewGroup()
		subgroup.Add(s1, s2, s3)
		g = NewGroup()
		g.Add(subgroup, s4)
		g.Divide(3)
		Expect(g.Objects[0]).To(Equal(subgroup))
		Expect(g.Objects[1]).To(Equal(s4))
		Expect(len(subgroup.Objects)).To(Equal(2))
		g1 = subgroup.Objects[0].(*Group)
		g2 = subgroup.Objects[1].(*Group)
		Expect(g1.Objects[0]).To(Equal(s1))
		Expect(g2.Objects).To(Equal([]Object{s2, s3}))
	})

	It("partitions children", func() {
		g := NewGroup()
		s1 := NewSphere()
		s1.SetTransform(base.Translate(-2, 0, 0))
		s2 := NewSphere()
		s2.SetTransform(base.Translate(2, 0, 0))
		s3 := NewSphere()
		g.Add(s1, s2, s3)

		left, right := g.partitionChildren()
		Expect(len(g.Objects)).To(Equal(1))
		Expect(g.Objects).To(ContainElement(s3))
		Expect(left[0]).To(Equal(s1))
		Expect(right[0]).To(Equal(s2))
	})

	It("creates a subgroup from a list of children", func() {
		s1 := NewSphere()
		s2 := NewSphere()
		g := NewGroup()
		g.makeSubgroup([]Object{s1, s2})
		Expect(len(g.Objects)).To(Equal(1))
		subgroup := g.Objects[0].(*Group)
		Expect(subgroup.Objects).To(Equal([]Object{s1, s2}))
	})
})
