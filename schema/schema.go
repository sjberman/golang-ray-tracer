package schema

import "encoding/json"

// Camera
type Camera struct {
	FieldOfView float64   `json:"field-of-view"`
	From        []float64 `json:"from"`
	Height      int       `json:"height"`
	To          []float64 `json:"to"`
	Up          []float64 `json:"up"`
	Width       int       `json:"width"`
}

// Csg
type Csg struct {
	LeftChild  ObjectShell  `json:"leftChild"`
	Material   *Material    `json:"material,omitempty"`
	Name       string       `json:"name"`
	Operation  string       `json:"operation"`
	RightChild ObjectShell  `json:"rightChild"`
	Transform  []*Transform `json:"transform,omitempty"`
}

// File
type File struct {
	File      string       `json:"file"`
	Material  *Material    `json:"material,omitempty"`
	Name      string       `json:"name"`
	Transform []*Transform `json:"transform,omitempty"`
}

// Group
type Group struct {
	Children  []ObjectShell `json:"children"`
	Material  *Material     `json:"material,omitempty"`
	Name      string        `json:"name"`
	Transform []*Transform  `json:"transform,omitempty"`
}

// Light
type Light struct {
	At        []float64 `json:"at"`
	Intensity []float64 `json:"intensity"`
}

// Material
type Material struct {
	Ambient         *float64   `json:"ambient,omitempty"`
	Color           *[]float64 `json:"color,omitempty"`
	Diffuse         *float64   `json:"diffuse,omitempty"`
	Pattern         *Pattern   `json:"pattern,omitempty"`
	Reflective      *float64   `json:"reflective,omitempty"`
	RefractiveIndex *float64   `json:"refractiveIndex,omitempty"`
	Shadow          *bool      `json:"shadow,omitempty"`
	Shininess       *float64   `json:"shininess,omitempty"`
	Specular        *float64   `json:"specular,omitempty"`
	Transparency    *float64   `json:"transparency,omitempty"`
}

// ObjectShell
type ObjectShell struct {
	Material  *Material    `json:"material,omitempty"`
	Name      string       `json:"name"`
	Transform []*Transform `json:"transform,omitempty"`
}

// Pattern
type Pattern struct {
	Color1    []float64    `json:"color1"`
	Color2    []float64    `json:"color2"`
	Transform []*Transform `json:"transform,omitempty"`
	Type      string       `json:"type"`
}

// RayTracerScene
type RayTracerScene struct {
	Camera *Camera  `json:"camera"`
	Csgs   []*Csg   `json:"csgs,omitempty"`
	Files  []*File  `json:"files,omitempty"`
	Groups []*Group `json:"groups,omitempty"`
	Lights []*Light `json:"lights"`
	Shapes []*Shape `json:"shapes,omitempty"`
}

// Shape
type Shape struct {
	Closed    *bool        `json:"closed,omitempty"`
	Inherits  *string      `json:"inherits,omitempty"`
	Material  *Material    `json:"material,omitempty"`
	Maximum   *float64     `json:"maximum,omitempty"`
	Minimum   *float64     `json:"minimum,omitempty"`
	Name      string       `json:"name"`
	Transform []*Transform `json:"transform,omitempty"`
	Type      string       `json:"type"`
}

// Transform
type Transform struct {
	Type   string    `json:"type"`
	Values []float64 `json:"values"`
}

func (strct *RayTracerScene) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "camera":
			if err := json.Unmarshal([]byte(v), &strct.Camera); err != nil {
				return err
			}
		case "csgs":
			if err := json.Unmarshal([]byte(v), &strct.Csgs); err != nil {
				return err
			}
		case "files":
			if err := json.Unmarshal([]byte(v), &strct.Files); err != nil {
				return err
			}
		case "groups":
			if err := json.Unmarshal([]byte(v), &strct.Groups); err != nil {
				return err
			}
		case "lights":
			if err := json.Unmarshal([]byte(v), &strct.Lights); err != nil {
				return err
			}
		case "shapes":
			if err := json.Unmarshal([]byte(v), &strct.Shapes); err != nil {
				return err
			}
		}
	}
	return nil
}
