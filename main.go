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
		RotateY(math.Sqrt(2)/2),
		Scale(2, 2, 2),
	)
	floor.SetPattern(cp)

	// Backdrop
	rp := NewRingPattern(NewColor(1, 0, 0), White)
	backdrop := NewPlane()
	backdrop.SetTransform(
		Translate(0, 0, 5),
		RotateX(math.Pi/2),
	)
	backdrop.SetPattern(rp)

	// Middle object
	sp := NewStripePattern(NewColor(1, 0.5, 0), NewColor(1, 0.3, 0))
	sp.SetTransform(
		Scale(0.05, 0.05, 0.05),
		RotateZ(math.Pi/2),
	)
	middle := NewSphere()
	middle.SetTransform(
		Translate(1.5, 0, 0),
		Scale(0.5, 4.5, 0.5),
	)
	middle.SetPattern(sp)
	middle.SetAmbient(0.15)

	// Right object
	right := NewSphere()
	right.SetTransform(
		Translate(0.2, 1, -1.5),
		Scale(0.3, 0.3, 0.3),
	)
	right.SetColor(NewColor(0, 1, 0))

	// Left object
	left := NewSphere()
	left.SetTransform(
		Translate(-1.5, 2, 2.5),
		Scale(0.7, 0.7, 0.7),
	)
	left.SetColor(NewColor(1, 0.8, 0.1))

	// World
	light := NewPointLight(NewPoint(-3, 6, -8), White.Multiply(1.2))
	world := NewWorld(light, []Object{
		floor,
		backdrop,
		middle,
		right,
		left,
	})

	// Camera
	camera := NewCamera(300, 300, math.Pi/3)
	camera.SetTransform(ViewTransform(NewPoint(-5, 2, -4), NewPoint(0, 2, 0), NewVector(0, 1, 0)))

	// Canvas
	canvas := Render(camera, world)

	err := canvas.WriteToFile("image.ppm")
	if err != nil {
		fmt.Println("error writing file: ", err.Error())
	}

	fmt.Println(time.Since(startTime))
}
