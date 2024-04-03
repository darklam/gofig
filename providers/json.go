package providers

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	json "github.com/titanous/json5"
)

type JSONProvider struct {
	parsedFile map[string]interface{}
}

func NewJSONProvider(filePath string) (*JSONProvider, error) {
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	parsed := map[string]interface{}{}

	err = json.Unmarshal(contents, &parsed)
	if err != nil {
		return nil, err
	}

	return &JSONProvider{parsedFile: parsed}, nil
}

func NewJSONProviderFromFs(fs fs.FS, filePath string) (*JSONProvider, error) {
	file, err := fs.Open(filePath)
	if err != nil {
		return nil, err
	}

	contents, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	parsed := map[string]interface{}{}

	err = json.Unmarshal(contents, &parsed)
	if err != nil {
		return nil, err
	}

	return &JSONProvider{parsedFile: parsed}, nil
}

func (jp JSONProvider) GetValue(fieldPath []string) (string, error) {
	var currentValue interface{} = jp.parsedFile

	for _, path := range fieldPath {
		m, ok := currentValue.(map[string]interface{})
		if !ok {
			return "", nil
		}

		currentValue, ok = m[path]
		if !ok {
			break
		}
	}

	if currentValue == nil {
		return "", nil
	}

	strValue, ok := currentValue.(string)
	if !ok {
		return "", errors.New(fmt.Sprintf("got invalid value: %+v", currentValue))
	}

	return strValue, nil
}
