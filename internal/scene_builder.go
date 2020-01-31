package internal

import (
	"fmt"
	"math"

	"github.com/sjberman/golang-ray-tracer/pkg/base"
	"github.com/sjberman/golang-ray-tracer/pkg/image"
	"github.com/sjberman/golang-ray-tracer/pkg/parser"
	"github.com/sjberman/golang-ray-tracer/pkg/scene"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/object"
	"github.com/sjberman/golang-ray-tracer/schema"
)

// CreateCamera builds a camera object using the spec
func CreateCamera(cam *schema.Camera) *scene.Camera {
	fov := cam.FieldOfView * math.Pi / 180
	camera := scene.NewCamera(cam.Width, cam.Height, fov)

	from := base.NewPoint(cam.From[0], cam.From[1], cam.From[2])
	to := base.NewPoint(cam.To[0], cam.To[1], cam.To[2])
	up := base.NewPoint(cam.Up[0], cam.Up[1], cam.Up[2])

	camera.SetTransform(base.ViewTransform(from, to, up))
	return camera
}

// CreateLights builds the light objects using the spec
func CreateLights(lights []*schema.Light) []*scene.PointLight {
	newLights := []*scene.PointLight{}
	for _, light := range lights {
		point := base.NewPoint(light.At[0], light.At[1], light.At[2])
		color := image.NewColor(light.Intensity[0], light.Intensity[1], light.Intensity[2])
		newLights = append(newLights, scene.NewPointLight(point, color))
	}
	return newLights
}

// CreateShapes builds the shape objects using the spec
func CreateShapes(shapes []*schema.Shape) ([]object.Object, map[string]*object.Object) {
	objs := []object.Object{}
	shapeMap := make(map[string]*object.Object)
	for _, shape := range shapes {
		var obj object.Object

		switch shape.Type {
		case "cone":
			cone := object.NewCone()
			if shape.Closed != nil {
				cone.Closed = *shape.Closed
			}
			if shape.Minimum != nil {
				cone.Minimum = *shape.Minimum
			}
			if shape.Maximum != nil {
				cone.Maximum = *shape.Maximum
			}
			obj = cone
		case "cube":
			cube := object.NewCube()
			obj = cube
		case "cylinder":
			cylinder := object.NewCylinder()
			if shape.Closed != nil {
				cylinder.Closed = *shape.Closed
			}
			if shape.Minimum != nil {
				cylinder.Minimum = *shape.Minimum
			}
			if shape.Maximum != nil {
				cylinder.Maximum = *shape.Maximum
			}
			obj = cylinder
		case "plane":
			plane := object.NewPlane()
			obj = plane
		case "sphere":
			sphere := object.NewSphere()
			obj = sphere
		case "glassSphere":
			sphere := object.GlassSphere()
			obj = sphere
		}

		obj.SetMaterial(getMaterial(shape.Material))
		obj.SetTransform(getTransforms(shape.Transform)...)
		objs = append(objs, obj)
		shapeMap[shape.Name] = &obj
	}
	return objs, shapeMap
}

// CreateGroupsAndCSGs builds the group and csg objects using the spec
func CreateGroupsAndCSGs(
	sceneStruct schema.RayTracerScene,
	shapes []object.Object,
	shapeMap map[string]*object.Object,
) []object.Object {
	groupMap := make(map[string]*object.Group)
	csgMap := make(map[string]*object.Csg)

	groups := make([]object.Object, 0, len(sceneStruct.Groups))
	csgs := make([]object.Object, 0, len(sceneStruct.Csgs))

	// First, populate the group map with default group
	for _, grp := range sceneStruct.Groups {
		groupMap[grp.Name] = object.NewGroup()
	}

	// Create CSGs
	for _, csg := range sceneStruct.Csgs {
		var left, right object.Object
		if o, ok := shapeMap[csg.LeftChild]; ok {
			left = *o
		} else if o, ok := groupMap[csg.LeftChild]; ok {
			left = o
		}
		if o, ok := shapeMap[csg.RightChild]; ok {
			right = *o
		} else if o, ok := groupMap[csg.RightChild]; ok {
			right = o
		}
		newCSG := object.NewCsg(csg.Operation, left, right)
		if csg.Material != nil {
			newCSG.SetMaterial(getMaterial(csg.Material))
		}
		newCSG.SetTransform(getTransforms(csg.Transform)...)
		csgMap[csg.Name] = newCSG
		newCSG.Divide(1)
		csgs = append(csgs, newCSG)
	}

	// Now update groups
	for _, grp := range sceneStruct.Groups {
		toAdd := make([]object.Object, 0, len(grp.Children))
		group := groupMap[grp.Name]
		for _, child := range grp.Children {
			if o, ok := shapeMap[child]; ok {
				toAdd = append(toAdd, *o)
			}
			if o, ok := csgMap[child]; ok {
				toAdd = append(toAdd, o)
			}
		}
		group.Add(toAdd...)
		if grp.Material != nil {
			group.SetMaterial(getMaterial(grp.Material))
		}
		group.SetTransform(getTransforms(grp.Transform)...)
		group.Divide(1)
		groups = append(groups, group)
	}
	return append(groups, csgs...)
}

