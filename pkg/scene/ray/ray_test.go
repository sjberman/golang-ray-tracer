package ray

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
)

func TestNewRay(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	o := base.NewPoint(1, 2, 3)
	d := base.NewVector(4, 5, 6)
	ray := NewRay(o, d)
	g.Expect(ray.Origin).To(Equal(o))
	g.Expect(ray.Direction).To(Equal(d))
}

func TestPosition(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	r := NewRay(base.NewPoint(2, 3, 4), base.NewVector(1, 0, 0))
	expPoint := base.NewPoint(2, 3, 4)
	g.Expect(r.Position(0)).To(Equal(expPoint))

	expPoint = base.NewPoint(3, 3, 4)
	g.Expect(r.Position(1)).To(Equal(expPoint))

	expPoint = base.NewPoint(1, 3, 4)
	g.Expect(r.Position(-1)).To(Equal(expPoint))

	expPoint = base.NewPoint(4.5, 3, 4)
	g.Expect(r.Position(2.5)).To(Equal(expPoint))
}

func TestTransform(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	r := NewRay(base.NewPoint(1, 2, 3), base.NewVector(0, 1, 0))
	m := base.Translate(3, 4, 5)
	r2 := r.Transform(m)
	g.Expect(r2.Origin).To(Equal(base.NewPoint(4, 6, 8)))
	g.Expect(r2.Direction).To(Equal(base.NewVector(0, 1, 0)))

	m = base.Scale(2, 3, 4)
	r2 = r.Transform(m)
	g.Expect(r2.Origin).To(Equal(base.NewPoint(2, 6, 12)))
	g.Expect(r2.Direction).To(Equal(base.NewVector(0, 3, 0)))
}
