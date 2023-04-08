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
				assert.Equal(t, tt.want, ans)
				assert.Equal(t, tt.a[len(tt.a)-1], 5, "last element not removed from list")
			},
		)
	}
}

func TestUtil_GetVisibleFieldPairs(t *testing.T) {
	// GIVEN
	type s struct {
		Field1              string `prop:"field1"`
		SuperSecretPassword string `prop:"super.secret.password"`
	}

	st := new(s)

	parentValue := reflect.ValueOf(st).Elem()

	// WHEN
	pairs := getFields(reflect.TypeOf(st).Elem(), nil, parentValue)

	// THEN
	assert.Equal(t, len(pairs), 2)
	fields := map[string]Field{}
	for _, pair := range pairs {
		fields[pair.field.Name] = pair
	}

	field1 := fields["Field1"]
	assert.NotNil(t, field1)
	assert.Equal(t, field1.parentValue, parentValue)
	assert.Equal(t, len(field1.fullPath), 1)
	assert.Equal(t, field1.fullPath[0], "field1")

	field2 := fields["SuperSecretPassword"]
	assert.NotNil(t, field2)
	assert.Equal(t, field2.parentValue, parentValue)
	assert.Equal(t, len(field2.fullPath), 3)
	assert.Equal(t, field2.fullPath[0], "super")
	assert.Equal(t, field2.fullPath[1], "secret")
	assert.Equal(t, field2.fullPath[2], "password")
}
