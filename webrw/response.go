package webrw

import "net/http"

type Response struct {
	rw        http.ResponseWriter
	marshaler JSONMarshaler
}

func (r *Response) JSON(status int, data any) error {
	r.rw.WriteHeader(status)
	return r.marshaler.Marshal(r.rw, data)
}
