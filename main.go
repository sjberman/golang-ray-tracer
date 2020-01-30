package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"runtime/pprof"
	"time"

	"github.com/sjberman/golang-ray-tracer/demo"
	. "github.com/sjberman/golang-ray-tracer/pkg/base"
	. "github.com/sjberman/golang-ray-tracer/pkg/image"
	. "github.com/sjberman/golang-ray-tracer/pkg/parser"
	. "github.com/sjberman/golang-ray-tracer/pkg/scene"
	. "github.com/sjberman/golang-ray-tracer/pkg/scene/object"
)

var objFile = flag.String("file", "", "OBJ file to render")

func usingObjFile(objfile *string, camera *Camera) (*Group, []*PointLight) {
	var group *Group
	parser, err := Parse(*objFile)
	if err != nil {
		fmt.Println("error parsing OBJ file: ", err)
		os.Exit(1)
	}
	group = parser.GetGroup()
	fmt.Println("OBJ minimum: ", group.Bounds().Minimum)
	fmt.Println("OBJ maximum: ", group.Bounds().Maximum)
	camera.SetTransform(ViewTransform(NewPoint(0, 2.5, -10), NewPoint(0, 1, 0), NewVector(0, 1, 0)))
	lights := []*PointLight{
		NewPointLight(NewPoint(-10, 100, -100), White),
	}
	return group, lights
}

func usingCustomObjects(camera *Camera) (Object, []*PointLight) {
	// return demo.FiveBallRainbow(camera)
	return demo.AllShapes(camera)
	// return demo.CSG(camera)
}

// Test program
func main() {
	startTime := time.Now()
	f1, _ := os.Create("perfFile")
	pprof.StartCPUProfile(f1)
	defer pprof.StopCPUProfile()

	flag.Parse()

	camera := NewCamera(1200, 1200, math.Pi/3)

	cp := NewCheckerPattern(Black, White)
	cp.SetTransform(RotateY(math.Sqrt(2) / 2))
	floor := NewPlane()
	floor.SetTransform(
		Translate(0, -1, 0),
	)
	floor.Pattern = cp
	floor.Reflective = 0.2

	wall := NewPlane()
	rp := NewRingPattern(Black, White)
	wall.SetTransform(
		RotateX(math.Pi/2),
		Translate(0, 5, 0),
	)
	wall.Pattern = rp

	// Build scene using OBJ file
	// group, lights := usingObjFile(objFile, camera)
	// group.Color = NewColor(0.5, 0.5, 0)
	// group.Shininess = 300
	// group.SetMaterial(group.GetMaterial())

	// Build scene using custom objects
	obj, lights := usingCustomObjects(camera)

	obj.Divide(1)

	world := NewWorld(lights, []Object{
		obj,
		floor,
		//wall,
	})

	canvas := Render(camera, world)
	err := canvas.WriteToFile("image.png")
	if err != nil {
		fmt.Println("error writing file: ", err.Error())
	}

	fmt.Println(time.Since(startTime).Round(time.Second))
}
