package main

import (
	"fmt"
	"math"
	"time"

	. "github.com/sjberman/golang-ray-tracer/pkg/base"
	. "github.com/sjberman/golang-ray-tracer/pkg/image"
	. "github.com/sjberman/golang-ray-tracer/pkg/scene"
	. "github.com/sjberman/golang-ray-tracer/pkg/scene/object"
)

// Test program
func main() {
	startTime := time.Now()
	// f1, _ := os.Create("perfFile")
	// pprof.StartCPUProfile(f1)
	// defer pprof.StopCPUProfile()

	// Floor
	cp := NewCheckerPattern(Black, White)
	floor := NewPlane()
	floor.SetTransform(
		RotateY(math.Sqrt(2) / 2),
	)
	floor.Pattern = cp
	floor.Reflective = 0.4
	floor.Specular = 0

	c1 := NewCube()
	c1.Color = NewColor(1, 0, 0)
	c1.SetTransform(
		Translate(1, 1, 2),
		RotateY(math.Pi/4),
	)
	c1.Reflective = 0.8

	// Backdrop
	// rp := NewRingPattern(NewColor(1, 0, 0), White)
	// backdrop := NewPlane()
	// backdrop.SetTransform(
	// 	Translate(0, 0, 5),
	// 	RotateX(math.Pi/2),
	// )
	// backdrop.SetPattern(rp)

	// s1 object
	// s1 := NewSphere()
	// s1.SetTransform(
	// 	Translate(1, 1, 2),
	// )
	// sColor := NewColor(1, 0, 0.5)
	// s1.SetColor(sColor)
	// s1.SetSpecular(0.5)
	// s1.SetShininess(5)

	// s2 object
	// s2 := NewSphere()
	// s2.SetMaterial(s1.GetMaterial())
	// s2.SetTransform(
	// 	Translate(3, 1, 2),
	// )

	blueSphere := GlassSphere()
	blueSphere.Color = NewColor(0, 0, 0.2)
	blueSphere.SetTransform(
		Translate(-2, 0.7, 1.4),
		Scale(0.7, 0.7, 0.7),
	)

	greenSphere := NewSphere()
	greenSphere.Transparency = 0.6
	greenSphere.Color = NewColor(0, 0.2, 0)
	greenSphere.SetTransform(
		Translate(1, 2.8, -3),
		Scale(0.5, 0.5, 0.5),
	)

	// World
	light := NewPointLight(NewPoint(-5, 5, -1), White)
	world := NewWorld(light, []Object{
		floor,
		c1,
		// backdrop,
		// s1,
		// s2,
		blueSphere,
		greenSphere,
	})

	// Camera
	camera := NewCamera(300, 300, math.Pi/3)
	camera.SetTransform(ViewTransform(NewPoint(0, 4, -7), NewPoint(1, 0, 0), NewVector(0, 1, 0)))

	// Canvas
	canvas := Render(camera, world)

	err := canvas.WriteToFile("image.ppm")
	if err != nil {
		fmt.Println("error writing file: ", err.Error())
	}

	fmt.Println(time.Since(startTime).Round(time.Second))
}
