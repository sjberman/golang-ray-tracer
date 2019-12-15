package scene

import "github.com/sjberman/golang-ray-tracer/pkg/image"

// Material contains the attributes of a surface material
type Material struct {
	color           *image.Color
	pattern         image.Pattern
	ambient         float64
	diffuse         float64
	specular        float64
	shininess       float64
	reflective      float64
	transparency    float64
	refractiveIndex float64
}

var defaultMaterial = Material{
	color:           image.White,
	ambient:         0.1,
	diffuse:         0.9,
	specular:        0.9,
	shininess:       200,
	reflective:      0,
	transparency:    0,
	refractiveIndex: 1.0,
}

// NewMaterial returns a new Material object
func NewMaterial(
	color *image.Color,
	pattern image.Pattern,
	ambient, diffuse, specular, shininess,
	reflective, transparency, refractiveIndex float64,
) *Material {
	return &Material{
		color:           color,
		pattern:         pattern,
		ambient:         ambient,
		diffuse:         diffuse,
		specular:        specular,
		shininess:       shininess,
		reflective:      reflective,
		transparency:    transparency,
		refractiveIndex: refractiveIndex,
	}
}

// SetColor sets the material's color field
func (m *Material) SetColor(color *image.Color) {
	m.color = color
}

// SetPattern sets the material's color pattern
func (m *Material) SetPattern(pattern image.Pattern) {
	m.pattern = pattern
}

// SetAmbient sets the material's ambient field
func (m *Material) SetAmbient(ambient float64) {
	m.ambient = ambient
}

// SetDiffuse sets the material's diffuse field
func (m *Material) SetDiffuse(diffuse float64) {
	m.diffuse = diffuse
}

// SetSpecular sets the material's specular field
func (m *Material) SetSpecular(specular float64) {
	m.specular = specular
}

// SetShininess sets the material's shininess field
func (m *Material) SetShininess(shininess float64) {
	m.shininess = shininess
}

// SetReflective sets the material's reflective field
func (m *Material) SetReflective(reflective float64) {
	m.reflective = reflective
}

// SetTransparency sets the material's transparency field
func (m *Material) SetTransparency(transparency float64) {
	m.transparency = transparency
}

// SetRefractiveIndex sets the material's refractiveIndex field
func (m *Material) SetRefractiveIndex(refractiveIndex float64) {
	m.refractiveIndex = refractiveIndex
}
