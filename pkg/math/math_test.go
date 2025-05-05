package math_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/math"
)

func TestMod(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	g.Expect(math.Mod(4, 2)).To(Equal(0.0))
	g.Expect(math.Mod(-5, 2)).To(Equal(1.0))
	g.Expect(math.Mod(5, -2)).To(Equal(-1.0))
}

func TestMax(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	g.Expect(math.Max(1.0, 4.0, -3.4, 5.8, 5.3)).To(Equal(5.8))
}

func TestMin(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	g.Expect(math.Min(1.0, 4.0, -3.4, 5.8, -3)).To(Equal(-3.4))
}
