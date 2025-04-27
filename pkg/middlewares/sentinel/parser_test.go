package sentinel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStringToHosts(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expected    []string
		description string
	}{
		{
			name:        "EmptyString",
			input:       "",
			expected:    []string{},
			description: "Empty string should return empty slice",
		},
		{
			name:        "SingleHost",
			input:       "example.com",
			expected:    []string{"example.com"},
			description: "Single host should be parsed correctly",
		},
		{
			name:        "MultipleHosts",
			input:       "example.com,api.example.com,localhost:8080",
			expected:    []string{"example.com", "api.example.com", "localhost:8080"},
			description: "Multiple hosts should be parsed correctly",
		},
		{
			name:        "HostsWithWhitespace",
			input:       "  example.com , api.example.com  ,  localhost:8080 ",
			expected:    []string{"example.com", "api.example.com", "localhost:8080"},
			description: "Whitespace should be trimmed",
		},
		{
			name:        "EmptyEntries",
			input:       "example.com,,localhost:8080,",
			expected:    []string{"example.com", "localhost:8080"},
			description: "Empty entries should be filtered out",
		},
		{
			name:        "OnlyWhitespace",
			input:       "   ,  ,    ",
			expected:    []string{},
			description: "Whitespace-only entries should be filtered out",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// When
			result := ParseStringToHosts(tc.input)

			// Then
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
