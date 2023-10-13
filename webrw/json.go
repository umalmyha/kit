package webrw

import (
	"encoding/json"
	"io"
)

type JSONMarshaler interface {
	Marshal(to io.Writer, data any) error
}

type DefaultJSONMarshaler struct{}

func (m *DefaultJSONMarshaler) Marshal(to io.Writer, data any) error {
	// TODO: think of JSON errors
	return json.NewEncoder(to).Encode(data)
}

type JSONUnmarshaler interface {
	Unmarshal(data io.Reader, to any) error
}

type DefaultJSONUnmarshaler struct{}

func (u *DefaultJSONUnmarshaler) Unmarshal(body io.Reader, to any) error {
	// TODO: think of JSON errors
	return json.NewDecoder(body).Decode(to)
}
