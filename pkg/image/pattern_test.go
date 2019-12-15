package image

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
)

var _ = Describe("pattern tests", func() {
	testNewPattern := func(p Pattern) {
		Expect(p.GetColors()[0]).To(Equal(White))
		Expect(p.GetColors()[1]).To(Equal(Black))
		Expect(p.GetTransform()).To(Equal(&base.Identity))

		t := base.Translate(2, 3, 4)
		p.SetTransform(t)
		Expect(p.GetTransform()).To(Equal(t))
	}

	It("creates patterns", func() {
		p := NewPattern(White, Black, nil)
		testNewPattern(p)
	})

	Context("stripe patterns", func() {
		It("creates stripe patterns", func() {
			sp := NewStripePattern(White, Black)
			testNewPattern(sp)
		})

		It("returns the stripe color at a point", func() {
			sp := NewStripePattern(White, Black)

			// constant in y
			stripe := sp.PatternAt(base.NewPoint(0, 0, 0))
			Expect(stripe).To(Equal(White))

			stripe = sp.PatternAt(base.NewPoint(0, 1, 0))
			Expect(stripe).To(Equal(White))

			stripe = sp.PatternAt(base.NewPoint(0, 2, 0))
			Expect(stripe).To(Equal(White))

			// constant in z
			stripe = sp.PatternAt(base.NewPoint(0, 0, 0))
			Expect(stripe).To(Equal(White))

			stripe = sp.PatternAt(base.NewPoint(0, 0, 1))
			Expect(stripe).To(Equal(White))

			stripe = sp.PatternAt(base.NewPoint(0, 0, 2))
			Expect(stripe).To(Equal(White))

			// alternates in x
			stripe = sp.PatternAt(base.NewPoint(0, 0, 0))
			Expect(stripe).To(Equal(White))

			stripe = sp.PatternAt(base.NewPoint(0.9, 0, 0))
			Expect(stripe).To(Equal(White))

			stripe = sp.PatternAt(base.NewPoint(1, 0, 0))
			Expect(stripe).To(Equal(Black))

			stripe = sp.PatternAt(base.NewPoint(-0.1, 0, 0))
			Expect(stripe).To(Equal(Black))

			stripe = sp.PatternAt(base.NewPoint(-1, 0, 0))
			Expect(stripe).To(Equal(Black))

			stripe = sp.PatternAt(base.NewPoint(-1.1, 0, 0))
			Expect(stripe).To(Equal(White))
		})
	})

	Context("gradient patterns", func() {
		It("creates gradient patterns", func() {
			gp := NewGradientPattern(White, Black)
			testNewPattern(gp)
		})

		It("returns the gradient color at a point", func() {
			gp := NewGradientPattern(White, Black)

			gradient := gp.PatternAt(base.NewPoint(0, 0, 0))
			Expect(gradient).To(Equal(White))

			gradient = gp.PatternAt(base.NewPoint(0.25, 0, 0))
			Expect(gradient).To(Equal(NewColor(0.75, 0.75, 0.75)))

			gradient = gp.PatternAt(base.NewPoint(0.5, 0, 0))
			Expect(gradient).To(Equal(NewColor(0.5, 0.5, 0.5)))

			gradient = gp.PatternAt(base.NewPoint(0.75, 0, 0))
			Expect(gradient).To(Equal(NewColor(0.25, 0.25, 0.25)))
		})
	})

	Context("ring patterns", func() {
		It("creates ring patterns", func() {
			rp := NewRingPattern(White, Black)
			testNewPattern(rp)
		})

		It("returns the ring color at a point", func() {
			rp := NewRingPattern(White, Black)

			ring := rp.PatternAt(base.NewPoint(0, 0, 0))
			Expect(ring).To(Equal(White))

			ring = rp.PatternAt(base.NewPoint(1, 0, 0))
			Expect(ring).To(Equal(Black))

			ring = rp.PatternAt(base.NewPoint(0, 0, 1))
			Expect(ring).To(Equal(Black))

			// 0.708 is just slightly more than sqrt(2)/2
			ring = rp.PatternAt(base.NewPoint(0.708, 0, 0.708))
			Expect(ring).To(Equal(Black))
		})
	})

	Context("checker patterns", func() {
		It("creates checker patterns", func() {
			cp := NewCheckerPattern(White, Black)
			testNewPattern(cp)
		})

		It("returns the checker color at a point", func() {
			cp := NewCheckerPattern(White, Black)

			// should repeat in x
			checker := cp.PatternAt(base.NewPoint(0, 0, 0))
			Expect(checker).To(Equal(White))

			checker = cp.PatternAt(base.NewPoint(0.99, 0, 0))
			Expect(checker).To(Equal(White))

			checker = cp.PatternAt(base.NewPoint(1.01, 0, 0))
			Expect(checker).To(Equal(Black))

			// should repeat in y
			checker = cp.PatternAt(base.NewPoint(0, 0, 0))
			Expect(checker).To(Equal(White))

			checker = cp.PatternAt(base.NewPoint(0, 0.99, 0))
			Expect(checker).To(Equal(White))

			checker = cp.PatternAt(base.NewPoint(0, 1.01, 0))
			Expect(checker).To(Equal(Black))

			// should repeat in z
			checker = cp.PatternAt(base.NewPoint(0, 0, 0))
			Expect(checker).To(Equal(White))

			checker = cp.PatternAt(base.NewPoint(0, 0, 0.99))
			Expect(checker).To(Equal(White))

			checker = cp.PatternAt(base.NewPoint(0, 0, 1.01))
			Expect(checker).To(Equal(Black))
		})
	})
})
