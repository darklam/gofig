package gofig

import (
	"reflect"
	"strings"
)

func pop[T any](slice *[]T) T {
	length := len(*slice)
	last := (*slice)[length-1]
	*slice = (*slice)[:length-1]
	return last
}

func getFields(t reflect.Type, parent *Field, parentValue reflect.Value) []Field {
	visible := reflect.VisibleFields(t)

	fields := make([]Field, len(visible))

	for i, field := range visible {
		currentPath := make([]string, 0)

		if parent != nil {
			parentPath := parent.fullPath
			for _, f := range parentPath {
				currentPath = append(currentPath, f)
			}
		}

		fieldPath := field.Tag.Get("prop")

		parts := strings.Split(fieldPath, ".")

		for _, part := range parts {
			currentPath = append(currentPath, part)
		}

		fields[i] = Field{field: field, parentValue: parentValue, fullPath: currentPath}
	}

	return fields
}
