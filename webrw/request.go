package webrw

import (
	"net/http"
	"strings"
)

type Request struct {
	req *http.Request
	ju  JSONUnmarshaler
	xu  XMLUnmarshaler
}

func (r *Request) ContentType() string {
	typ := r.req.Header.Get(HeaderContentType)
	idx := strings.IndexByte(typ, ';')
	if idx == -1 {
		return typ
	}
	return typ[:idx]
}

func (r *Request) Headers() http.Header {
	// TODO: maybe own type
	return r.req.Header
}

func (r *Request) Cookies() []*http.Cookie {
	// TODO: maybe own type
	return r.req.Cookies()
}

func (r *Request) Cookie(name string) (*http.Cookie, error) {
	// TODO: maybe own type
	return r.req.Cookie(name)
}

func (r *Request) Body(to any) error {
	switch r.ContentType() {
	case MIMEApplicationJSON:
		return r.JSON(to)
	case MIMEApplicationXML:
		return r.XML(to)
	case MIMEApplicationFormURLEncoded, MIMEMultipartFormData:
	default:
		return ErrUnsupportedMediaType
	}
}

func (r *Request) JSON(to any) error {
	return r.ju.Unmarshal(r.req.Body, to)
}

func (r *Request) XML(to any) error {
	return r.xu.Unmarshal(r.req.Body, to)
}
