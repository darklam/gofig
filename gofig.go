package gofig

import (
	"errors"
	"fmt"
	"github.com/darklam/gofig/interfaces"
	"os"
	"reflect"
)

// Gofig is the struct containing the registered providers and responsible for populating the provided configuration
type Gofig struct {
	providers map[string]interfaces.Provider
}

// NewGofig returns a new Gofig instance without any provider
func NewGofig() *Gofig {
	return &Gofig{providers: map[string]interfaces.Provider{}}
}

// RegisterProvider registers a new provider for Gofig to use. If the provider's name collides with another
// provider, then only the one registered last will be used
func (gofig *Gofig) RegisterProvider(provider interfaces.Provider) {
	gofig.providers[provider.Name()] = provider
}

// PopulateConfig populates the values of the given config
// The cfg parameter must be a pointer to a struct (not nil)
// and the struct can only contain string values or pointers to other structs
// (these can and should not be initialized)
func (gofig Gofig) PopulateConfig(cfg interface{}) error {
	t := reflect.TypeOf(cfg)
	v := reflect.ValueOf(cfg)

	if v.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	fields := getVisibleFieldPairs(t, nil, v)

	for len(fields) != 0 {
		current := pop(&fields)
		field := current.field
		parent := current.parent
		parentValue := current.parentValue

		fieldValue := parentValue.FieldByName(field.Name)

		if fieldValue.Kind() == reflect.Ptr {
			// We first create an instance of the pointer
			pointerInstance := reflect.New(field.Type)
			fieldPointerInterface := pointerInstance.Elem().Interface()
			fieldPointerInterfaceType := reflect.TypeOf(fieldPointerInterface).Elem()

			if fieldPointerInterfaceType.Kind() != reflect.Struct {
				return errors.New("only struct pointers and strings are allowed")
			}

			// After we have the pointer, we can instantiate the struct
			structInstance := reflect.New(fieldPointerInterfaceType).Interface()

			fieldValue.Set(reflect.ValueOf(structInstance))

			currentFields := getVisibleFieldPairs(reflect.TypeOf(structInstance).Elem(), &field, fieldValue.Elem())
			fields = append(fields, currentFields...)

			continue
		} else if fieldValue.Kind() != reflect.String {
			return errors.New("only struct pointers and strings are allowed")
		}

		tags := field.Tag

		value := ""

		providerName := tags.Get("provider")
		parentProviderName := ""
		if parent != nil {
			parentProviderName = parent.Tag.Get("provider")
		}

		var provider interfaces.Provider = nil

		if parentProviderName != "" {
			provider = gofig.providers[parentProviderName]
		}

		if providerName != "" {
			provider = gofig.providers[providerName]
		}

		env := tags.Get("env")
		defaultValue := tags.Get("default")

		if env == "" && defaultValue == "" && provider == nil {
			errorText := fmt.Sprintf("no provider, env or default for field %s", field.Name)
			return errors.New(errorText)
		}

		if env != "" {
			value = os.Getenv(env)
		}

		if defaultValue != "" && value == "" {
			value = defaultValue
		}

		if provider != nil {
			providerValue, err := provider.GetValue(field, parent)
			if err != nil {
				return err
			}

			if providerValue != "" {
				value = providerValue
			}
		}

		fieldValue.SetString(value)
	}

	return nil
}
