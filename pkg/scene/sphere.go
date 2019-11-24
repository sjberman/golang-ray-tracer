package scene

import "github.com/sjberman/golang-ray-tracer/pkg/base"

// Sphere is a sphere object for ray intersection
type Sphere struct {
	transform base.Matrix
	material  Material
}

// NewSphere returns a new Sphere object
func NewSphere() *Sphere {
	return &Sphere{
		transform: base.Identity,
		material:  defaultMaterial,
	}
}

// SetTransform sets the sphere's transform to the supplied matrix
func (s *Sphere) SetTransform(matrix *base.Matrix) {
	s.transform = *matrix
}

// SetMaterial sets the sphere's material
func (s *Sphere) SetMaterial(material *Material) {
	s.material = *material
}

// GetMaterial gets the sphere's material
func (s *Sphere) GetMaterial() *Material {
	return &s.material
}

// NormalAt returns the surface normal at a position on the sphere
func (s *Sphere) NormalAt(worldPoint *base.Tuple) *base.Tuple {
	// convert the point from world space to object space
	// (sphere is likely not at the world origin)
	inverse, _ := s.transform.Inverse()
	objectPoint := inverse.MultiplyTuple(worldPoint)

	objectNormal, _ := objectPoint.Subtract(base.Origin)
	// convert normal back to world space
	worldNormal := inverse.Transpose().MultiplyTuple(objectNormal)
	// ensure this is a vector
	worldNormal.SetW(0)
	return worldNormal.Normalize()
}
