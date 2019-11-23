package main

import (
	"fmt"
	"math"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
)

// Test program
func main() {
	// f1, _ := os.Create("perfFile")
	// pprof.StartCPUProfile(f1)
	// defer pprof.StopCPUProfile()

	canvas := image.NewCanvas(500, 500)
	clockRadius := float64((3.0 / 8.0) * 500.0)
	originX, originY := 250, 250
	twelve := base.NewPoint(0, 1, 0)

	for i := 1.0; i < 13.0; i++ {
		rz := base.ZRotationMatrix(i * (math.Pi / 6))
		iTime := rz.MultiplyTuple(twelve)

		x := int(math.Round(iTime.GetX()*clockRadius)) + originX
		y := int(math.Round(iTime.GetY()*clockRadius)) + originY
		canvas.WritePixel(x, y, image.NewColor(1, 1, 1))
	}

	err := canvas.WriteToFile("image.ppm")
	if err != nil {
		fmt.Println("error writing file: ", err.Error())
	}
}
