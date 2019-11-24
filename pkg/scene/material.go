package scene

import "github.com/sjberman/golang-ray-tracer/pkg/image"

// Material contains the attributes of a surface material
type Material struct {
	color     *image.Color
	ambient   float64
	diffuse   float64
	specular  float64
	shininess float64
}

var defaultMaterial = Material{
	color:     image.NewColor(1, 1, 1),
	ambient:   0.1,
	diffuse:   0.9,
	specular:  0.9,
	shininess: 200,
}

// NewMaterial returns a new Material object
func NewMaterial(color *image.Color, ambient, diffuse, specular, shininess float64) *Material {
	return &Material{
		color:     color,
		ambient:   ambient,
		diffuse:   diffuse,
		specular:  specular,
		shininess: shininess,
	}
}

// SetColor sets the material's color field
func (m *Material) SetColor(color *image.Color) {
	m.color = color
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
