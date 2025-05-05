package scene

import (
	"math"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/object"
)

func TestNewPointLight(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	point := base.Origin
	color := image.White
	p := NewPointLight(point, color)
	g.Expect(p.position).To(Equal(point))
	g.Expect(p.intensity).To(Equal(color))
}

func TestLighting(t *testing.T) {
	t.Parallel()

	s := object.NewSphere()
	m := object.DefaultMaterial
	position := base.Origin

	tests := []struct {
		name     string
		eyev     *base.Tuple
		light    *PointLight
		inShadow bool
		expColor *image.Color
	}{
		{
			name:     "lighting with eye between light and surface",
			eyev:     base.NewVector(0, 0, -1),
			light:    NewPointLight(base.NewPoint(0, 0, -10), image.White),
			expColor: image.NewColor(1.9000000000000001, 1.9000000000000001, 1.9000000000000001),
		},
		{
			name:     "lighting with surface in shadow",
			eyev:     base.NewVector(0, 0, -1),
			light:    NewPointLight(base.NewPoint(0, 0, -10), image.White),
			inShadow: true,
			expColor: image.NewColor(0.1, 0.1, 0.1),
		},
		{
			name:     "lighting with eye between light and surface, eye offset 45 degrees",
			eyev:     base.NewVector(0, math.Sqrt(2)/2, -math.Sqrt(2)/2),
			light:    NewPointLight(base.NewPoint(0, 0, -10), image.White),
			expColor: image.NewColor(1.0, 1.0, 1.0),
		},
		{
			name:     "lighting with eye opposite surface, light offset 45 degrees",
			eyev:     base.NewVector(0, 0, -1),
			light:    NewPointLight(base.NewPoint(0, 10, -10), image.White),
			expColor: image.NewColor(0.7363961030678927, 0.7363961030678927, 0.7363961030678927),
		},
		{
			name:     "lighting with eye in the path of the reflection vector",
			eyev:     base.NewVector(0, -math.Sqrt(2)/2, -math.Sqrt(2)/2),
			light:    NewPointLight(base.NewPoint(0, 10, -10), image.White),
			expColor: image.NewColor(1.6363961030678928, 1.6363961030678928, 1.6363961030678928),
		},
		{
			name:     "lighting with light behind the surface",
			eyev:     base.NewVector(0, 0, -1),
			light:    NewPointLight(base.NewPoint(0, 0, 10), image.White),
			expColor: image.NewColor(0.1, 0.1, 0.1),
		},
	}

	normalv := base.NewVector(0, 0, -1)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			result := lighting(test.light, s, &m, position, test.eyev, normalv, test.inShadow)
			g.Expect(result).To(Equal(test.expColor))
		})
	}
}

func TestLighting_WithPattern(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	s := object.NewSphere()
	m := object.DefaultMaterial

	eyev := base.NewVector(0, 0, -1)
	normalv := base.NewVector(0, 0, -1)

	m.Pattern = image.NewStripePattern(image.White, image.Black)
	m.Ambient = 1
	m.Diffuse = 0
	m.Specular = 0
	light := NewPointLight(base.NewPoint(0, 0, -10), image.White)
	c1 := lighting(light, s, &m, base.NewPoint(0.9, 0, 0), eyev, normalv, false)
	c2 := lighting(light, s, &m, base.NewPoint(1.1, 0, 0), eyev, normalv, false)
	g.Expect(c1).To(Equal(image.White))
	g.Expect(c2).To(Equal(image.Black))
}
