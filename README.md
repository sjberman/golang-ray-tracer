# Golang Ray Tracer

A [ray tracer](https://en.wikipedia.org/wiki/Ray_tracing_(graphics)) written in Golang. Implemented with much help from the tremendous book by Jamis Buck, [The Ray Tracer Challenge](https://pragprog.com/book/jbtracer/the-ray-tracer-challenge).

### Building/Testing

#### Pre-requisites:
- go 1.13 installed
- golangci-lint installed (for linting purposes only)

```make build``` builds the ray tracer binary.

```make test``` runs the unit tests and prints code coverage.

```make lint``` runs the linting tool on the code.

```make deps``` updates the third-party dependencies used by the code.

**Note:** main.go is currently a sandbox for testing out various scenes, and requires manually updating the code to render a desired scene. It can accept an OBJ file via the `--file` flag.
