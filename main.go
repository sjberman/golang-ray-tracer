package main

import (
	"fmt"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
	"github.com/sjberman/golang-ray-tracer/pkg/ray"
)

// Test program
func main() {
	// f1, _ := os.Create("perfFile")
	// pprof.StartCPUProfile(f1)
	// defer pprof.StopCPUProfile()

	canvas := image.NewCanvas(100, 100)
	sphere := ray.NewSphere()
	sphere.SetTransform(base.ScalingMatrix(1, .5, 1))
	rayOrigin := base.NewPoint(0, 0, -5)
	wallZ := 10.0
	wallSize := 7.0
	pixelSize := wallSize / 100.0
	half := wallSize / 2
	red := image.NewColor(1, 0, 0)

	for y := 0; y < 99; y++ {
		worldY := half - pixelSize*float64(y)
		for x := 0; x < 99; x++ {
			worldX := -half + pixelSize*float64(x)
			position := base.NewPoint(worldX, worldY, wallZ)
			sub, _ := position.Subtract(rayOrigin)
			r := ray.NewRay(rayOrigin, sub.Normalize())
			ints := r.Intersect(sphere)
			if ray.Hit(ints) != nil {
				canvas.WritePixel(x, y, red)
			}
		}
	}

	err := canvas.WriteToFile("image.ppm")
	if err != nil {
		fmt.Println("error writing file: ", err.Error())
	}
}
