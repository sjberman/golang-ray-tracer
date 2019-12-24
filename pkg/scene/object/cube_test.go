package object

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

var _ = Describe("cube tests", func() {
	It("creates cubes", func() {
		c := NewCube()
		testNewObject(c)
		Expect(c.Bounds()).To(Equal(&bounds{
			minimum: base.NewPoint(-1, -1, -1),
			maximum: base.NewPoint(1, 1, 1)}))
	})

	It("calculates a cube intersection", func() {
		c := NewCube()

		testCases := []data{
			{
				// +x face
				ray:     ray.NewRay(base.NewPoint(5, 0.5, 0), base.NewVector(-1, 0, 0)),
				expVal1: 4.0,
				expVal2: 6.0,
			},
			{
				// -x face
				ray:     ray.NewRay(base.NewPoint(-5, 0.5, 0), base.NewVector(1, 0, 0)),
				expVal1: 4.0,
				expVal2: 6.0,
			},
			{
				// +y face
				ray:     ray.NewRay(base.NewPoint(0.5, 5, 0), base.NewVector(0, -1, 0)),
				expVal1: 4.0,
				expVal2: 6.0,
			},
			{
				// -y face
				ray:     ray.NewRay(base.NewPoint(0.5, -5, 0), base.NewVector(0, 1, 0)),
				expVal1: 4.0,
				expVal2: 6.0,
			},
			{
				// +z face
				ray:     ray.NewRay(base.NewPoint(0.5, 0, 5), base.NewVector(0, 0, -1)),
				expVal1: 4.0,
				expVal2: 6.0,
			},
			{
				// -z face
				ray:     ray.NewRay(base.NewPoint(0.5, 0, -5), base.NewVector(0, 0, 1)),
				expVal1: 4.0,
				expVal2: 6.0,
			},
			{
				// inside
				ray:     ray.NewRay(base.NewPoint(0, 0.5, 0), base.NewVector(0, 0, 1)),
				expVal1: -1.0,
				expVal2: 1.0,
			},
		}

		for _, t := range testCases {
			ints := c.Intersect(t.ray)
			Expect(len(ints)).To(Equal(2),
				fmt.Sprintf("origin: %v, direction: %v", t.ray.Origin, t.ray.Direction))
			Expect(ints[0].Value).To(Equal(t.expVal1))
			Expect(ints[1].Value).To(Equal(t.expVal2))
		}

		// ray misses the cube
		rays := []*ray.Ray{
			ray.NewRay(base.NewPoint(-2, 0, 0), base.NewVector(0.2673, 0.5345, 0.8018)),
			ray.NewRay(base.NewPoint(0, -2, 0), base.NewVector(0.8018, 0.2673, 0.5345)),
			ray.NewRay(base.NewPoint(0, 0, -2), base.NewVector(0.5345, 0.8018, 0.2673)),
			ray.NewRay(base.NewPoint(2, 0, 2), base.NewVector(0, 0, -1)),
			ray.NewRay(base.NewPoint(0, 2, 2), base.NewVector(0, -1, 0)),
			ray.NewRay(base.NewPoint(2, 2, 0), base.NewVector(-1, 0, 0)),
		}

		for _, r := range rays {
			ints := c.Intersect(r)
			Expect(len(ints)).To(BeZero())
		}
	})

	It("computes the surface normal", func() {
		c := NewCube()

		n := c.NormalAt(base.NewPoint(1, 0.5, -0.8))
		Expect(n).To(Equal(base.NewVector(1, 0, 0)))

		n = c.NormalAt(base.NewPoint(-1, -0.2, 0.9))
		Expect(n).To(Equal(base.NewVector(-1, 0, 0)))

		n = c.NormalAt(base.NewPoint(-0.4, 1, -0.1))
		Expect(n).To(Equal(base.NewVector(0, 1, 0)))

		n = c.NormalAt(base.NewPoint(0.3, -1, -0.7))
		Expect(n).To(Equal(base.NewVector(0, -1, 0)))

		n = c.NormalAt(base.NewPoint(-0.6, 0.3, 1))
		Expect(n).To(Equal(base.NewVector(0, 0, 1)))

		n = c.NormalAt(base.NewPoint(0.4, 0.4, -1))
		Expect(n).To(Equal(base.NewVector(0, 0, -1)))

		n = c.NormalAt(base.NewPoint(1, 1, 1))
		Expect(n).To(Equal(base.NewVector(1, 0, 0)))

		n = c.NormalAt(base.NewPoint(-1, -1, -1))
		Expect(n).To(Equal(base.NewVector(-1, 0, 0)))
	})
})
