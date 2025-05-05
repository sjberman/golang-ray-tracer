package parser

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/object"
)

func setup(g *WithT, contents string) *Parser {
	file, err := os.CreateTemp(".", "testOBJ")
	g.Expect(err).ToNot(HaveOccurred())
	defer os.Remove(file.Name())

	_, err = file.WriteString(contents)
	g.Expect(err).ToNot(HaveOccurred())
	parser, err := Parse(file.Name())
	g.Expect(err).ToNot(HaveOccurred())

	return parser
}

func TestParse_IgnoreUnrecognizedLines(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	parser := setup(g, `There was a young lady named Bright
who traveled much faster than light.
She set out one day
in a relative way,
and came back the previous night`)
	g.Expect(parser.ignoredLines).To(Equal(5))
}

func TestParse_RecordVertices(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	parser := setup(g, `v -1 1 0
v -1.0000 0.5000 0.0000
v 1 0 0
v 1 1 0`)
	g.Expect(len(parser.vertices)).To(Equal(4))
	g.Expect(parser.vertices[1]).To(Equal(base.NewPoint(-1, 1, 0)))
	g.Expect(parser.vertices[2]).To(Equal(base.NewPoint(-1, 0.5, 0)))
	g.Expect(parser.vertices[3]).To(Equal(base.NewPoint(1, 0, 0)))
	g.Expect(parser.vertices[4]).To(Equal(base.NewPoint(1, 1, 0)))
}

func TestParse_RecordVertexNormals(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	parser := setup(g, `vn 0 0 1
vn 0.707 0 -0.707
vn 1 2 3`)
	g.Expect(len(parser.normals)).To(Equal(3))
	g.Expect(parser.normals[1]).To(Equal(base.NewVector(0, 0, 1)))
	g.Expect(parser.normals[2]).To(Equal(base.NewVector(0.707, 0, -0.707)))
	g.Expect(parser.normals[3]).To(Equal(base.NewVector(1, 2, 3)))
}

func TestParse_TriangleFaces(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	parser := setup(g, `v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0

f 1 2 3
f 1 3 4`)
	t1, ok := parser.groups[0].Objects[0].(*object.Triangle)
	g.Expect(ok).To(BeTrue())
	t2, ok := parser.groups[0].Objects[1].(*object.Triangle)
	g.Expect(ok).To(BeTrue())
	g.Expect(t1.P1).To(Equal(parser.vertices[1]))
	g.Expect(t1.P2).To(Equal(parser.vertices[2]))
	g.Expect(t1.P3).To(Equal(parser.vertices[3]))
	g.Expect(t2.P1).To(Equal(parser.vertices[1]))
	g.Expect(t2.P2).To(Equal(parser.vertices[3]))
	g.Expect(t2.P3).To(Equal(parser.vertices[4]))

	// with normals
	parser = setup(g, `v 0 1 0
v -1 0 0
v 1 0 0

vn -1 0 0
vn 1 0 0
vn 0 1 0

f 1//3 2//1 3//2
f 1/0/3 2/102/1 3/14/2`)
	st1, ok := parser.groups[0].Objects[0].(*object.SmoothTriangle)
	g.Expect(ok).To(BeTrue())
	st2, ok := parser.groups[0].Objects[1].(*object.SmoothTriangle)
	g.Expect(ok).To(BeTrue())
	g.Expect(st1.P1).To(Equal(parser.vertices[1]))
	g.Expect(st1.P2).To(Equal(parser.vertices[2]))
	g.Expect(st1.P3).To(Equal(parser.vertices[3]))
	g.Expect(st1.N1).To(Equal(parser.normals[3]))
	g.Expect(st1.N2).To(Equal(parser.normals[1]))
	g.Expect(st1.N3).To(Equal(parser.normals[2]))

	g.Expect(st2).To(Equal(st1))
}

func TestParse_TriangulatePolygons(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	parser := setup(g, `v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0
v 0 2 0

f 1 2 3 4 5`)
	t1, ok := parser.groups[0].Objects[0].(*object.Triangle)
	g.Expect(ok).To(BeTrue())
	t2, ok := parser.groups[0].Objects[1].(*object.Triangle)
	g.Expect(ok).To(BeTrue())
	t3, ok := parser.groups[0].Objects[2].(*object.Triangle)
	g.Expect(ok).To(BeTrue())
	g.Expect(t1.P1).To(Equal(parser.vertices[1]))
	g.Expect(t1.P2).To(Equal(parser.vertices[2]))
	g.Expect(t1.P3).To(Equal(parser.vertices[3]))
	g.Expect(t2.P1).To(Equal(parser.vertices[1]))
	g.Expect(t2.P2).To(Equal(parser.vertices[3]))
	g.Expect(t2.P3).To(Equal(parser.vertices[4]))
	g.Expect(t3.P1).To(Equal(parser.vertices[1]))
	g.Expect(t3.P2).To(Equal(parser.vertices[4]))
	g.Expect(t3.P3).To(Equal(parser.vertices[5]))
}

func TestParse_TrianglesInGroups(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	parser := setup(g, `v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0

g FirstGroup
f 1 2 3
g SecondGroup
f 1 3 4`)

	g.Expect(parser.groups).To(HaveLen(2))
	t1, ok := parser.groups[0].Objects[0].(*object.Triangle)
	g.Expect(ok).To(BeTrue())
	t2, ok := parser.groups[1].Objects[0].(*object.Triangle)
	g.Expect(ok).To(BeTrue())
	g.Expect(t1.P1).To(Equal(parser.vertices[1]))
	g.Expect(t1.P2).To(Equal(parser.vertices[2]))
	g.Expect(t1.P3).To(Equal(parser.vertices[3]))
	g.Expect(t2.P1).To(Equal(parser.vertices[1]))
	g.Expect(t2.P2).To(Equal(parser.vertices[3]))
	g.Expect(t2.P3).To(Equal(parser.vertices[4]))

	g.Expect(len(parser.GetGroup().Objects)).To(Equal(2))
}
