package object

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

func TestNewCube(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	c := NewCube()
	testNewObject(g, c)
	g.Expect(c.Bounds()).To(Equal(&Bounds{
		Minimum: base.NewPoint(-1, -1, -1),
		Maximum: base.NewPoint(1, 1, 1),
	}))
}

func TestCubeIntersect(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	c := NewCube()

	tests := []struct {
		name             string
		ray              *ray.Ray
		expVal1, expVal2 float64
	}{
		{
			name:    "+x face",
			ray:     ray.NewRay(base.NewPoint(5, 0.5, 0), base.NewVector(-1, 0, 0)),
			expVal1: 4.0,
			expVal2: 6.0,
		},
		{
			name:    "-x face",
			ray:     ray.NewRay(base.NewPoint(-5, 0.5, 0), base.NewVector(1, 0, 0)),
			expVal1: 4.0,
			expVal2: 6.0,
		},
		{
			name:    "+y face",
			ray:     ray.NewRay(base.NewPoint(0.5, 5, 0), base.NewVector(0, -1, 0)),
			expVal1: 4.0,
			expVal2: 6.0,
		},
		{
			name:    "-y face",
			ray:     ray.NewRay(base.NewPoint(0.5, -5, 0), base.NewVector(0, 1, 0)),
			expVal1: 4.0,
			expVal2: 6.0,
		},
		{
			name:    "+z face",
			ray:     ray.NewRay(base.NewPoint(0.5, 0, 5), base.NewVector(0, 0, -1)),
			expVal1: 4.0,
			expVal2: 6.0,
		},
		{
			name:    "-z face",
			ray:     ray.NewRay(base.NewPoint(0.5, 0, -5), base.NewVector(0, 0, 1)),
			expVal1: 4.0,
			expVal2: 6.0,
		},
		{
			name:    "inside",
			ray:     ray.NewRay(base.NewPoint(0, 0.5, 0), base.NewVector(0, 0, 1)),
			expVal1: -1.0,
			expVal2: 1.0,
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			ints := c.Intersect(test.ray)
			g.Expect(len(ints)).To(Equal(2))
			g.Expect(ints[0].Value).To(Equal(test.expVal1))
			g.Expect(ints[1].Value).To(Equal(test.expVal2))
		})
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
		g.Expect(len(ints)).To(BeZero())
	}
}

func TestCubeNormalAt(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	c := NewCube()

	n := c.NormalAt(base.NewPoint(1, 0.5, -0.8), nil)
	g.Expect(n).To(Equal(base.NewVector(1, 0, 0)))

	n = c.NormalAt(base.NewPoint(-1, -0.2, 0.9), nil)
	g.Expect(n).To(Equal(base.NewVector(-1, 0, 0)))

	n = c.NormalAt(base.NewPoint(-0.4, 1, -0.1), nil)
	g.Expect(n).To(Equal(base.NewVector(0, 1, 0)))

	n = c.NormalAt(base.NewPoint(0.3, -1, -0.7), nil)
	g.Expect(n).To(Equal(base.NewVector(0, -1, 0)))

	n = c.NormalAt(base.NewPoint(-0.6, 0.3, 1), nil)
	g.Expect(n).To(Equal(base.NewVector(0, 0, 1)))

	n = c.NormalAt(base.NewPoint(0.4, 0.4, -1), nil)
	g.Expect(n).To(Equal(base.NewVector(0, 0, -1)))

	n = c.NormalAt(base.NewPoint(1, 1, 1), nil)
	g.Expect(n).To(Equal(base.NewVector(1, 0, 0)))

	n = c.NormalAt(base.NewPoint(-1, -1, -1), nil)
	g.Expect(n).To(Equal(base.NewVector(-1, 0, 0)))
}
