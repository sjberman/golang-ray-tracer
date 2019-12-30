package object

import "sort"

// Intersection keeps track of the value and object of an intersection
type Intersection struct {
	Value  float64
	Object Object
	u, v   float64 // only used for triangles
}

// NewIntersection returns a new Intersection object
func NewIntersection(value float64, object Object) *Intersection {
	return &Intersection{
		Value:  value,
		Object: object,
	}
}

// Hit returns the closest intersection to the origin
func Hit(intersections []*Intersection) *Intersection {
	if len(intersections) == 0 {
		return nil
	}

	var min *Intersection
	for _, i := range intersections {
		// set initial value of min
		if min == nil && i.Value > 0 {
			min = i
		}
		// update min if new value is less
		if i.Value > 0 && i.Value < min.Value {
			min = i
		}
	}
	return min
}

// sorts a list of intersections based on value
func sortIntersections(ints []*Intersection) []*Intersection {
	sort.Slice(ints, func(i, j int) bool {
		return ints[i].Value < ints[j].Value
	})
	return ints
}

// Intersections returns a combined list of the supplied intersections
func Intersections(intersections ...*Intersection) []*Intersection {
	return sortIntersections(intersections)
}
