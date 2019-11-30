package image

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestImage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Image Suite")
}

var _ = Describe("canvas tests", func() {
	It("creates a new canvas", func() {
		c := NewCanvas(10, 20)
		Expect(c.width).To(Equal(10))
		Expect(c.height).To(Equal(20))
		Expect(len(c.pixels)).To(Equal(10))
		for _, column := range c.pixels {
			Expect(len(column)).To(Equal(20))
			for _, row := range column {
				Expect(row).To(Equal(Black))
			}
		}
	})

	It("sets the color of a pixel", func() {
		c := NewCanvas(10, 10)
		color := NewColor(1, 0.5, 0.2)
		c.WritePixel(2, 3, color)
		Expect(c.PixelAt(3, 2)).To(Equal(&Black))
		Expect(c.PixelAt(2, 3)).To(Equal(color))
		Expect(c.PixelAt(20, 20)).To(Equal(&Black))
	})

	It("builds a ppm string", func() {
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
		Expect(c.toPPM()).To(Equal(expPPM))

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
		Expect(c.toPPM()).To(Equal(expPPM))
	})
})
