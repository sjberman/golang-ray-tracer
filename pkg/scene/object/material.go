package object

import "github.com/sjberman/golang-ray-tracer/pkg/image"

// Material contains the attributes of a surface material
type Material struct {
	Color           *image.Color
	Pattern         image.Pattern
	Ambient         float64
	Diffuse         float64
	Specular        float64
	Shininess       float64
	Reflective      float64
	Transparency    float64
	RefractiveIndex float64
	Shadow          bool
}

var DefaultMaterial = Material{
	Color:           image.White,
	Ambient:         0.1,
	Diffuse:         0.9,
	Specular:        0.9,
	Shininess:       200,
	Reflective:      0,
	Transparency:    0,
	RefractiveIndex: 1.0,
	Shadow:          true,
}
