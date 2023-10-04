package webrw

import "net/http"

type Response struct {
	rw        http.ResponseWriter
	marshaler JSONMarshaler
}

func (r *Response) JSON(status int, data any) error {
	r.rw.WriteHeader(status)
	r.rw.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	return r.marshaler.Marshal(r.rw, data)
}

func (r *Response) Status(status int) {
	r.rw.WriteHeader(status)
}

func (r *Response) NoContent() {
	r.rw.WriteHeader(http.StatusNoContent)
}
