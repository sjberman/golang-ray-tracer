package scene

import (
	"fmt"
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
)

// type MockObject struct {
// 	*ObjectImpl
// }

// func NewMockObject() *MockObject {
// 	return &MockObject{}
// }

// func (m *MockObject) GetMaterial() *Material                            { return &m.Material }
// func (m *MockObject) GetTransform() *base.Matrix                        { return &m.transform }
// func (m *MockObject) SetTransform(matrix ...*base.Matrix)               {}
// func (m *MockObject) SetMaterial(mat *Material)                         {}
// func (m *MockObject) patternAt(*base.Tuple, image.Pattern) *image.Color { return nil }
// func (m *MockObject) Intersect(*Ray) []*Intersection {
// 	return nil
// }
// func (m *MockObject) normalAt(*base.Tuple) *base.Tuple {
// 	return nil
// }

type MockPattern struct {
	*image.PatternObject
}

func NewMockPattern() *MockPattern {
	return &MockPattern{
		PatternObject: image.NewPattern(nil, nil, mockFunc),
	}
}

func mockFunc(point *base.Tuple, p *image.PatternObject) *image.Color {
	return image.NewColor(point.GetX(), point.GetY(), point.GetZ())
}

type data struct {
	ray     *Ray
	expVal1 float64
	expVal2 float64
}

