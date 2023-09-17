package webrw

import (
	"net/http"
	"strings"
)

type Request struct {
	req         *http.Request
	unmarshaler JSONUnmarshaler
}

func (r *Request) ContentType() string {
	typ := r.Headers().Get(HeaderContentType)
	parts := strings.SplitN(typ, ";", 2)
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}
	return ""
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

func (r *Request) Request() *http.Request {
	return r.req
}

func (r *Request) JSON(to any) error {
	if r.ContentType() == MIMEApplicationJSON {

	}
	return r.unmarshaler.Unmarshal(r.req.Body, to)
}
