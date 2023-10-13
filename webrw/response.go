package webrw

import (
	"io"
	"net/http"
)

type Response struct {
	rw   http.ResponseWriter
	jm   JSONMarshaler
	xm   XMLMarshaler
	sent bool
}

func (r *Response) Status(status int) {
	if r.sent {
		return
	}
	r.sent = true
	r.rw.WriteHeader(status)
}

func (r *Response) SetContentType(typ string) {
	r.rw.Header().Set(HeaderContentType, typ)
}

func (r *Response) Write(b []byte) error {
	_, err := r.rw.Write(b)
	return err
}

func (r *Response) Content(status int, typ string, b []byte) error {
	r.SetContentType(typ)
	r.Status(status)
	return r.Write(b)
}

func (r *Response) NoContent() {
	r.Status(http.StatusNoContent)
}

func (r *Response) PlainText(status int, text string) error {
	return r.Content(status, MIMETextPlainCharsetUTF8, []byte(text))
}

func (r *Response) HTML(status int, html string) error {
	return r.Content(status, MIMETextHTMLCharsetUTF8, []byte(html))
}

func (r *Response) JSON(status int, data any) error {
	r.SetContentType(MIMEApplicationJSONCharsetUTF8)
	r.Status(status)
	return r.jm.Marshal(r.rw, data)
}

func (r *Response) XML(status int, data any) error {
	r.SetContentType(MIMEApplicationXMLCharsetUTF8)
	r.Status(status)
	return r.xm.Marshal(r.rw, data)
}

func (r *Response) Redirect(status int, url string) error {
	// standard defines codes from 300 to 309, but custom HTTP clients might handle their own codes
	// correspondingly for 3xx range
	if status < 300 || status > 399 {
		return ErrInvalidRedirectStatusCode
	}
	r.rw.Header().Set(HeaderLocation, url)
	r.Status(status)
	return nil
}

func (r *Response) Stream(status int, typ string, rdr io.Reader) error {
	r.SetContentType(typ)
	r.Status(status)
	_, err := io.Copy(r.rw, rdr)
	return err
}

func (r *Response) Sent() bool {
	return r.sent
}
