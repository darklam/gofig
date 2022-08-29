package gofig

import "reflect"

type fieldPair struct {
	field       reflect.StructField
	parent      *reflect.StructField
	parentValue reflect.Value
}
