package gofig

import "reflect"

type Field struct {
	field       reflect.StructField
	parentValue reflect.Value
	fullPath    []string
}