// ParseOBJ parses the supplied OBJ files and creates groupss
func ParseOBJ(files []*schema.File) ([]object.Object, error) {
	groups := make([]object.Object, 0, len(files))
	for _, group := range files {
		parser, err := parser.Parse(group.File)
		if err != nil {
			return nil, fmt.Errorf("error parsing OBJ file: %v", err)
		}
		groups = append(groups, parser.GetGroup())
	}
	return groups, nil
}

func getMaterial(material *schema.Material) *object.Material {
	objMaterial := object.DefaultMaterial
	if material == nil {
		return &objMaterial
	}

	if material.Color != nil {
		rgb := *material.Color
		objMaterial.Color = image.NewColor(rgb[0], rgb[1], rgb[2])
	}
	if material.Pattern != nil {
		var pattern image.Pattern
		rgb1 := material.Pattern.Color1
		rgb2 := material.Pattern.Color2
		color1 := image.NewColor(rgb1[0], rgb1[1], rgb1[2])
		color2 := image.NewColor(rgb2[0], rgb2[1], rgb2[2])

		switch material.Pattern.Type {
		case "checker":
			pattern = image.NewCheckerPattern(color1, color2)
		case "gradient":
			pattern = image.NewGradientPattern(color1, color2)
		case "ring":
			pattern = image.NewRingPattern(color1, color2)
		case "stripe":
			pattern = image.NewStripePattern(color1, color2)
		}
		pattern.SetTransform(getTransforms(material.Pattern.Transform)...)
		objMaterial.Pattern = pattern
	}
	if material.Ambient != nil {
		objMaterial.Ambient = *material.Ambient
	}
	if material.Diffuse != nil {
		objMaterial.Diffuse = *material.Diffuse
	}
	if material.Specular != nil {
		objMaterial.Specular = *material.Specular
	}
	if material.Shininess != nil {
		objMaterial.Shininess = *material.Shininess
	}
	if material.Reflective != nil {
		objMaterial.Reflective = *material.Reflective
	}
	if material.Transparency != nil {
		objMaterial.Ambient = *material.Transparency
	}
	if material.RefractiveIndex != nil {
		objMaterial.RefractiveIndex = *material.RefractiveIndex
	}
	if material.Shadow != nil {
		objMaterial.Shadow = *material.Shadow
	}
	return &objMaterial
}

func getTransforms(transforms []*schema.Transform) []*base.Matrix {
	t := make([]*base.Matrix, 0, len(transforms))
	for _, transform := range transforms {
		switch transform.Type {
		case "translate":
			t = append(t, base.Translate(transform.Values[0], transform.Values[1], transform.Values[2]))
		case "scale":
			t = append(t, base.Scale(transform.Values[0], transform.Values[1], transform.Values[2]))
		case "rotate":
			x := transform.Values[0] * math.Pi / 180
			y := transform.Values[1] * math.Pi / 180
			z := transform.Values[2] * math.Pi / 180
			if transform.Values[0] != 0 {
				t = append(t, base.RotateX(x))
			}
			if transform.Values[1] != 0 {
				t = append(t, base.RotateY(y))
			}
			if transform.Values[2] != 0 {
				t = append(t, base.RotateZ(z))
			}
			// case "shear": still need to figure out the JSON for this
		}
	}
	return t
}