var _ = Describe("object tests", func() {
	testNewObject := func(o Object) {
		Expect(o.GetTransform()).To(Equal(&base.Identity))
		Expect(o.GetMaterial()).To(Equal(&DefaultMaterial))

		t := base.Translate(2, 3, 4)
		o.SetTransform(t)
		Expect(o.GetTransform()).To(Equal(t))

		m := DefaultMaterial
		m.Ambient = 1
		o.SetMaterial(&m)
		Expect(o.GetMaterial()).To(Equal(&m))
	}

	It("creates objects", func() {
		o := newObject()
		Expect(o.GetTransform()).To(Equal(&base.Identity))
		Expect(o.GetMaterial()).To(Equal(&DefaultMaterial))

		t := base.Translate(2, 3, 4)
		o.SetTransform(t)
		Expect(o.GetTransform()).To(Equal(t))

		m := DefaultMaterial
		m.Ambient = 1
		o.SetMaterial(&m)
		Expect(o.GetMaterial()).To(Equal(&m))
	})

	It("returns the pattern at a point", func() {
		// with object transformation
		s := NewSphere()
		s.SetTransform(base.Scale(2, 2, 2))
		p := NewMockPattern()
		c := s.patternAt(base.NewPoint(2, 3, 4), p)
		Expect(c).To(Equal(image.NewColor(1, 1.5, 2)))

		// with pattern transformation
		s = NewSphere()
		p = NewMockPattern()
		p.SetTransform(base.Scale(2, 2, 2))
		c = s.patternAt(base.NewPoint(2, 3, 4), p)
		Expect(c).To(Equal(image.NewColor(1, 1.5, 2)))

		// with both object and pattern transformation
		s = NewSphere()
		s.SetTransform(base.Scale(2, 2, 2))
		p = NewMockPattern()
		p.SetTransform(base.Translate(0.5, 1, 1.5))
		c = s.patternAt(base.NewPoint(2.5, 3, 3.5), p)
		Expect(c).To(Equal(image.NewColor(0.75, 0.5, 0.25)))
	})

	Context("planes", func() {
		It("creates planes", func() {
			p := NewPlane()
			testNewObject(p)
		})

		It("calculates a plane intersection", func() {
			// ray parallel to plane
			p := NewPlane()
			r := NewRay(base.NewPoint(0, 10, 0), base.NewVector(0, 0, 1))
			ints := p.intersect(r)
			Expect(len(ints)).To(BeZero())

			// coplanar ray
			r = NewRay(base.NewPoint(0, 0, 0), base.NewVector(0, 0, 1))
			ints = p.intersect(r)
			Expect(len(ints)).To(BeZero())

			// intersects plane from above
			r = NewRay(base.NewPoint(0, 1, 0), base.NewVector(0, -1, 0))
			ints = p.intersect(r)
			Expect(len(ints)).To(Equal(1))
			Expect(ints[0].value).To(Equal(1.0))
			Expect(ints[0].GetObject()).To(Equal(p))

			// intersects plane from below
			r = NewRay(base.NewPoint(0, -1, 0), base.NewVector(0, 1, 0))
			ints = p.intersect(r)
			Expect(len(ints)).To(Equal(1))
			Expect(ints[0].value).To(Equal(1.0))
			Expect(ints[0].GetObject()).To(Equal(p))
		})

		It("computes the surface normal", func() {
			p := NewPlane()
			constVector := base.NewVector(0, 1, 0)

			n := p.normalAt(base.NewPoint(0, 0, 0))
			Expect(n).To(Equal(constVector))
			n = p.normalAt(base.NewPoint(10, 0, -10))
			Expect(n).To(Equal(constVector))
			n = p.normalAt(base.NewPoint(-5, 0, 150))
			Expect(n).To(Equal(constVector))
		})
	})

	Context("cubes", func() {
		It("creates cubes", func() {
			c := NewCube()
			testNewObject(c)
		})

		It("calculates a cube intersection", func() {
			c := NewCube()

			testCases := []data{
				{
					// +x face
					ray:     NewRay(base.NewPoint(5, 0.5, 0), base.NewVector(-1, 0, 0)),
					expVal1: 4.0,
					expVal2: 6.0,
				},
				{
					// -x face
					ray:     NewRay(base.NewPoint(-5, 0.5, 0), base.NewVector(1, 0, 0)),
					expVal1: 4.0,
					expVal2: 6.0,
				},
				{
					// +y face
					ray:     NewRay(base.NewPoint(0.5, 5, 0), base.NewVector(0, -1, 0)),
					expVal1: 4.0,
					expVal2: 6.0,
				},
				{
					// -y face
					ray:     NewRay(base.NewPoint(0.5, -5, 0), base.NewVector(0, 1, 0)),
					expVal1: 4.0,
					expVal2: 6.0,
				},
				{
					// +z face
					ray:     NewRay(base.NewPoint(0.5, 0, 5), base.NewVector(0, 0, -1)),
					expVal1: 4.0,
					expVal2: 6.0,
				},
				{
					// -z face
					ray:     NewRay(base.NewPoint(0.5, 0, -5), base.NewVector(0, 0, 1)),
					expVal1: 4.0,
					expVal2: 6.0,
				},
				{
					// inside
					ray:     NewRay(base.NewPoint(0, 0.5, 0), base.NewVector(0, 0, 1)),
					expVal1: -1.0,
					expVal2: 1.0,
				},
			}

			for _, t := range testCases {
				ints := c.intersect(t.ray)
				Expect(len(ints)).To(Equal(2),
					fmt.Sprintf("origin: %v, direction: %v", t.ray.Origin, t.ray.Direction))
				Expect(ints[0].value).To(Equal(t.expVal1))
				Expect(ints[1].value).To(Equal(t.expVal2))
			}

			// ray misses the cube
			rays := []*Ray{
				NewRay(base.NewPoint(-2, 0, 0), base.NewVector(0.2673, 0.5345, 0.8018)),
				NewRay(base.NewPoint(0, -2, 0), base.NewVector(0.8018, 0.2673, 0.5345)),
				NewRay(base.NewPoint(0, 0, -2), base.NewVector(0.5345, 0.8018, 0.2673)),
				NewRay(base.NewPoint(2, 0, 2), base.NewVector(0, 0, -1)),
				NewRay(base.NewPoint(0, 2, 2), base.NewVector(0, -1, 0)),
				NewRay(base.NewPoint(2, 2, 0), base.NewVector(-1, 0, 0)),
			}

			for _, r := range rays {
				ints := c.intersect(r)
				Expect(len(ints)).To(BeZero())
			}
		})

		It("computes the surface normal", func() {
			c := NewCube()

			n := c.normalAt(base.NewPoint(1, 0.5, -0.8))
			Expect(n).To(Equal(base.NewVector(1, 0, 0)))

			n = c.normalAt(base.NewPoint(-1, -0.2, 0.9))
			Expect(n).To(Equal(base.NewVector(-1, 0, 0)))

			n = c.normalAt(base.NewPoint(-0.4, 1, -0.1))
			Expect(n).To(Equal(base.NewVector(0, 1, 0)))

			n = c.normalAt(base.NewPoint(0.3, -1, -0.7))
			Expect(n).To(Equal(base.NewVector(0, -1, 0)))

			n = c.normalAt(base.NewPoint(-0.6, 0.3, 1))
			Expect(n).To(Equal(base.NewVector(0, 0, 1)))

			n = c.normalAt(base.NewPoint(0.4, 0.4, -1))
			Expect(n).To(Equal(base.NewVector(0, 0, -1)))

			n = c.normalAt(base.NewPoint(1, 1, 1))
			Expect(n).To(Equal(base.NewVector(1, 0, 0)))

			n = c.normalAt(base.NewPoint(-1, -1, -1))
			Expect(n).To(Equal(base.NewVector(-1, 0, 0)))
		})
	})

	Context("cylinders", func() {
		It("creates cylinders", func() {
			c := NewCylinder()
			testNewObject(c)
			Expect(c.Minimum).To(Equal(math.Inf(-1)))
			Expect(c.Maximum).To(Equal(math.Inf(0)))
			Expect(c.Closed).To(BeFalse())
		})

		It("calculates a cylinder intersection", func() {
			// misses
			c := NewCylinder()
			r := NewRay(base.NewPoint(1, 0, 0), base.NewVector(0, 1, 0))
			ints := c.intersect(r)
			Expect(len(ints)).To(BeZero())

			r = NewRay(base.NewPoint(0, 0, 0), base.NewVector(0, 1, 0))
			ints = c.intersect(r)
			Expect(len(ints)).To(BeZero())

			r = NewRay(base.NewPoint(0, 0, -5), base.NewVector(1, 1, 1))
			ints = c.intersect(r)
			Expect(len(ints)).To(BeZero())

			// hits
			testCases := []data{
				{
					ray:     NewRay(base.NewPoint(1, 0, -5), base.NewVector(0, 0, 1)),
					expVal1: 5.0,
					expVal2: 5.0,
				},
				{
					ray:     NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1)),
					expVal1: 4.0,
					expVal2: 6.0,
				},
				{
					ray:     NewRay(base.NewPoint(0.5, 0, -5), base.NewVector(0.1, 1, 1).Normalize()),
					expVal1: 6.80798191702732,
					expVal2: 7.088723439378861,
				},
			}

			for _, t := range testCases {
				ints := c.intersect(t.ray)
				Expect(len(ints)).To(Equal(2),
					fmt.Sprintf("origin: %v, direction: %v", t.ray.Origin, t.ray.Direction))
				Expect(ints[0].value).To(Equal(t.expVal1))
				Expect(ints[1].value).To(Equal(t.expVal2))
			}

			// constrained cylinder
			c.Minimum = 1
			c.Maximum = 2
			rays := []*Ray{
				NewRay(base.NewPoint(0, 1.5, 0), base.NewVector(0.1, 1, 0).Normalize()),
				NewRay(base.NewPoint(0, 3, -5), base.NewVector(0, 0, 1)),
				NewRay(base.NewPoint(0, 0, -5), base.NewVector(0, 0, 1)),
				NewRay(base.NewPoint(0, 2, -5), base.NewVector(0, 0, 1)),
				NewRay(base.NewPoint(0, 1, -5), base.NewVector(0, 0, 1)),
			}

			for _, r := range rays {
				ints := c.intersect(r)
				Expect(len(ints)).To(BeZero(),
					fmt.Sprintf("origin: %v, direction: %v", r.Origin, r.Direction))
			}
			r = NewRay(base.NewPoint(0, 1.5, -2), base.NewVector(0, 0, 1))
			Expect(len(c.intersect(r))).To(Equal(2))

			// intersection of caps on a closed cylinder
			c.Closed = true
			rays = []*Ray{
				NewRay(base.NewPoint(0, 3, 0), base.NewVector(0, -1, 0)),
				NewRay(base.NewPoint(0, 3, -2), base.NewVector(0, -1, 2)),
				NewRay(base.NewPoint(0, 4, -2), base.NewVector(0, -1, 1)),
				NewRay(base.NewPoint(0, 0, -2), base.NewVector(0, 1, 2)),
				NewRay(base.NewPoint(0, -1, -2), base.NewVector(0, 1, 1)),
			}

			for _, r := range rays {
				ints := c.intersect(r)
				Expect(len(ints)).To(Equal(2),
					fmt.Sprintf("origin: %v, direction: %v", r.Origin, r.Direction))
			}
		})

		It("computes the surface normal", func() {
			c := NewCylinder()

			n := c.normalAt(base.NewPoint(1, 0, 0))
			Expect(n).To(Equal(base.NewVector(1, 0, 0)))

			n = c.normalAt(base.NewPoint(0, 5, -1))
			Expect(n).To(Equal(base.NewVector(0, 0, -1)))

			n = c.normalAt(base.NewPoint(0, -2, 1))
			Expect(n).To(Equal(base.NewVector(0, 0, 1)))

			n = c.normalAt(base.NewPoint(-1, 1, 0))
			Expect(n).To(Equal(base.NewVector(-1, 0, 0)))

			// cylinder caps
			c.Minimum = 1
			c.Maximum = 2
			c.Closed = true
			n = c.normalAt(base.NewPoint(0, 1, 0))
			Expect(n).To(Equal(base.NewVector(0, -1, 0)))

			n = c.normalAt(base.NewPoint(0.5, 1, 0))
			Expect(n).To(Equal(base.NewVector(0, -1, 0)))

			n = c.normalAt(base.NewPoint(0, 1, 0.5))
			Expect(n).To(Equal(base.NewVector(0, -1, 0)))

			n = c.normalAt(base.NewPoint(0, 2, 0))
			Expect(n).To(Equal(base.NewVector(0, 1, 0)))

			n = c.normalAt(base.NewPoint(0.5, 2, 0))
			Expect(n).To(Equal(base.NewVector(0, 1, 0)))

			n = c.normalAt(base.NewPoint(0, 2, 0.5))
			Expect(n).To(Equal(base.NewVector(0, 1, 0)))
		})
	})
})
