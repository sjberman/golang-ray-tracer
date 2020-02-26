# Golang Ray Tracer

A [ray tracer](https://en.wikipedia.org/wiki/Ray_tracing_(graphics)) written in Golang. Implemented with much help from the tremendous book by Jamis Buck, [The Ray Tracer Challenge](https://pragprog.com/book/jbtracer/the-ray-tracer-challenge).

#### Advantages of using Go
- Easy concurrency. It's almost a necessity to use multithreading in a ray tracer in order to improve performance, and Go's [sync.WaitGroup](https://golang.org/pkg/sync/#WaitGroup) makes this very simple.
- Garbage collection. There are a LOT of objects created while the ray tracer is running, and having those cleaned up automatically is a great convenience.

#### Disadvantages of using Go
- Inheritance/polymorphism is hard. The object classes (shapes) in this ray tracer have a lot of common code, and are designed in such a way that polymorphism is required. I managed to make it work, but it isn't as pretty as it could be using an object-oriented language.

### Building/Running

#### Pre-requisites:
- go 1.13 installed

```make build``` builds the ray tracer binary.

Specify a scene file when running:
```./gtracer --scene my-scene.yaml```

Both YAML and JSON file types are supported. See the `demo/` directory for some example scenes. The schema for the scene file can be viewed [here](schema/README.md).
