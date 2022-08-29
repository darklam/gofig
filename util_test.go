package gofig

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestUtil_Pop(t *testing.T) {
	tests := []struct {
		a    []int
		want int
	}{
		{[]int{1, 2, 3, 4, 5, 6}, 6},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%+v", tt.a)
		t.Run(
			testname, func(t *testing.T) {
				ans := pop(&tt.a)
				if ans != tt.want {
					t.Errorf("got %d, want %d", tt.a, tt.want)
				}

				if tt.a[len(tt.a)-1] != 5 {
					t.Error("last element not removed from list")
				}
			},
		)
	}
}

func TestUtil_GetVisibleFieldPairs(t *testing.T) {
	type s struct {
		Field1 string
		Field2 string
	}

	st := new(s)

	pairs := getVisibleFieldPairs(reflect.TypeOf(st).Elem(), nil, reflect.ValueOf(st).Elem())

	assert.Equal(t, len(pairs), 2)
	fields := map[string]*fieldPair{}
	for _, pair := range pairs {
		fields[pair.field.Name] = pair
	}

	assert.NotNil(t, fields["Field1"])
	assert.Nil(t, fields["Field1"].parent)

	assert.NotNil(t, fields["Field2"])
	assert.Nil(t, fields["Field2"].parent)
}
