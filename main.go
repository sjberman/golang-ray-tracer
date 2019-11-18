package main

import (
	"io"
	"os"
	"strings"

	"github.com/sjberman/golang-ray-tracer/pkg/types"
)

type projectile struct {
	position *types.Tuple
	velocity *types.Tuple
}

type environment struct {
	gravity *types.Tuple
	wind    *types.Tuple
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
	p := &projectile{
		position: types.NewPoint(0, 1, 0),
		velocity: types.NewVector(1, 5, 0).Normalize().Multiply(11.25),
	}
	e := &environment{
		gravity: types.NewVector(0, -0.1, 0),
		wind:    types.NewVector(0, 0, 0),
	}

	canvas := types.NewCanvas(500, 550)
	for p.position.GetY() > 0 {
		p = tick(e, p)
		x := int(p.position.GetX())
		y := int(900 - p.position.GetY())
		canvas.WritePixel(x, y, types.NewColor(1, 0, 0))
	}
	ppm := canvas.ToPPM()
	f, _ := os.Create("image.ppm")
	defer f.Close()

	io.Copy(f, strings.NewReader(ppm))
}
