package object

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

func TestObject(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Object Suite")
}

func testNewObject(o Object) {
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

type data struct {
	ray     *ray.Ray
	expVal1 float64
	expVal2 float64
}

var _ = Describe("object tests", func() {
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
		p := image.NewMockPattern()
		c := s.PatternAt(base.NewPoint(2, 3, 4), p)
		Expect(c).To(Equal(image.NewColor(1, 1.5, 2)))

		// with pattern transformation
		s = NewSphere()
		p = image.NewMockPattern()
		p.SetTransform(base.Scale(2, 2, 2))
		c = s.PatternAt(base.NewPoint(2, 3, 4), p)
		Expect(c).To(Equal(image.NewColor(1, 1.5, 2)))

		// with both object and pattern transformation
		s = NewSphere()
		s.SetTransform(base.Scale(2, 2, 2))
		p = image.NewMockPattern()
		p.SetTransform(base.Translate(0.5, 1, 1.5))
		c = s.PatternAt(base.NewPoint(2.5, 3, 3.5), p)
		Expect(c).To(Equal(image.NewColor(0.75, 0.5, 0.25)))
	})
})
