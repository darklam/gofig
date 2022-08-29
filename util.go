package gofig

import "reflect"

func pop[T any](slice *[]T) T {
	length := len(*slice)
	last := (*slice)[length-1]
	*slice = (*slice)[:length-1]
	return last
}

func getVisibleFieldPairs(t reflect.Type, parent *reflect.StructField, parentValue reflect.Value) []*fieldPair {
	visible := reflect.VisibleFields(t)

	fields := make([]*fieldPair, len(visible))

	for i, field := range visible {
		fields[i] = &fieldPair{field: field, parent: parent, parentValue: parentValue}
	}

	return fields
}
