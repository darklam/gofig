package providers

import (
	"errors"
	json "github.com/titanous/json5"
	"os"
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

func (jp JSONProvider) GetValue(fieldPath []string) (string, error) {
	var currentValue interface{} = jp.parsedFile

	for _, path := range fieldPath {
		m, ok := currentValue.(map[string]interface{})
		if !ok {
			return "", errors.New("received invalid path or JSON file")
		}

		currentValue, ok = m[path]
		if !ok {
			break
		}
	}

	if currentValue == nil {
		return "", errors.New("key not found in JSON")
	}

	strValue, ok := currentValue.(string)
	if !ok {
		return "", errors.New("got invalid value")
	}

	return strValue, nil
}
