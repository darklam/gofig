// Package interfaces includes common interfaces used in the project
package interfaces

import "reflect"

// Provider is the interface that every provider must implement to be used with Gofig
type Provider interface {
	// Name returns the name of the provider (so the right provider is used when annotating a struct field with
	// `provider:"<name>" where <name> is the string the provider's Name method returns
	Name() string

	// GetValue returns the value of struct field (optionally given the parent)
	// If a value is not found, it should return an empty string without an error
	GetValue(field reflect.StructField, parentField *reflect.StructField) (string, error)
}

// Valuer is used to wrap the Vault API, so it can be easily unit tested
type Valuer interface {
	// Value returns the map present in the mountPath/secretPath
	Value(mountPath string, secretPath string) (map[string]interface{}, error)
}
