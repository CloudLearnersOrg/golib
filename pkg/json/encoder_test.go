package json_test

import (
	"bytes"
	"sync"
	"testing"

	"github.com/CloudLearnersOrg/golib/pkg/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "string map",
			input:    map[string]string{"key": "value"},
			expected: "{\"key\":\"value\"}\n",
		},
		{
			name:     "int map",
			input:    map[string]int{"count": 42},
			expected: "{\"count\":42}\n",
		},
		{
			name:     "struct",
			input:    struct{ Name string }{"John"},
			expected: "{\"Name\":\"John\"}\n",
		},
		{
			name:     "nil",
			input:    nil,
			expected: "null\n",
		},
	}

	var wg sync.WaitGroup
	wg.Add(len(testCases))

	for _, tc := range testCases {
		tc := tc // Capture range variable
		go func() {
			defer wg.Done()
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				// Given
				buf := &bytes.Buffer{}

				// When
				err := json.Encode(buf, tc.input)

				// Then
				require.NoError(t, err)
				assert.Equal(t, tc.expected, buf.String())
			})
		}()
	}

	wg.Wait()
}

func TestDecode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		target   interface{}
		expected interface{}
	}{
		{
			name:     "string map",
			input:    "{\"key\":\"value\"}",
			target:   &map[string]string{},
			expected: map[string]string{"key": "value"},
		},
		{
			name:     "int map",
			input:    "{\"count\":42}",
			target:   &map[string]int{},
			expected: map[string]int{"count": 42},
		},
		{
			name:     "struct",
			input:    "{\"Name\":\"John\"}",
			target:   &struct{ Name string }{},
			expected: struct{ Name string }{"John"},
		},
	}

	var wg sync.WaitGroup
	wg.Add(len(testCases))

	for _, tc := range testCases {
		tc := tc // Capture range variable
		go func() {
			defer wg.Done()
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				// Given
				buf := bytes.NewBufferString(tc.input)

				// When
				err := json.Decode(buf, tc.target)

				// Then
				require.NoError(t, err)

				// Dereference pointers before comparing
				switch v := tc.target.(type) {
				case *map[string]string:
					assert.Equal(t, tc.expected, *v)
				case *map[string]int:
					assert.Equal(t, tc.expected, *v)
				case *struct{ Name string }:
					assert.Equal(t, tc.expected, *v)
				default:
					t.Fatalf("Unsupported target type: %T", tc.target)
				}
			})
		}()
	}

	wg.Wait()
}

func TestMarshal(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "string map",
			input:    map[string]string{"key": "value"},
			expected: "{\"key\":\"value\"}",
		},
		{
			name:     "int map",
			input:    map[string]int{"count": 42},
			expected: "{\"count\":42}",
		},
		{
			name:     "struct",
			input:    struct{ Name string }{"John"},
			expected: "{\"Name\":\"John\"}",
		},
		{
			name:     "nil",
			input:    nil,
			expected: "null",
		},
	}

	var wg sync.WaitGroup
	wg.Add(len(testCases))

	for _, tc := range testCases {
		tc := tc // Capture range variable
		go func() {
			defer wg.Done()
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				// Given / When
				result, err := json.Marshal(tc.input)

				// Then
				require.NoError(t, err)
				assert.Equal(t, tc.expected, string(result))
			})
		}()
	}

	wg.Wait()
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		target   interface{}
		expected interface{}
	}{
		{
			name:     "string map",
			input:    "{\"key\":\"value\"}",
			target:   &map[string]string{},
			expected: map[string]string{"key": "value"},
		},
		{
			name:     "int map",
			input:    "{\"count\":42}",
			target:   &map[string]int{},
			expected: map[string]int{"count": 42},
		},
		{
			name:     "struct",
			input:    "{\"Name\":\"John\"}",
			target:   &struct{ Name string }{},
			expected: struct{ Name string }{"John"},
		},
	}

	var wg sync.WaitGroup
	wg.Add(len(testCases))

	for _, tc := range testCases {
		tc := tc // Capture range variable
		go func() {
			defer wg.Done()
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				// Given / When
				err := json.Unmarshal([]byte(tc.input), tc.target)

				// Then
				require.NoError(t, err)

				// Dereference pointers before comparing
				switch v := tc.target.(type) {
				case *map[string]string:
					assert.Equal(t, tc.expected, *v)
				case *map[string]int:
					assert.Equal(t, tc.expected, *v)
				case *struct{ Name string }:
					assert.Equal(t, tc.expected, *v)
				default:
					t.Fatalf("Unsupported target type: %T", tc.target)
				}
			})
		}()
	}

	wg.Wait()
}

func TestErrorCases(t *testing.T) {
	t.Parallel()

	t.Run("decode invalid json", func(t *testing.T) {
		t.Parallel()
		// Given
		invalidJSON := "{invalid:json}"
		var result map[string]string

		// When
		err := json.Decode(bytes.NewBufferString(invalidJSON), &result)

		// Then
		assert.Error(t, err)
	})

	t.Run("unmarshal invalid json", func(t *testing.T) {
		t.Parallel()
		// Given
		invalidJSON := "{invalid:json}"
		var result map[string]string

		// When
		err := json.Unmarshal([]byte(invalidJSON), &result)

		// Then
		assert.Error(t, err)
	})

	t.Run("encode invalid value", func(t *testing.T) {
		t.Parallel()
		// Given
		// Create a circular reference which can't be marshaled
		type circular struct {
			Self *circular
		}
		value := &circular{}
		value.Self = value // circular reference
		buf := &bytes.Buffer{}

		// When
		err := json.Encode(buf, value)

		// Then
		assert.Error(t, err)
	})

	t.Run("marshal invalid value", func(t *testing.T) {
		t.Parallel()
		// Given
		// Create a circular reference which can't be marshaled
		type circular struct {
			Self *circular
		}
		value := &circular{}
		value.Self = value // circular reference

		// When
		_, err := json.Marshal(value)

		// Then
		assert.Error(t, err)
	})
}
