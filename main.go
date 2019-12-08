package main

import (
	"fmt"
	"math"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
	"github.com/sjberman/golang-ray-tracer/pkg/scene"
)

// Test program
func main() {
	// f1, _ := os.Create("perfFile")
	// pprof.StartCPUProfile(f1)
	// defer pprof.StopCPUProfile()

	floor := scene.NewPlane()
	floor.SetTransform(base.ScalingMatrix(10, 0.01, 10))
	floor.GetMaterial().SetColor(image.NewColor(1, 0.9, 0.9))
	floor.GetMaterial().SetSpecular(0)

	// backdrop := scene.NewPlane()
	// transform := base.TranslationMatrix(0, 0, 4).Multiply(
	// 	base.XRotationMatrix(math.Pi / 2))
	// backdrop.SetTransform(transform)
	// backdrop.SetMaterial(floor.GetMaterial())

	middle := scene.NewSphere()
	middle.SetTransform(base.TranslationMatrix(-0.5, 1, 0.5))
	middle.GetMaterial().SetColor(image.NewColor(0.1, 1, 0.5))
	middle.GetMaterial().SetDiffuse(0.7)
	middle.GetMaterial().SetSpecular(0.3)

	right := scene.NewSphere()
	right.SetTransform(base.TranslationMatrix(1.5, 0.5, -0.5).Multiply(base.ScalingMatrix(0.5, 0.5, 0.5)))
	right.GetMaterial().SetColor(image.NewColor(0.5, 1, 0.1))
	right.GetMaterial().SetDiffuse(0.7)
	right.GetMaterial().SetSpecular(0.3)

	left := scene.NewSphere()
	left.SetTransform(base.TranslationMatrix(-1.5, 0.33, -0.75).Multiply(base.ScalingMatrix(0.33, 0.33, 0.33)))
	left.GetMaterial().SetColor(image.NewColor(1, 0.8, 0.1))
	left.GetMaterial().SetDiffuse(0.7)
	left.GetMaterial().SetSpecular(0.3)

	light := scene.NewPointLight(base.NewPoint(-10, 10, -10), image.NewColor(1, 1, 1))
	world := scene.NewWorld(light, []scene.Object{floor, middle, right, left})
	camera := scene.NewCamera(300, 300, math.Pi/3)

	camera.SetTransform(base.ViewTransform(base.NewPoint(0, 1.5, -5), base.NewPoint(0, 1, 0), base.NewVector(0, 1, 0)))

	canvas := scene.Render(camera, world)

	err := canvas.WriteToFile("image.ppm")
	if err != nil {
		fmt.Println("error writing file: ", err.Error())
	}
}
