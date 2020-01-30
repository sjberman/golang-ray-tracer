package demo

import (
	"math"

	. "github.com/sjberman/golang-ray-tracer/pkg/base"
	. "github.com/sjberman/golang-ray-tracer/pkg/image"
	. "github.com/sjberman/golang-ray-tracer/pkg/scene"
	. "github.com/sjberman/golang-ray-tracer/pkg/scene/object"
)

func FiveBallRainbow(camera *Camera) (Object, []*PointLight) {
	s1 := NewSphere()
	s1.Color = NewColor(1, 0, 0)
	s1.Reflective = 0.9

	s2 := NewSphere()
	s2.SetTransform(
		Scale(0.5, 0.5, 0.5),
		Translate(-3.5, 0, 0),
	)
	s2.Color = NewColor(0, 1, 1)
	s2.Reflective = 0.9

	s3 := NewSphere()
	s3.SetTransform(
		Scale(0.5, 0.5, 0.5),
		Translate(0, -3.5, 0),
	)
	s3.Color = NewColor(0.1, 0.4, 1)
	s3.Reflective = 0.9

	s4 := NewSphere()
	s4.SetTransform(
		Scale(0.5, 0.5, 0.5),
		Translate(3.5, 0, 0),
	)
	s4.Color = NewColor(1, 1, 0)
	s4.Reflective = 0.9

	s5 := NewSphere()
	s5.SetTransform(
		Scale(0.5, 0.5, 0.5),
		Translate(0, 3.5, 0),
	)
	s5.Color = NewColor(0.6, 1, 0)
	s5.Reflective = 0.9

	group := NewGroup()
	group.Add(
		s1,
		s2,
		s3,
		s4,
		s5,
	)
	camera.SetTransform(ViewTransform(NewPoint(0, 0, -5), NewPoint(0, 0, 0), NewVector(0, 1, 0)))
	lights := []*PointLight{
		NewPointLight(NewPoint(-6, 6, -5), White),
		NewPointLight(NewPoint(10, -10, 3), White),
	}
	return group, lights
}

func AllShapes(camera *Camera) (Object, []*PointLight) {
	s1 := NewSphere()
	s1.SetTransform(
		Translate(0, 2, 0),
	)
	s1.Color = NewColor(1, 0, 0)

	s2 := GlassSphere()
	s2.SetTransform(
		Translate(-2.5, 2, -1),
	)
	s2.Color = NewColor(0, 1, 0)

	s3 := NewSphere()
	s3.SetTransform(
		Translate(2.5, 2, -1),
	)
	s3.Color = NewColor(0, 0, 1)
	s3.Specular = 0

	c1 := NewCylinder()
	c1.Closed = true
	c1.Maximum = 1
	c1.Minimum = -1
	c1.SetTransform(
		Translate(-2.5, 0, -1),
	)
	c1.Color = NewColor(0.5, 0.2, 0.7)
	c1.Reflective = 0.1

	c2 := NewCube()
	c2.Color = NewColor(0.2, 0.7, 0.4)
	c2.SetTransform(
		RotateY(math.Pi / 6),
	)
	c2.Reflective = 0.2

	c3 := NewCone()
	c3.Closed = true
	c3.Maximum = 1
	c3.Minimum = -1
	c3.SetTransform(
		Translate(2.5, 0, -1),
	)
	c3.Color = NewColor(0.8, 0.3, 0.4)
	c3.Reflective = 0.5

	group := NewGroup()
	group.Add(
		s1,
		s2,
		s3,
		c1,
		c2,
		c3,
	)
	camera.SetTransform(ViewTransform(NewPoint(0, 5, -10), NewPoint(0, 0, 0), NewVector(0, 1, 0)))
	lights := []*PointLight{
		NewPointLight(NewPoint(-5, 10, -10), White),
	}
	return group, lights
}
