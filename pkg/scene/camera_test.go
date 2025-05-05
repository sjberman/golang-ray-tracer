package scene

import (
	"fmt"
	"math"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
)

func TestNewCamera(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	c := NewCamera(160, 120, math.Pi/2)
	g.Expect(c.hsize).To(Equal(160))
	g.Expect(c.vsize).To(Equal(120))
	g.Expect(c.fieldOfView).To(Equal(math.Pi / 2))
	g.Expect(c.transform).To(Equal(&base.Identity))

	c.SetTransform(base.Scale(1, 2, 3))
	g.Expect(c.transform).To(Equal(base.Scale(1, 2, 3)))

	c = NewCamera(200, 125, math.Pi/2)
	g.Expect(c.pixelSize).To(BeNumerically("~", 0.01))

	c = NewCamera(125, 200, math.Pi/2)
	g.Expect(c.pixelSize).To(BeNumerically("~", 0.01))
}

func TestRayForPixel(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	// through center of canvas
	c := NewCamera(201, 101, math.Pi/2)
	ray := c.RayForPixel(100, 50)
	g.Expect(ray.Origin).To(Equal(base.Origin))
	g.Expect(ray.Direction).To(Equal(base.NewVector(0, 0, -1)))

	// through corner of canvas
	ray = c.RayForPixel(0, 0)
	g.Expect(ray.Origin).To(Equal(base.Origin))
	g.Expect(ray.Direction.GetX()).To(BeNumerically("~", 0.6651864261194509))
	g.Expect(ray.Direction.GetY()).To(BeNumerically("~", 0.33259321305972545))
	g.Expect(ray.Direction.GetZ()).To(BeNumerically("~", -0.6685123582500481))

	// when camera is transformed
	c.SetTransform(base.RotateY(math.Pi / 4).Multiply(base.Translate(0, -2, 5)))
	ray = c.RayForPixel(100, 50)
	g.Expect(ray.Origin).To(Equal(base.NewPoint(0, 2, -5)))
	expVector := base.NewVector(math.Sqrt(2)/2, 0, -math.Sqrt(2)/2)
	g.Expect(ray.Direction.Equals(expVector)).To(BeTrue(), fmt.Sprintf("%v", ray.Direction))
}
