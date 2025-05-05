package image

import (
	"fmt"
	"math"
	"os"
	"strings"
)

// Canvas represents a grid of pixels for displaying an image.
type Canvas struct {
	width  int
	height int
	// two dimensional array (matrix) of Colors
	pixels [][]Color
}

// NewCanvas returns a new Canvas object.
func NewCanvas(width, height int) *Canvas {
	pixels := make([][]Color, width)
	for i := range pixels {
		pixels[i] = make([]Color, height)
	}

	return &Canvas{
		width:  width,
		height: height,
		pixels: pixels,
	}
}

// WritePixel sets a Canvas's pixel to a color.
func (c *Canvas) WritePixel(x, y int, color *Color) {
	if (x <= c.width-1) && (y <= c.height-1) {
		c.pixels[x][y] = *color
	}
}

// PixelAt returns the Color of a Canvas's pixel.
func (c *Canvas) PixelAt(x, y int) *Color {
	if (x <= c.width-1) && (y <= c.height-1) {
		return &c.pixels[x][y]
	}

	return Black
}

// returns a PPM (portable pixelmap) string of the canvas.
func (c *Canvas) toPPM() string {
	header := fmt.Sprintf("P3\n%d %d\n%d\n", c.width, c.height, 255)
	var body strings.Builder
	for i := 0; i < c.height; i++ {
		var line strings.Builder
		for j := 0; j < c.width; j++ {
			color := c.PixelAt(j, i)
			red, green, blue := scalePixel(color)
			pixelVal := fmt.Sprintf("%d %d %d ", red, green, blue)
			// lines should not exceed 70 chars
			lineMax := 70
			if len(line.String()+pixelVal) > lineMax {
				body.WriteString(strings.TrimRight(line.String(), " ") + "\n")
				line.Reset()
			}
			line.WriteString(pixelVal)
		}
		body.WriteString(strings.TrimRight(line.String(), " ") + "\n")
	}
	ppm := header + body.String()

	return ppm
}

// WriteToFile writes the canvas ppm to a file.
func (c *Canvas) WriteToFile(name string) error {
	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(c.toPPM()); err != nil {
		return fmt.Errorf("error writing PPM string: %w", err)
	}

	return nil
}

// scales a color's values to be from 0 to 255.
func scalePixel(color *Color) (int64, int64, int64) {
	return scaleColor(color.red), scaleColor(color.green), scaleColor(color.blue)
}

// scales a color float value to be between 0 and 255.
func scaleColor(color float64) int64 {
	if color < 0 {
		return 0
	}
	if color > 1 {
		return 255
	}

	return int64(math.Round(255 * color))
}
