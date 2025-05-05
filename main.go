package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/sjberman/golang-ray-tracer/internal"
	"github.com/sjberman/golang-ray-tracer/pkg/scene"
	"github.com/sjberman/golang-ray-tracer/pkg/scene/object"
	"github.com/sjberman/golang-ray-tracer/schema"
	"github.com/xeipuuv/gojsonschema"
)

var (
	schemaFile = flag.String("schema", "schema/schema.json", "Relative path to the schema.json file")
	sceneFile  = flag.String("scene", "", "JSON or YAML file containing scene info")
	outputFile = flag.String("output", "image.ppm", "Image output file (.ppm)")
)

func parseArgs() {
	flag.Parse()
	if *sceneFile == "" {
		log.Fatalf("scene file is a required argument")
	}

	if *schemaFile == "schema/schema.json" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("error getting current working directory: %v\n", err)
		}
		*schemaFile = path.Join(cwd, *schemaFile)
	}

	if !strings.HasSuffix(*sceneFile, ".json") && !strings.HasSuffix(*sceneFile, ".yaml") {
		log.Fatal("scene file must be of type .json or .yaml")
	}
}

// Builds all of the objects defined in the scene.
func getSceneObjects(sceneStruct schema.RayTracerScene) (*scene.Camera, []*scene.PointLight, []object.Object) {
	camera := internal.CreateCamera(sceneStruct.Camera)
	lights := internal.CreateLights(sceneStruct.Lights)
	shapes, shapeMap := internal.CreateShapes(sceneStruct.Shapes)
	objGroups, objMap, err := internal.ParseOBJ(sceneStruct.Files)
	if err != nil {
		log.Fatal(err.Error())
	}

	groups, usedShapes, usedOBJGroups := internal.CreateGroupsAndCSGs(sceneStruct, shapeMap, objMap)

	// De-dupe any objects that are included in a group definition
	shapes = internal.DeDupe(shapes, shapeMap, usedShapes)
	objGroups = internal.DeDupe(objGroups, objMap, usedOBJGroups)

	objects := append(shapes, objGroups...)
	objects = append(objects, groups...)

	return camera, lights, objects
}

func main() {
	startTime := time.Now()
	// f1, _ := os.Create("perfFile")
	// pprof.StartCPUProfile(f1)
	// defer pprof.StopCPUProfile()
	parseArgs()

	sceneBytes, err := os.ReadFile(*sceneFile)
	if err != nil {
		log.Fatalf("error reading scene file: %v\n", err)
	}
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + *schemaFile)
	if strings.HasSuffix(*sceneFile, ".yaml") {
		sceneBytes, err = yaml.YAMLToJSON(sceneBytes)
		if err != nil {
			log.Fatalf("error converting YAML to JSON: %v", err)
		}
	}
	docLoader := gojsonschema.NewStringLoader(string(sceneBytes))

	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		log.Fatalf("could not validate scene with schema: %v\n", err)
	}

	if !result.Valid() {
		fmt.Println("Scene is invalid:")
		for _, desc := range result.Errors() {
			fmt.Println("- ", desc)
		}
		os.Exit(1)
	}

	var sceneStruct schema.RayTracerScene
	err = json.Unmarshal(sceneBytes, &sceneStruct)
	if err != nil {
		log.Fatalf("error unmarshaling scene JSON: %v\n", err)
	}
	camera, lights, objects := getSceneObjects(sceneStruct)

	world := scene.NewWorld(lights, objects)
	canvas := scene.Render(camera, world)
	err = canvas.WriteToFile(*outputFile)
	if err != nil {
		fmt.Println("error writing file: ", err.Error())
	}

	fmt.Println("Total runtime: ", time.Since(startTime).Round(time.Second))
}
