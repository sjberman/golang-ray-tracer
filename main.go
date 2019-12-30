package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"runtime/pprof"
	"time"

	. "github.com/sjberman/golang-ray-tracer/pkg/base"
	. "github.com/sjberman/golang-ray-tracer/pkg/image"
	. "github.com/sjberman/golang-ray-tracer/pkg/parser"
	. "github.com/sjberman/golang-ray-tracer/pkg/scene"
	. "github.com/sjberman/golang-ray-tracer/pkg/scene/object"
)

var objFile = flag.String("file", "", "OBJ file to render")

func usingObjFile(objfile *string, camera *Camera) (*Group, *PointLight) {
	var group *Group
	parser, err := Parse(*objFile)
	if err != nil {
		fmt.Println("error parsing OBJ file: ", err)
		os.Exit(1)
	}
	group = parser.GetGroup()
	fmt.Println("OBJ minimum: ", group.Bounds().Minimum)
	fmt.Println("OBJ maximum: ", group.Bounds().Maximum)
	camera.SetTransform(ViewTransform(NewPoint(-4, 4, -10), NewPoint(0, 2.5, 0.5), NewVector(0, 1, 0)))
	return group, NewPointLight(NewPoint(0, 15, -350), White)
}

func usingCustomObjects(camera *Camera) (*Group, *PointLight) {
	r1 := NewSphere()
	r1.SetTransform(
		Translate(6, 1, 4),
	)
	r1.Color = NewColor(1, 0.3, 0.2)
	r1.Specular = 0.4
	r1.Shininess = 5

	r2 := NewSphere()
	r2.SetMaterial(r1.GetMaterial())
	r2.SetTransform(
		Translate(2, 1, 3),
	)

	r3 := NewSphere()
	r3.SetMaterial(r1.GetMaterial())
	r3.SetTransform(
		Translate(-1, 1, 2),
	)

	blueSphere := GlassSphere()
	blueSphere.Color = NewColor(0, 0, 0.2)
	blueSphere.Diffuse = 0.4
	blueSphere.Specular = 0.9
	blueSphere.Reflective = 0.9
	blueSphere.Shininess = 300
	blueSphere.SetTransform(
		Translate(0.6, 0.7, -0.6),
		Scale(0.7, 0.7, 0.7),
	)

	greenSphere := GlassSphere()
	greenSphere.SetMaterial(blueSphere.GetMaterial())
	greenSphere.Color = NewColor(0, 0.2, 0)
	greenSphere.SetTransform(
		Translate(-0.7, 0.5, -0.8),
		Scale(0.5, 0.5, 0.5),
	)
	group := NewGroup()
	group.Add(
		r1,
		r2,
		r3,
		blueSphere,
		greenSphere,
	)
	camera.SetTransform(ViewTransform(NewPoint(-2.6, 1.5, -3.9), NewPoint(-0.6, 1, -0.8), NewVector(0, 1, 0)))
	return group, NewPointLight(NewPoint(-5, 5, -1), White)
}

// Test program
func main() {
	startTime := time.Now()
	f1, _ := os.Create("perfFile")
	pprof.StartCPUProfile(f1)
	defer pprof.StopCPUProfile()

	flag.Parse()

	camera := NewCamera(1000, 800, math.Pi/3)

	cp := NewCheckerPattern(Black, White)
	cp.SetTransform(RotateY(math.Sqrt(2) / 2))
	floor := NewPlane()
	floor.Pattern = cp
	floor.Reflective = 0.4

	// Build scene using OBJ file
	// group, light := usingObjFile(objFile, camera)
	// group.Color = NewColor(0.5, 0.5, 0)
	// group.Shininess = 300
	// group.SetMaterial(group.GetMaterial())

	// Build scene using custom objects
	group, light := usingCustomObjects(camera)

	group.Divide(1)

	world := NewWorld(light, []Object{group, floor})

	canvas := Render(camera, world)
	err := canvas.WriteToFile("image.ppm")
	if err != nil {
		fmt.Println("error writing file: ", err.Error())
	}

	fmt.Println(time.Since(startTime).Round(time.Second))
}
