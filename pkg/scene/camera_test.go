package scene

import (
	"fmt"
	"math"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sjberman/golang-ray-tracer/pkg/base"
)

func TestScene(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scene Suite")
}

var _ = Describe("camera tests", func() {
	It("creates cameras", func() {
		c := NewCamera(160, 120, math.Pi/2)
		Expect(c.hsize).To(Equal(160))
		Expect(c.vsize).To(Equal(120))
		Expect(c.fieldOfView).To(Equal(math.Pi / 2))
		Expect(c.transform).To(Equal(&base.Identity))

		c.SetTransform(base.Scale(1, 2, 3))
		Expect(c.transform).To(Equal(base.Scale(1, 2, 3)))

		c = NewCamera(200, 125, math.Pi/2)
		Expect(c.pixelSize).To(Equal(0.01))

		c = NewCamera(125, 200, math.Pi/2)
		Expect(c.pixelSize).To(Equal(0.01))
	})

	It("constructs rays to a canvas", func() {
		// through center of canvas
		c := NewCamera(201, 101, math.Pi/2)
		ray := c.RayForPixel(100, 50)
		Expect(ray.Origin).To(Equal(base.Origin))
		Expect(ray.Direction).To(Equal(base.NewVector(0, 0, -1)))

		// through corner of canvas
		ray = c.RayForPixel(0, 0)
		Expect(ray.Origin).To(Equal(base.Origin))
		expVector := base.NewVector(0.6651864261194509, 0.33259321305972545, -0.6685123582500481)
		Expect(ray.Direction).To(Equal(expVector))

		// when camera is transformed
		c.SetTransform(base.RotateY(math.Pi / 4).Multiply(base.Translate(0, -2, 5)))
		ray = c.RayForPixel(100, 50)
		Expect(ray.Origin).To(Equal(base.NewPoint(0, 2, -5)))
		expVector = base.NewVector(math.Sqrt(2)/2, 0, -math.Sqrt(2)/2)
		Expect(ray.Direction.Equals(expVector)).To(BeTrue(), fmt.Sprintf("%v", ray.Direction))
	})
})
