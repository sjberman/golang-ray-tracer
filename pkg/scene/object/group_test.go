package object

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

func TestNewGroup(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	group := NewGroup()
	g.Expect(group.transform).To(Equal(base.Identity))
	g.Expect(len(group.Objects)).To(BeZero())
}

func TestAdd(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	group := NewGroup()
	s := NewSphere()
	group.Add(s)
	g.Expect(group.Objects).To(ContainElement(s))
	g.Expect(s.parent).To(Equal(group))
}

func TestGroupIntersect(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// empty group
	group := NewGroup()
	r := ray.NewRay(base.NewPoint(0, 0, 0), base.NewVector(0, 0, 1))
	ints := group.Intersect(r)
	g.Expect(len(ints)).To(BeZero())

	// non-empty group
	s1 := NewSphere()
	s2 := NewSphere()
	s2.SetTransform(base.Translate(0, 0, -3))
	s3 := NewSphere()
	s3.SetTransform(base.Translate(5, 0, 0))
	group.Add(s1, s2, s3)
	r = ray.NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1))
	ints = group.Intersect(r)
	g.Expect(len(ints)).To(Equal(4))
	g.Expect(ints[0].Object).To(Equal(s2))
	g.Expect(ints[1].Object).To(Equal(s2))
	g.Expect(ints[2].Object).To(Equal(s1))
	g.Expect(ints[3].Object).To(Equal(s1))

	// transformed group
	group = NewGroup()
	group.SetTransform(base.Scale(2, 2, 2))
	s := NewSphere()
	s.SetTransform(base.Translate(5, 0, 0))
	group.Add(s)
	r = ray.NewRay(base.NewPoint(10, 0, -10), base.NewVector(0, 0, 1))
	ints = group.Intersect(r)
	g.Expect(len(ints)).To(Equal(2))
}

func TestBounds(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	group := NewGroup()
	s := NewSphere()
	s.SetTransform(base.Translate(2, 5, -3), base.Scale(2, 2, 2))
	c := NewCylinder()
	c.Minimum = -2
	c.Maximum = 2
	c.SetTransform(base.Translate(-4, -1, 4), base.Scale(0.5, 1, 0.5))
	group.Add(s, c)

	box := group.Bounds()
	g.Expect(box.Minimum).To(Equal(base.NewPoint(-4.5, -3, -5)))
	g.Expect(box.Maximum).To(Equal(base.NewPoint(4, 7, 4.5)))
}

func TestDivide(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	s1 := NewSphere()
	s1.SetTransform(base.Translate(-2, -2, 0))
	s2 := NewSphere()
	s2.SetTransform(base.Translate(-2, 2, 0))
	s3 := NewSphere()
	s3.SetTransform(base.Scale(4, 4, 4))

	group := NewGroup()
	group.Add(s1, s2, s3)
	group.Divide(1)
	g.Expect(group.Objects[0]).To(Equal(s3))
	subgroup, ok := group.Objects[1].(*Group)
	g.Expect(ok).To(BeTrue())
	g.Expect(len(subgroup.Objects)).To(Equal(2))
	g1, ok := subgroup.Objects[0].(*Group)
	g.Expect(ok).To(BeTrue())
	g2, ok := subgroup.Objects[1].(*Group)
	g.Expect(ok).To(BeTrue())
	g.Expect(g1.Objects[0]).To(Equal(s1))
	g.Expect(g2.Objects[0]).To(Equal(s2))

	// too few children in top group
	s1.SetTransform(base.Translate(-2, 0, 0))
	s2.SetTransform(base.Translate(2, 1, 0))
	s3.SetTransform(base.Translate(2, -1, 0))
	s4 := NewSphere()

	subgroup = NewGroup()
	subgroup.Add(s1, s2, s3)
	group = NewGroup()
	group.Add(subgroup, s4)
	group.Divide(3)
	g.Expect(group.Objects[0]).To(Equal(subgroup))
	g.Expect(group.Objects[1]).To(Equal(s4))
	g.Expect(len(subgroup.Objects)).To(Equal(2))
	g1, ok = subgroup.Objects[0].(*Group)
	g.Expect(ok).To(BeTrue())
	g2, ok = subgroup.Objects[1].(*Group)
	g.Expect(ok).To(BeTrue())
	g.Expect(g1.Objects[0]).To(Equal(s1))
	g.Expect(g2.Objects).To(Equal([]Object{s2, s3}))
}

func TestPartitionChildren(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	group := NewGroup()
	s1 := NewSphere()
	s1.SetTransform(base.Translate(-2, 0, 0))
	s2 := NewSphere()
	s2.SetTransform(base.Translate(2, 0, 0))
	s3 := NewSphere()
	group.Add(s1, s2, s3)

	left, right := group.partitionChildren()
	g.Expect(len(group.Objects)).To(Equal(1))
	g.Expect(group.Objects).To(ContainElement(s3))
	g.Expect(left[0]).To(Equal(s1))
	g.Expect(right[0]).To(Equal(s2))
}

func TestMakeSubgroup(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	s1 := NewSphere()
	s2 := NewSphere()
	group := NewGroup()
	group.makeSubgroup([]Object{s1, s2})
	g.Expect(len(group.Objects)).To(Equal(1))
	subgroup, ok := group.Objects[0].(*Group)
	g.Expect(ok).To(BeTrue())
	g.Expect(subgroup.Objects).To(Equal([]Object{s1, s2}))
}
