package main

import (
	"fmt"
	"math"
	"time"

	. "github.com/sjberman/golang-ray-tracer/pkg/base"
	. "github.com/sjberman/golang-ray-tracer/pkg/image"
	. "github.com/sjberman/golang-ray-tracer/pkg/scene"
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
	floor.SetPattern(cp)
	floor.SetReflective(0.4)
	floor.SetSpecular(0)

	// Backdrop
	// rp := NewRingPattern(NewColor(1, 0, 0), White)
	// backdrop := NewPlane()
	// backdrop.SetTransform(
	// 	Translate(0, 0, 5),
	// 	RotateX(math.Pi/2),
	// )
	// backdrop.SetPattern(rp)

	// s1 object
	s1 := NewSphere()
	s1.SetTransform(
		Translate(1, 1, 2),
	)
	sColor := NewColor(1, 0, 0.5)
	s1.SetColor(sColor)
	s1.SetSpecular(0.5)
	// s1.SetShininess(5)

	// s2 object
	s2 := NewSphere()
	s2.SetMaterial(s1.GetMaterial())
	s2.SetTransform(
		Translate(3, 1, 2),
	)

	blueSphere := GlassSphere()
	blueSphere.SetColor(NewColor(0, 0, 0.2))
	blueSphere.SetTransform(
		Translate(2.5, 0.7, -1),
		Scale(0.7, 0.7, 0.7),
	)

	// greenSphere := GlassSphere()
	// greenSphere.SetColor(NewColor(0, 0.2, 0))
	// greenSphere.SetTransform(
	// 	Translate(-3, 2, 4),
	// 	Scale(2, 2, 2),
	// )

	// World
	light := NewPointLight(NewPoint(-5, 5, -1), White)
	world := NewWorld(light, []Object{
		floor,
		// backdrop,
		s1,
		s2,
		blueSphere,
		// greenSphere,
	})

	// Camera
	camera := NewCamera(640, 400, math.Pi/3)
	camera.SetTransform(ViewTransform(NewPoint(0, 2, -7), NewPoint(1, 0, 0), NewVector(0, 1, 0)))

	// Canvas
	canvas := Render(camera, world)

	err := canvas.WriteToFile("image.ppm")
	if err != nil {
		fmt.Println("error writing file: ", err.Error())
	}

	fmt.Println(time.Since(startTime))
}
