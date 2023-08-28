package gofig

import (
	"errors"
	"reflect"

	"github.com/darklam/gofig/interfaces"
)

// Gofig is the struct containing the registered providers and responsible for populating the provided configuration
type Gofig struct {
	providers []interfaces.Provider
}

// NewGofig returns a new Gofig instance without any provider
func NewGofig() *Gofig {
	return &Gofig{providers: make([]interfaces.Provider, 0)}
}

// RegisterProvider registers a new provider for Gofig to use. If the provider's name collides with another
// provider, then only the one registered last will be used
func (gofig *Gofig) RegisterProvider(provider interfaces.Provider) {
	gofig.providers = append(gofig.providers, provider)
}

// PopulateConfig populates the values of the given config
// The cfg parameter must be a pointer to a struct (not nil)
// and the struct can only contain string values or pointers to other structs
// (these can and should not be initialized)
func (gofig *Gofig) PopulateConfig(cfg interface{}) error {
	// Get the reflect.Type and reflect.Value of the provided configuration struct
	t := reflect.TypeOf(cfg)
	v := reflect.ValueOf(cfg)

	// If the provided configuration is a pointer, dereference it to get the underlying type and value
	if v.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	// Get all the top-level fields of the provided configuration struct
	fields := getFields(t, nil, v)

	// Iterate over the fields to populate their values
	for len(fields) != 0 {
		// Get the current field and remove it from the fields list
		current := pop(&fields)
		field := current.field

		// Get the reflect.Value of the current field
		fieldValue := current.parentValue.FieldByName(field.Name)

		// Check if the current field is a pointer to another struct
		if fieldValue.Kind() == reflect.Ptr {
			// We first create an instance of the pointer
			pointerInstance := reflect.New(field.Type)
			fieldPointerInterface := pointerInstance.Elem().Interface()
			fieldPointerInterfaceType := reflect.TypeOf(fieldPointerInterface).Elem()

			// Ensure the pointed type is a struct, otherwise return an error
			if fieldPointerInterfaceType.Kind() != reflect.Struct {
				return errors.New("only struct pointers and strings are allowed")
			}

			// Instantiate the struct and assign it to the pointer field
			structInstance := reflect.New(fieldPointerInterfaceType).Interface()
			fieldValue.Set(reflect.ValueOf(structInstance))

			// Get the fields of the newly instantiated struct and add them to the fields list
			currentFields := getFields(reflect.TypeOf(structInstance).Elem(), &current, fieldValue.Elem())
			fields = append(fields, currentFields...)

			continue
		} else if fieldValue.Kind() != reflect.String {
			// Ensure the field is a string, otherwise return an error
			return errors.New("only struct pointers and strings are allowed")
		}

		// Get the default value for the field from its tag
		value := field.Tag.Get("default")

		// Iterate over the registered providers to resolve the value for the current field
		for _, provider := range gofig.providers {
			resolved, err := provider.GetValue(current.fullPath)
			if err != nil {
				return err
			}
			if resolved == "" {
				continue
			}

			// Use the resolved value if it's not empty
			value = resolved
		}

		fieldValue.SetString(value)
	}

	return nil
}
