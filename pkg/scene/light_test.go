package scene

import (
	"fmt"
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
)

var _ = Describe("light tests", func() {
	It("creates point lights", func() {
		point := base.Origin
		color := image.NewColor(1, 1, 1)
		p := NewPointLight(point, color)
		Expect(p.position).To(Equal(point))
		Expect(p.intensity).To(Equal(color))
	})

	It("calculates the lighting", func() {
		m := &defaultMaterial
		position := base.Origin

		// lighting with eye between light and surface
		eyev := base.NewVector(0, 0, -1)
		normalv := base.NewVector(0, 0, -1)
		light := NewPointLight(base.NewPoint(0, 0, -10), image.NewColor(1, 1, 1))

		result := Lighting(light, m, position, eyev, normalv, false)
		expColor := image.NewColor(1.9, 1.9, 1.9)
		Expect(result.Equals(expColor)).To(BeTrue(), fmt.Sprintf("%v", result))

		// lighting with surface in shadow
		result = Lighting(light, m, position, eyev, normalv, true)
		expColor = image.NewColor(0.1, 0.1, 0.1)
		Expect(result.Equals(expColor)).To(BeTrue(), fmt.Sprintf("%v", result))

		// lighting with eye between light and surface, eye offset 45 degrees
		eyev = base.NewVector(0, math.Sqrt(2)/2, -math.Sqrt(2)/2)
		normalv = base.NewVector(0, 0, -1)
		light = NewPointLight(base.NewPoint(0, 0, -10), image.NewColor(1, 1, 1))

		result = Lighting(light, m, position, eyev, normalv, false)
		expColor = image.NewColor(1.0, 1.0, 1.0)
		Expect(result.Equals(expColor)).To(BeTrue(), fmt.Sprintf("%v", result))

		// lighting with eye opposite surface, light offset 45 degrees
		eyev = base.NewVector(0, 0, -1)
		normalv = base.NewVector(0, 0, -1)
		light = NewPointLight(base.NewPoint(0, 10, -10), image.NewColor(1, 1, 1))

		result = Lighting(light, m, position, eyev, normalv, false)
		expColor = image.NewColor(0.7363961030678927, 0.7363961030678927, 0.7363961030678927)
		Expect(result.Equals(expColor)).To(BeTrue(), fmt.Sprintf("%v", result))

		// lighting with eye in the path of the reflection vector
		eyev = base.NewVector(0, -math.Sqrt(2)/2, -math.Sqrt(2)/2)
		normalv = base.NewVector(0, 0, -1)
		light = NewPointLight(base.NewPoint(0, 10, -10), image.NewColor(1, 1, 1))

		result = Lighting(light, m, position, eyev, normalv, false)
		expColor = image.NewColor(1.6363961030678928, 1.6363961030678928, 1.6363961030678928)
		Expect(result.Equals(expColor)).To(BeTrue(), fmt.Sprintf("%v", result))

		// lighting with light behind the surface
		eyev = base.NewVector(0, 0, -1)
		normalv = base.NewVector(0, 0, -1)
		light = NewPointLight(base.NewPoint(0, 0, 10), image.NewColor(1, 1, 1))

		result = Lighting(light, m, position, eyev, normalv, false)
		expColor = image.NewColor(0.1, 0.1, 0.1)
		Expect(result.Equals(expColor)).To(BeTrue(), fmt.Sprintf("%v", result))
	})
})
