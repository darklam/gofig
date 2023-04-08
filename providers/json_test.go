package providers

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testJSON = `
{
	"key1": "value1",
	"key2": {
		"key3": "value2",
		"key4": {
			"key5": "value3"
		}
	},
	"key6": 123
}
`

func TestJSONProvider_GetValue(t *testing.T) {
	// Create a temporary file with test JSON content
	tmpFile, err := os.CreateTemp("", "tests")
	assert.Nil(t, err)

	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(testJSON)
	assert.Nil(t, err)

	jp, err := NewJSONProvider(tmpFile.Name())
	assert.Nil(t, err)

	tests := []struct {
		name        string
		fieldPath   []string
		expected    string
		expectedErr error
	}{
		{
			name:        "Valid key",
			fieldPath:   []string{"key1"},
			expected:    "value1",
			expectedErr: nil,
		},
		{
			name:        "Nested key",
			fieldPath:   []string{"key2", "key3"},
			expected:    "value2",
			expectedErr: nil,
		},
		{
			name:        "Deeply nested key",
			fieldPath:   []string{"key2", "key4", "key5"},
			expected:    "value3",
			expectedErr: nil,
		},
		{
			name:        "Non-string value",
			fieldPath:   []string{"key6"},
			expected:    "",
			expectedErr: errors.New("got invalid value"),
		},
		{
			name:        "Invalid key",
			fieldPath:   []string{"nonexistent"},
			expected:    "",
			expectedErr: errors.New("key not found in JSON"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, err := jp.GetValue(test.fieldPath)
			assert.Equal(t, test.expected, value)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
