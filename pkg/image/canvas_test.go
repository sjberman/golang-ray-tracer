package image

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestNewCanvas(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	c := NewCanvas(10, 20)
	g.Expect(c.width).To(Equal(10))
	g.Expect(c.height).To(Equal(20))
	g.Expect(len(c.pixels)).To(Equal(10))
	for _, column := range c.pixels {
		g.Expect(len(column)).To(Equal(20))
		for _, row := range column {
			g.Expect(row).To(Equal(*Black))
		}
	}
}

func TestWritePixel(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	c := NewCanvas(10, 10)
	color := NewColor(1, 0.5, 0.2)
	c.WritePixel(2, 3, color)
	g.Expect(c.PixelAt(3, 2)).To(Equal(Black))
	g.Expect(c.PixelAt(2, 3)).To(Equal(color))
	g.Expect(c.PixelAt(20, 20)).To(Equal(Black))
}

func TestToPPM(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	c := NewCanvas(5, 3)
	c.WritePixel(0, 0, NewColor(1.5, 0, 0))
	c.WritePixel(2, 1, NewColor(0, 0.5, 0))
	c.WritePixel(4, 2, NewColor(-0.5, 0, 1))

	expPPM := `P3
5 3
255
255 0 0 0 0 0 0 0 0 0 0 0 0 0 0
0 0 0 0 0 0 0 128 0 0 0 0 0 0 0
0 0 0 0 0 0 0 0 0 0 0 0 0 0 255
`
	g.Expect(c.toPPM()).To(Equal(expPPM))

	// Verify we don't have lines over 70 characters
	c = NewCanvas(10, 2)
	for i, columns := range c.pixels {
		for j := range columns {
			c.WritePixel(i, j, NewColor(1, 0.8, 0.6))
		}
	}

	expPPM = `P3
10 2
255
255 204 153 255 204 153 255 204 153 255 204 153 255 204 153
255 204 153 255 204 153 255 204 153 255 204 153 255 204 153
255 204 153 255 204 153 255 204 153 255 204 153 255 204 153
255 204 153 255 204 153 255 204 153 255 204 153 255 204 153
`
	g.Expect(c.toPPM()).To(Equal(expPPM))
}
