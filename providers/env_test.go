package providers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvProvider_GetValue(t *testing.T) {
	// Set up some test environment variables
	err := os.Setenv("TEST_VAR1", "value1")
	assert.Nil(t, err)

	err = os.Setenv("TEST_VAR2_VAR3", "value2")
	assert.Nil(t, err)

	defer os.Clearenv()

	ep := NewEnvProvider()

	tests := []struct {
		name      string
		fieldPath []string
		expected  string
	}{
		{
			name:      "Valid environment variable",
			fieldPath: []string{"test", "var1"},
			expected:  "value1",
		},
		{
			name:      "Nested environment variable",
			fieldPath: []string{"test", "var2", "var3"},
			expected:  "value2",
		},
		{
			name:      "Non-existent environment variable",
			fieldPath: []string{"nonexistent"},
			expected:  "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, err := ep.GetValue(test.fieldPath)
			assert.Equal(t, test.expected, value)
			assert.NoError(t, err)
		})
	}
}
