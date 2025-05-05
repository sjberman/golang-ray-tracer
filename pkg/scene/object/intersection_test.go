package object

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestNewIntersection(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	s := NewSphere()
	i := NewIntersection(3.5, s)
	g.Expect(i.Value).To(Equal(3.5))
	g.Expect(i.Object).To(Equal(s))
}

func TestIntersections(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	s := NewSphere()
	i1 := NewIntersection(1, s)
	i2 := NewIntersection(2, s)
	ints := Intersections(i1, i2)
	g.Expect(len(ints)).To(Equal(2))
	g.Expect(ints[0].Value).To(Equal(1.0))
	g.Expect(ints[1].Value).To(Equal(2.0))
}

func TestHit(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	s := NewSphere()
	i1 := NewIntersection(1, s)
	i2 := NewIntersection(2, s)
	ints := Intersections(i1, i2)
	i := Hit(ints)
	g.Expect(i).To(Equal(i1))

	i1 = NewIntersection(-1, s)
	i2 = NewIntersection(1, s)
	ints = Intersections(i1, i2)
	i = Hit(ints)
	g.Expect(i).To(Equal(i2))

	i1 = NewIntersection(-2, s)
	i2 = NewIntersection(-1, s)
	ints = Intersections(i1, i2)
	i = Hit(ints)
	g.Expect(i).To(BeNil())

	i1 = NewIntersection(5, s)
	i2 = NewIntersection(7, s)
	i3 := NewIntersection(-3, s)
	i4 := NewIntersection(2, s)
	ints = Intersections(i1, i2, i3, i4)
	i = Hit(ints)
	g.Expect(i).To(Equal(i4))
}
