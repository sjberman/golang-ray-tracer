package parser

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/object"
)

func TestParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parser Suite")
}

func setup(contents string) *Parser {
	file, err := ioutil.TempFile(".", "testOBJ")
	Expect(err).ToNot(HaveOccurred())
	defer os.Remove(file.Name())

	_, err = file.WriteString(contents)
	Expect(err).ToNot(HaveOccurred())
	parser, err := Parse(file.Name())
	Expect(err).ToNot(HaveOccurred())

	return parser
}

var _ = Describe("parser tests", func() {
	It("ignores unrecognized lines", func() {
		parser := setup(`There was a young lady named Bright​
who traveled much faster than light.​
She set out one day​
in a relative way,​
and came back the previous night.​`)
		Expect(parser.ignoredLines).To(Equal(5))
	})

	It("records vertices", func() {
		parser := setup(`v -1 1 0
v -1.0000 0.5000 0.0000
v 1 0 0
v 1 1 0`)
		Expect(len(parser.vertices)).To(Equal(4))
		Expect(parser.vertices[1]).To(Equal(base.NewPoint(-1, 1, 0)))
		Expect(parser.vertices[2]).To(Equal(base.NewPoint(-1, 0.5, 0)))
		Expect(parser.vertices[3]).To(Equal(base.NewPoint(1, 0, 0)))
		Expect(parser.vertices[4]).To(Equal(base.NewPoint(1, 1, 0)))
	})

	It("records vertex normals", func() {
		parser := setup(`vn 0 0 1
vn 0.707 0 -0.707
vn 1 2 3`)
		Expect(len(parser.normals)).To(Equal(3))
		Expect(parser.normals[1]).To(Equal(base.NewVector(0, 0, 1)))
		Expect(parser.normals[2]).To(Equal(base.NewVector(0.707, 0, -0.707)))
		Expect(parser.normals[3]).To(Equal(base.NewVector(1, 2, 3)))
	})

	It("parses triangle faces", func() {
		parser := setup(`v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0

f 1 2 3
f 1 3 4`)
		t1, ok := parser.groups[0].Objects[0].(*object.Triangle)
		Expect(ok).To(BeTrue())
		t2, ok := parser.groups[0].Objects[1].(*object.Triangle)
		Expect(ok).To(BeTrue())
		Expect(t1.P1).To(Equal(parser.vertices[1]))
		Expect(t1.P2).To(Equal(parser.vertices[2]))
		Expect(t1.P3).To(Equal(parser.vertices[3]))
		Expect(t2.P1).To(Equal(parser.vertices[1]))
		Expect(t2.P2).To(Equal(parser.vertices[3]))
		Expect(t2.P3).To(Equal(parser.vertices[4]))

		// with normals
		parser = setup(`v 0 1 0
v -1 0 0
v 1 0 0

vn -1 0 0
vn 1 0 0
vn 0 1 0

f 1//3 2//1 3//2
f 1/0/3 2/102/1 3/14/2`)
		st1, ok := parser.groups[0].Objects[0].(*object.SmoothTriangle)
		Expect(ok).To(BeTrue())
		st2, ok := parser.groups[0].Objects[1].(*object.SmoothTriangle)
		Expect(ok).To(BeTrue())
		Expect(st1.P1).To(Equal(parser.vertices[1]))
		Expect(st1.P2).To(Equal(parser.vertices[2]))
		Expect(st1.P3).To(Equal(parser.vertices[3]))
		Expect(st1.N1).To(Equal(parser.normals[3]))
		Expect(st1.N2).To(Equal(parser.normals[1]))
		Expect(st1.N3).To(Equal(parser.normals[2]))

		Expect(st2).To(Equal(st1))
	})

	It("triangulates polygons", func() {
		parser := setup(`v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0
v 0 2 0

f 1 2 3 4 5`)
		t1, ok := parser.groups[0].Objects[0].(*object.Triangle)
		Expect(ok).To(BeTrue())
		t2, ok := parser.groups[0].Objects[1].(*object.Triangle)
		Expect(ok).To(BeTrue())
		t3, ok := parser.groups[0].Objects[2].(*object.Triangle)
		Expect(ok).To(BeTrue())
		Expect(t1.P1).To(Equal(parser.vertices[1]))
		Expect(t1.P2).To(Equal(parser.vertices[2]))
		Expect(t1.P3).To(Equal(parser.vertices[3]))
		Expect(t2.P1).To(Equal(parser.vertices[1]))
		Expect(t2.P2).To(Equal(parser.vertices[3]))
		Expect(t2.P3).To(Equal(parser.vertices[4]))
		Expect(t3.P1).To(Equal(parser.vertices[1]))
		Expect(t3.P2).To(Equal(parser.vertices[4]))
		Expect(t3.P3).To(Equal(parser.vertices[5]))
	})

	It("parses triangles in groups", func() {
		parser := setup(`v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0

g FirstGroup
f 1 2 3
g SecondGroup
f 1 3 4`)
		t1, ok := parser.groups[0].Objects[0].(*object.Triangle)
		Expect(ok).To(BeTrue())
		t2, ok := parser.groups[1].Objects[0].(*object.Triangle)
		Expect(ok).To(BeTrue())
		Expect(t1.P1).To(Equal(parser.vertices[1]))
		Expect(t1.P2).To(Equal(parser.vertices[2]))
		Expect(t1.P3).To(Equal(parser.vertices[3]))
		Expect(t2.P1).To(Equal(parser.vertices[1]))
		Expect(t2.P2).To(Equal(parser.vertices[3]))
		Expect(t2.P3).To(Equal(parser.vertices[4]))

		Expect(len(parser.GetGroup().Objects)).To(Equal(2))
	})
})
