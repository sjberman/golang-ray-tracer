package internal

import (
	"fmt"
	"log"
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
func CreateShapes(shapes []*schema.Shape) ([]object.Object, map[string]object.Object) {
	objs := []object.Object{}
	shapeMap := make(map[string]object.Object)
	for _, shape := range shapes {
		var obj, parent object.Object
		var inheritedMaterial *object.Material
		var inheritedTform *base.Matrix

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
			obj = object.NewCube()
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
			obj = object.NewPlane()
		case "sphere":
			obj = object.NewSphere()
		case "glassSphere":
			obj = object.GlassSphere()
		}
		if shape.Inherits != nil {
			parent = shapeMap[*shape.Inherits]
			if parent == nil {
				log.Fatalf("Shape '%s' inherits from shape '%s', which must be defined prior.",
					shape.Name, *shape.Inherits)
			}
			inheritedMaterial = parent.GetMaterial()
			inheritedTform = parent.GetTransform()
		}

		obj.SetMaterial(getMaterial(shape.Material, inheritedMaterial))
		obj.SetTransform(getTransforms(shape.Transform, inheritedTform)...)
		objs = append(objs, obj)
		shapeMap[shape.Name] = obj
	}
	return objs, shapeMap
}

// CreateGroupsAndCSGs builds the group and csg objects using the spec
func CreateGroupsAndCSGs(
	sceneStruct schema.RayTracerScene,
	shapeMap,
	objMap map[string]object.Object,
) ([]object.Object, []string, []string) {
	var usedShapes, usedOBJGroups []string
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
		left := getChild(&csg.LeftChild, shapeMap, objMap, groupMap, nil, &usedShapes, &usedOBJGroups)
		right := getChild(&csg.RightChild, shapeMap, objMap, groupMap, nil, &usedShapes, &usedOBJGroups)

		newCSG := object.NewCsg(csg.Operation, left, right)
		if csg.Material != nil {
			newCSG.SetMaterial(getMaterial(csg.Material, nil))
		}
		newCSG.SetTransform(getTransforms(csg.Transform, nil)...)
		newCSG.Divide(2)
		csgMap[csg.Name] = newCSG
		csgs = append(csgs, newCSG)
	}

	// Now update groups
	for _, grp := range sceneStruct.Groups {
		toAdd := make([]object.Object, 0, len(grp.Children))
		group := groupMap[grp.Name]

		for _, child := range grp.Children {
			childObj := getChild(&child, shapeMap, objMap, groupMap, csgMap, &usedShapes, &usedOBJGroups)
			toAdd = append(toAdd, childObj)
		}
		group.Add(toAdd...)
		if grp.Material != nil {
			group.SetMaterial(getMaterial(grp.Material, nil))
		}
		group.SetTransform(getTransforms(grp.Transform, nil)...)
		group.Divide(2)
		groups = append(groups, group)
	}
	return append(groups, csgs...), usedShapes, usedOBJGroups
}

func getChild(
	child *schema.ObjectShell,
	shapeMap,
	objMap map[string]object.Object,
	groupMap map[string]*object.Group,
	csgMap map[string]*object.Csg,
	usedShapes,
	usedOBJGroups *[]string,
) object.Object {
	var childObject object.Object
	var inheritedMaterial *object.Material
	var inheritedTform *base.Matrix

	if o, ok := shapeMap[child.Name]; ok {
		childObject = o.DeepCopy()
		*usedShapes = append(*usedShapes, child.Name)
		inheritedMaterial = o.GetMaterial()
		inheritedTform = o.GetTransform()
	}
	if csgMap != nil {
		if o, ok := csgMap[child.Name]; ok {
			childObject = o.DeepCopy()
			inheritedMaterial = o.GetMaterial()
			inheritedTform = o.GetTransform()
		}
	}
	if groupMap != nil {
		if o, ok := groupMap[child.Name]; ok {
			childObject = o.DeepCopy()
			inheritedMaterial = o.GetMaterial()
			inheritedTform = o.GetTransform()
		}
	}
	if o, ok := objMap[child.Name]; ok {
		childObject = o.DeepCopy()
		*usedOBJGroups = append(*usedOBJGroups, child.Name)
		inheritedMaterial = o.GetMaterial()
		inheritedTform = o.GetTransform()
	}

	if child.Material != nil {
		childObject.SetMaterial(getMaterial(child.Material, inheritedMaterial))
	}
	if child.Transform != nil {
		childObject.SetTransform(getTransforms(child.Transform, inheritedTform)...)
	}
	return childObject
}

// ParseOBJ parses the supplied OBJ files and creates groups
func ParseOBJ(files []*schema.File) ([]object.Object, map[string]object.Object, error) {
	groups := make([]object.Object, 0, len(files))
	objMap := make(map[string]object.Object)

	for _, grp := range files {
		parser, err := parser.Parse(grp.File)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing OBJ file: %v", err)
		}
		group := parser.GetGroup()
		if grp.Material != nil {
			group.SetMaterial(getMaterial(grp.Material, nil))
		}
		group.SetTransform(getTransforms(grp.Transform, nil)...)
		group.Divide(2)
		groups = append(groups, group)
		objMap[grp.Name] = group
	}
	return groups, objMap, nil
}

func getMaterial(material *schema.Material, inheritedMaterial *object.Material) *object.Material {
	objMaterial := object.DefaultMaterial
	if inheritedMaterial != nil {
		objMaterial = *inheritedMaterial
	}
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
		pattern.SetTransform(getTransforms(material.Pattern.Transform, nil)...)
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
		objMaterial.Transparency = *material.Transparency
	}
	if material.RefractiveIndex != nil {
		objMaterial.RefractiveIndex = *material.RefractiveIndex
	}
	if material.Shadow != nil {
		objMaterial.Shadow = *material.Shadow
	}
	return &objMaterial
}

func getTransforms(transforms []*schema.Transform, inheritedTform *base.Matrix) []*base.Matrix {
	var t []*base.Matrix
	tLength := len(transforms)
	if inheritedTform != nil {
		tLength++
		t = make([]*base.Matrix, 0, tLength)
		t = append(t, inheritedTform)
	} else {
		t = make([]*base.Matrix, 0, tLength)
	}

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
		case "shear":
			xy := transform.Values[0] * math.Pi / 180
			xz := transform.Values[1] * math.Pi / 180
			yx := transform.Values[2] * math.Pi / 180
			yz := transform.Values[3] * math.Pi / 180
			zx := transform.Values[4] * math.Pi / 180
			zy := transform.Values[5] * math.Pi / 180
			t = append(t, base.Shear(xy, xz, yx, yz, zx, zy))
		}
	}
	return t
}
