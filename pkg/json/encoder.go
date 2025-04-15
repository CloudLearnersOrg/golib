package json

import (
	"encoding/json"
	"io"
)

func Encode(w io.Writer, data any) error {
	return json.NewEncoder(w).Encode(data)
}

func Decode(r io.Reader, data any) error {
	return json.NewDecoder(r).Decode(data)
}

func Marshal(data any) ([]byte, error) {
	return json.Marshal(data)
}

func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
