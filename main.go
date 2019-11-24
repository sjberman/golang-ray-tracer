package main

import (
	"fmt"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
	"github.com/sjberman/golang-ray-tracer/pkg/scene"
)

// Test program
func main() {
	// f1, _ := os.Create("perfFile")
	// pprof.StartCPUProfile(f1)
	// defer pprof.StopCPUProfile()

	canvasSize := 300
	canvas := image.NewCanvas(canvasSize, canvasSize)
	sphere := scene.NewSphere()
	// sphere.SetTransform(base.ScalingMatrix(1, .5, 1))
	sphere.GetMaterial().SetColor(image.NewColor(1, 0.2, 1))
	rayOrigin := base.NewPoint(0, 0, -5)
	wallZ := 10.0
	wallSize := 7.0
	pixelSize := wallSize / float64(canvasSize)
	half := wallSize / 2

	light := scene.NewPointLight(base.NewPoint(0, 10, -10), image.NewColor(1, 1, 1))

	for y := 0; y < canvasSize-1; y++ {
		worldY := half - pixelSize*float64(y)
		for x := 0; x < canvasSize-1; x++ {
			worldX := -half + pixelSize*float64(x)
			position := base.NewPoint(worldX, worldY, wallZ)
			sub, _ := position.Subtract(rayOrigin)
			r := scene.NewRay(rayOrigin, sub.Normalize())
			ints := r.Intersect(sphere)
			hit := scene.Hit(ints)
			if hit != nil {
				point := r.Position(hit.GetValue())
				normal := hit.GetObject().NormalAt(point)
				eye := r.GetDirection().Negate()

				color := scene.Lighting(light, sphere.GetMaterial(), point, eye, normal)
				canvas.WritePixel(x, y, color)
			}
		}
	}

	err := canvas.WriteToFile("image.ppm")
	if err != nil {
		fmt.Println("error writing file: ", err.Error())
	}
}
