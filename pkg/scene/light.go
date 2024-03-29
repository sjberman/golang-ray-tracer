package scene

import (
	"math"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/object"
)

// PointLight is a light with no size, existing at a single point.
type PointLight struct {
	position  *base.Tuple
	intensity *image.Color
}

// NewPointLight returns a new PointLight object.
func NewPointLight(pos *base.Tuple, intensity *image.Color) *PointLight {
	return &PointLight{
		position:  pos,
		intensity: intensity,
	}
}

// lighting returns the color at a point based on the light, material, and the eye/normal vectors.
func lighting(
	light *PointLight,
	obj object.Object,
	material *object.Material,
	point, eyev, normalv *base.Tuple,
	inShadow bool,
) *image.Color {
	color := material.Color
	if material.Pattern != nil {
		color = obj.PatternAt(point, material.Pattern)
	}
	diffuse, specular := image.Black, image.Black
	// combine surface color with light's color
	effectiveColor := color.MultiplyColor(light.intensity)

	// find the direction to the light source
	lightv := light.position.Subtract(point).Normalize()

	// compute the ambient contribution
	ambient := effectiveColor.Multiply(material.Ambient)
	if inShadow {
		return ambient
	}

	// lightDotNormal represents the cosine of the angle between the light vector
	// and the normal vector. A negative number means the light is on the
	// other side of the surface.
	lightDotNormal := lightv.DotProduct(normalv)
	if lightDotNormal >= 0 {
		// compute the diffuse contribution
		diffuse = effectiveColor.Multiply(material.Diffuse).Multiply(lightDotNormal)

		// reflectDotEye represents the cosine of the angle between the reflection vector
		// and the eye vector. A negative number means the light reflects away from the eye.
		reflectv := lightv.Negate().Reflect(normalv)
		reflectDotEye := reflectv.DotProduct(eyev)
		if reflectDotEye > 0 {
			// compute the specular contribution
			factor := math.Pow(reflectDotEye, material.Shininess)
			specular = light.intensity.Multiply(material.Specular).Multiply(factor)
		}
	}
	// Add the three contributions together to get the final shading
	return ambient.Add(diffuse.Add(specular))
}
