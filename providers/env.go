package providers

import (
	"os"
	"strings"
)

type EnvProvider struct{}

func (ep EnvProvider) GetValue(fieldPath []string) (string, error) {
	transformedPath := make([]string, len(fieldPath))
	for i, path := range fieldPath {
		transformedPath[i] = strings.ReplaceAll(path, ".", "_")
	}

	envVar := strings.Join(transformedPath, "_")
	envVar = strings.ToUpper(envVar)

	return os.Getenv(envVar), nil
}

func NewEnvProvider() EnvProvider {
	return EnvProvider{}
}
