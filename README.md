# Golang Ray Tracer

A [ray tracer](https://en.wikipedia.org/wiki/Ray_tracing_(graphics)) written in Golang. Implemented with much help from the tremendous book by Jamis Buck, [The Ray Tracer Challenge](https://pragprog.com/titles/jbtracer/the-ray-tracer-challenge/).

### Building/Running

#### Pre-requisites:
- Golang [installed](https://go.dev/doc/install)

`make build` builds the ray tracer binary.

Specify a scene file when running:
`./gtracer --scene my-scene.yaml`

Both YAML and JSON file types are supported. See the `demo/` directory for some example scenes. The schema for the scene file can be viewed [here](schema/README.md).

**Important Notes:**
1. In a scene definition, children listed in either a group or csg need to be defined as a top level object (either as shape, file, group, or csg) in order to be properly referenced.
2. Child objects must be defined before their parents.
3. As of now, materials defined on a child within nested groups may not be honored. For best results, avoid nested groups and csgs.
