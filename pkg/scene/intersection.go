package scene

import "sort"

// Intersection keeps track of the value and object of an intersection
type Intersection struct {
	value  float64
	object *Sphere
}

// NewIntersection returns a new Intersection object
func NewIntersection(value float64, object *Sphere) *Intersection {
	return &Intersection{
		value:  value,
		object: object,
	}
}

// GetValue returns the intersection's value
func (i *Intersection) GetValue() float64 {
	return i.value
}

// GetObject returns the intersection's object
func (i *Intersection) GetObject() *Sphere {
	return i.object
}

// Hit returns the closest intersection to the origin
func Hit(intersections []*Intersection) *Intersection {
	if len(intersections) == 0 {
		return nil
	}

	var min *Intersection
	for _, i := range intersections {
		// set initial value of min
		if min == nil && i.value > 0 {
			min = i
		}
		// update min if new value is less
		if i.value > 0 && i.value < min.value {
			min = i
		}
	}
	return min
}

// sorts a list of intersections based on value
func sortIntersections(ints []*Intersection) []*Intersection {
	sort.Slice(ints, func(i, j int) bool {
		return ints[i].value < ints[j].value
	})
	return ints
}

// returns a combined list of the supplied intersections
func intersections(intersections ...*Intersection) []*Intersection {
	return sortIntersections(intersections)
}
