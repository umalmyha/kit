package webrw

import "net/http"

type Response struct {
	rw        http.ResponseWriter
	marshaler JSONMarshaler
}

func (r *Response) Status(status int) {
	r.rw.WriteHeader(status)
}

func (r *Response) SetContentType(typ string) {
	r.rw.Header().Set(HeaderContentType, typ)
}

func (r *Response) Write(b []byte) error {
	_, err := r.rw.Write(b)
	return err
}

func (r *Response) NoContent() {
	r.Status(http.StatusNoContent)
}

func (r *Response) PlainText(status int, text string) error {
	r.Status(status)
	r.SetContentType(MIMETextPlainCharsetUTF8)
	return r.Write([]byte(text))
}

func (r *Response) HTML(status int, html string) error {
	r.Status(status)
	r.SetContentType(MIMETextHTMLCharsetUTF8)
	return r.Write([]byte(html))
}

func (r *Response) JSON(status int, data any) error {
	r.Status(status)
	r.SetContentType(MIMEApplicationJSONCharsetUTF8)
	return r.marshaler.Marshal(r.rw, data)
}
