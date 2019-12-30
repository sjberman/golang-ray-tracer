package parser

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/object"
)

// Parser is an OBJ file reader
type Parser struct {
	ignoredLines int
	vertices     map[int]*base.Tuple
	normals      map[int]*base.Tuple
	groups       []object.Group
}

// Parse reads and parses an OBJ file and returns a Parser
func Parse(filename string) (*Parser, error) {
	p := &Parser{
		ignoredLines: 0,
		vertices:     make(map[int]*base.Tuple),
		normals:      make(map[int]*base.Tuple),
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	re := regexp.MustCompile(`[a-z]\s{1,}([-+]?[0-9.]+\s?)+|^g (?i)[a-z]+$`)
	vIndex := 1
	vnIndex := 1
	triangles := make(map[int][]object.Object)
	// smoothTriangles := make(map[int][]object.Object)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !re.MatchString(line) {
			p.ignoredLines++
			continue
		}
		split := strings.Fields(line)
		switch split[0] {
		case "v":
			// vertex
			points := make([]float64, 0, 3)
			for i := 1; i < 4; i++ {
				p, err := strconv.ParseFloat(split[i], 64)
				if err != nil {
					return nil, err
				}
				points = append(points, p)
			}
			p.vertices[vIndex] = base.NewPoint(points[0], points[1], points[2])
			vIndex++
		case "vn":
			// vertex normal
			vals := make([]float64, 0, 3)
			for i := 1; i < 4; i++ {
				p, err := strconv.ParseFloat(split[i], 64)
				if err != nil {
					return nil, err
				}
				vals = append(vals, p)
			}
			p.normals[vnIndex] = base.NewVector(vals[0], vals[1], vals[2])
			vnIndex++
		case "g":
			// group
			p.groups = append(p.groups, *object.NewGroup())
		case "f":
			// face
			var smooth bool
			v1, n1, err := getFields(split, 1)
			if err != nil {
				return nil, err
			}
			if _, ok := p.normals[n1]; ok {
				smooth = true
			}
			for i := 2; i < len(split)-1; i++ {
				v2, n2, err := getFields(split, i)
				if err != nil {
					return nil, err
				}
				v3, n3, err := getFields(split, i+1)
				if err != nil {
					return nil, err
				}

				if len(p.groups) == 0 {
					p.groups = append(p.groups, *object.NewGroup())
				}

				if smooth {
					triangle := object.NewSmoothTriangle(
						p.vertices[v1], p.vertices[v2], p.vertices[v3],
						p.normals[n1], p.normals[n2], p.normals[n3])
					// add to most recent group
					grpIdx := len(p.groups) - 1
					triangles[grpIdx] = append(triangles[grpIdx], triangle)
				} else {
					triangle := object.NewTriangle(p.vertices[v1], p.vertices[v2], p.vertices[v3])
					// add to most recent group
					grpIdx := len(p.groups) - 1
					triangles[grpIdx] = append(triangles[grpIdx], triangle)
				}
			}
		}
	}
	for idx, grp := range triangles {
		p.groups[idx].Add(grp...)
	}
	return p, nil
}

// GetGroup returns the single Group instance of all objects in the OBJ file
func (p *Parser) GetGroup() *object.Group {
	group := object.NewGroup()
	for _, o := range p.groups {
		o := o
		group.Add(&o)
	}
	return group
}

func getFields(list []string, index int) (int, int, error) {
	var err error
	var nIdx int
	faceSplit := strings.Split(list[index], "/")
	vIdx, err := strconv.Atoi(faceSplit[0])
	if err != nil {
		return 0, 0, err
	}
	if len(faceSplit) > 1 {
		nIdx, err = strconv.Atoi(faceSplit[2])
	}
	return vIdx, nIdx, err
}
