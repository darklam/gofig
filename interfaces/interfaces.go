// Package interfaces includes common interfaces used in the project
package interfaces

// Provider is the interface that every provider must implement to be used with Gofig
//
//counterfeiter:generate . Provider
type Provider interface {
	// GetValue returns the value for a struct field given its path in the struct
	// If a value is not found, it should return an empty string without an error
	GetValue(fieldPath []string) (string, error)
}
