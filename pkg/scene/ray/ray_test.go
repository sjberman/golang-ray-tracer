package ray

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sjberman/golang-ray-tracer/pkg/base"
)

func TestRay(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ray Suite")
}

var _ = Describe("ray tests", func() {
	It("creates rays", func() {
		o := base.NewPoint(1, 2, 3)
		d := base.NewVector(4, 5, 6)
		ray := NewRay(o, d)
		Expect(ray.Origin).To(Equal(o))
		Expect(ray.Direction).To(Equal(d))
	})

	It("computes the position of a point on the ray", func() {
		r := NewRay(base.NewPoint(2, 3, 4), base.NewVector(1, 0, 0))
		expPoint := base.NewPoint(2, 3, 4)
		Expect(r.Position(0)).To(Equal(expPoint))

		expPoint = base.NewPoint(3, 3, 4)
		Expect(r.Position(1)).To(Equal(expPoint))

		expPoint = base.NewPoint(1, 3, 4)
		Expect(r.Position(-1)).To(Equal(expPoint))

		expPoint = base.NewPoint(4.5, 3, 4)
		Expect(r.Position(2.5)).To(Equal(expPoint))
	})

	It("transforms a ray", func() {
		r := NewRay(base.NewPoint(1, 2, 3), base.NewVector(0, 1, 0))
		m := base.Translate(3, 4, 5)
		r2 := r.Transform(m)
		Expect(r2.Origin).To(Equal(base.NewPoint(4, 6, 8)))
		Expect(r2.Direction).To(Equal(base.NewVector(0, 1, 0)))

		m = base.Scale(2, 3, 4)
		r2 = r.Transform(m)
		Expect(r2.Origin).To(Equal(base.NewPoint(2, 6, 12)))
		Expect(r2.Direction).To(Equal(base.NewVector(0, 3, 0)))
	})
})
