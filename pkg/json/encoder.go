// Package json provides a set of methods for encoding and marshalling of JSON data more efficiently.
// It wraps the standard encoding/json package to provide a simpler interface for common JSON operations.
package json

import (
	"encoding/json"
	"io"
)

// Encode writes the JSON encoding of data to the given io.Writer.
// It returns an error if the encoding fails.
//
// Example:
//
//	var data struct { Name string `json:"name"` }
//	data.Name = "example"
//	err := json.Encode(writer, data)
func Encode(w io.Writer, data any) error {
	return json.NewEncoder(w).Encode(data)
}

// Decode reads JSON-encoded data from the given io.Reader and stores it in the value pointed to by data.
// It returns an error if the decoding fails.
//
// Example:
//
//	var data struct { Name string `json:"name"` }
//	err := json.Decode(reader, &data)
func Decode(r io.Reader, data any) error {
	return json.NewDecoder(r).Decode(data)
}

// Marshal returns the JSON encoding of data as a byte slice.
// It returns an error if the marshaling fails.
//
// Example:
//
//	var data struct { Name string `json:"name"` }
//	data.Name = "example"
//	bytes, err := json.Marshal(data)
func Marshal(data any) ([]byte, error) {
	return json.Marshal(data)
}

// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
// It returns an error if the unmarshaling fails.
//
// Example:
//
//	var data struct { Name string `json:"name"` }
//	err := json.Unmarshal(bytes, &data)
func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
