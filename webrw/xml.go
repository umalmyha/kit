package webrw

import (
	"encoding/xml"
	"io"
)

type XMLMarshaler interface {
	Marshal(to io.Writer, data any) error
}

type DefaultXMLMarshaler struct{}

func (m *DefaultXMLMarshaler) Marshal(to io.Writer, data any) error {
	// TODO: think of XML errors
	return xml.NewEncoder(to).Encode(data)
}

type XMLUnmarshaler interface {
	Unmarshal(data io.Reader, to any) error
}

type DefaultXMLUnmarshaler struct{}

func (m *DefaultXMLUnmarshaler) Unmarshal(data io.Reader, to any) error {
	// TODO: think of XML errors
	return xml.NewDecoder(data).Decode(to)
}
