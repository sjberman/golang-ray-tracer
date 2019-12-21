package utils

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils Suite")
}

var _ = Describe("helpers tests", func() {
	It("returns the proper mod value", func() {
		Expect(Mod(4, 2)).To(Equal(0.0))
		Expect(Mod(-5, 2)).To(Equal(1.0))
		Expect(Mod(5, -2)).To(Equal(-1.0))
	})

	It("returns the maximum value in a list", func() {
		Expect(Max(1.0, 4.0, -3.4, 5.8, 5.3)).To(Equal(5.8))
	})

	It("returns the minimum value in a list", func() {
		Expect(Min(1.0, 4.0, -3.4, 5.8, -3)).To(Equal(-3.4))
	})
})
