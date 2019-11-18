package main

import (
	"io"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/sjberman/golang-ray-tracer/pkg/image"
	"github.com/sjberman/golang-ray-tracer/pkg/base"
)

type projectile struct {
	position *base.Tuple
	velocity *base.Tuple
}

type environment struct {
	gravity *base.Tuple
	wind    *base.Tuple
}

func tick(e *environment, p *projectile) *projectile {
	newPos, _ := p.position.Add(p.velocity)
	step, _ := p.velocity.Add(e.gravity)
	newVel, _ := step.Add(e.wind)
	return &projectile{
		position: newPos,
		velocity: newVel,
	}
}

// Test program
func main() {
	f1, _ := os.Create("perfFile")
	pprof.StartCPUProfile(f1)
	defer pprof.StopCPUProfile()

	p := &projectile{
		position: base.NewPoint(0, 1, 0),
		velocity: base.NewVector(1, 5, 0).Normalize().Multiply(11.25),
	}
	e := &environment{
		gravity: base.NewVector(0, -0.1, 0),
		wind:    base.NewVector(0, 0, 0),
	}

	canvas := image.NewCanvas(500, 550)
	for p.position.GetY() > 0 {
		p = tick(e, p)
		x := int(p.position.GetX())
		y := int(900 - p.position.GetY())
		canvas.WritePixel(x, y, image.NewColor(1, 0, 0))
	}
	ppm := canvas.ToPPM()
	f, _ := os.Create("image.ppm")
	defer f.Close()

	io.Copy(f, strings.NewReader(ppm))
}
