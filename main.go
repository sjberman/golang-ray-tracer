package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"

	"github.com/sjberman/golang-ray-tracer/internal"
	"github.com/sjberman/golang-ray-tracer/pkg/scene"
	"github.com/sjberman/golang-ray-tracer/schema"
)

var schemaFile = flag.String("schema", "schema/schema.json", "Relative path to the schema.json file")
var sceneFile = flag.String("scene", "", "JSON or YAML file containing scene info")
var outputFile = flag.String("output", "image.ppm", "Image output file (.ppm)")

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

func main() {
	startTime := time.Now()
	f1, _ := os.Create("perfFile")
	pprof.StartCPUProfile(f1)
	defer pprof.StopCPUProfile()
	parseArgs()

	sceneBytes, err := ioutil.ReadFile(*sceneFile)
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

	camera := internal.CreateCamera(sceneStruct.Camera)
	lights := internal.CreateLights(sceneStruct.Lights)
	shapes, shapeMap := internal.CreateShapes(sceneStruct.Shapes)
	// Note: this could be very buggy; needs some testing
	groups := internal.CreateGroupsAndCSGs(sceneStruct, shapes, shapeMap)
	objGroups, err := internal.ParseOBJ(sceneStruct.Files)
	if err != nil {
		log.Fatal(err.Error())
	}
	objects := append(shapes, objGroups...)
	objects = append(objects, groups...)

	world := scene.NewWorld(lights, objects)
	canvas := scene.Render(camera, world)
	err = canvas.WriteToFile(*outputFile)
	if err != nil {
		fmt.Println("error writing file: ", err.Error())
	}

	fmt.Println("Total runtime: ", time.Since(startTime).Round(time.Second))
}
