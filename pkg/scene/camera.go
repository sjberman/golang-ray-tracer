package scene

import (
	"math"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/ray"
)

// Camera is the viewpoint of a scene
type Camera struct {
	hsize       int // horizontal size
	vsize       int // vertical size
	fieldOfView float64
	pixelSize   float64
	halfWidth   float64
	halfHeight  float64
	transform   *base.Matrix
}

// NewCamera returns a new Camera object
func NewCamera(hsize, vsize int, fieldOfView float64) *Camera {
	c := &Camera{
		hsize:       hsize,
		vsize:       vsize,
		fieldOfView: fieldOfView,
		transform:   &base.Identity,
	}
	halfView := math.Tan(fieldOfView / 2)
	aspect := float64(hsize) / float64(vsize)
	if aspect >= 1 {
		c.halfWidth = halfView
		c.halfHeight = halfView / aspect
	} else {
		c.halfWidth = halfView * aspect
		c.halfHeight = halfView
	}
	c.pixelSize = (c.halfWidth * 2) / float64(c.hsize)
	return c
}

// SetTransform sets the transform matrix of the camera
func (c *Camera) SetTransform(matrix *base.Matrix) {
	c.transform = matrix
}

// RayForPixel returns a ray starting at the camera and going to x,y on the canvas
func (c *Camera) RayForPixel(x, y int) *ray.Ray {
	// the offset from the edge of the canvas to the pixel's center
	xOffset := (float64(x) + 0.5) * c.pixelSize
	yOffset := (float64(y) + 0.5) * c.pixelSize

	// the untransformed coordinates of the pixel in world space.
	// (camera looks towards -z, so +x is to the left)
	worldX := c.halfWidth - xOffset
	worldY := c.halfHeight - yOffset

	// using the camera matrix, transform the canvas point and the origin,
	// and then compute the ray's direction vector
	// (canvas is at z = -1)
	inverse := c.transform.Inverse()
	pixel := inverse.MultiplyTuple(base.NewPoint(worldX, worldY, -1))
	origin := inverse.MultiplyTuple(base.Origin)
	direction := pixel.Subtract(origin).Normalize()

	return ray.NewRay(origin, direction)
}
