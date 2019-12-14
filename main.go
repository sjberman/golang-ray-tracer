package main

import (
	"fmt"
	"math"
	"time"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
	"github.com/sjberman/golang-ray-tracer/pkg/scene"
)

// Test program
func main() {
	startTime := time.Now()
	// f1, _ := os.Create("perfFile")
	// pprof.StartCPUProfile(f1)
	// defer pprof.StopCPUProfile()

	floor := scene.NewPlane()
	floor.SetTransform(base.YRotationMatrix(math.Sqrt(2) / 2).Multiply(
		base.ScalingMatrix(2, 2, 2)))
	cp := image.NewCheckerPattern(image.Black, image.White)
	// cp.SetTransform(base.ScalingMatrix(0.25, 0.25, 0.25))
	floor.GetMaterial().SetPattern(cp)

	backdrop := scene.NewPlane()
	transform := base.TranslationMatrix(0, 0, 5).Multiply(
		base.XRotationMatrix(math.Pi / 2))
	backdrop.SetTransform(transform)
	rp := image.NewRingPattern(image.NewColor(1, 0, 0), image.White)
	// sp.SetTransform(base.ScalingMatrix(0.5, 0.5, 0.5))
	backdrop.GetMaterial().SetPattern(rp)

	middle := scene.NewSphere()
	middle.SetTransform(base.TranslationMatrix(1.5, 0, 0).Multiply(
		base.ScalingMatrix(0.5, 4.5, 0.5)))
	sp := image.NewStripePattern(image.NewColor(1, 0.5, 0), image.NewColor(1, 0.3, 0))
	sp.SetTransform(base.ScalingMatrix(0.05, 0.05, 0.05).Multiply(base.ZRotationMatrix(math.Pi / 2)))
	middle.GetMaterial().SetPattern(sp)
	middle.GetMaterial().SetAmbient(0.15)

	// right := scene.NewSphere()
	// right.SetTransform(base.TranslationMatrix(1.5, 0.5, -0.5).Multiply(base.ScalingMatrix(0.5, 0.5, 0.5)))
	// // right.GetMaterial().SetColor(image.NewColor(0.5, 1, 0.1))
	// right.GetMaterial().SetDiffuse(0.7)
	// right.GetMaterial().SetSpecular(0.3)
	// gp := image.NewGradientPattern(image.NewColor(0.5, 0.5, 0.5), image.NewColor(0.7, 0.8, 0.2))
	// // gp.SetTransform(base.ScalingMatrix(0.5, 0.5, 0.5))
	// right.GetMaterial().SetPattern(gp)

	// left := scene.NewSphere()
	// left.SetTransform(base.TranslationMatrix(-1.5, 0.33, -0.75).Multiply(base.ScalingMatrix(0.33, 0.33, 0.33)))
	// left.GetMaterial().SetColor(image.NewColor(1, 0.8, 0.1))
	// left.GetMaterial().SetDiffuse(0.7)
	// left.GetMaterial().SetSpecular(0.3)

	light := scene.NewPointLight(base.NewPoint(-3, 6, -8), image.White.Multiply(1.2))
	world := scene.NewWorld(light, []scene.Object{
		floor,
		backdrop,
		middle,
		//right,
		//left,
	})
	camera := scene.NewCamera(400, 400, math.Pi/3)

	camera.SetTransform(base.ViewTransform(base.NewPoint(-5, 2, -4), base.NewPoint(0, 2, 0), base.NewVector(0, 1, 0)))

	canvas := scene.Render(camera, world)

	err := canvas.WriteToFile("image.ppm")
	if err != nil {
		fmt.Println("error writing file: ", err.Error())
	}

	fmt.Println(time.Since(startTime))
}
