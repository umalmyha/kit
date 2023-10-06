package webrw

const (
	charsetUTF8 = "charset=UTF-8"

	HeaderContentType = "Content-Type"

	MIMEApplicationJSON            = "application/json"
	MIMEApplicationJSONCharsetUTF8 = MIMEApplicationJSON + "; " + charsetUTF8

	MIMETextPlain            = "text/plain"
	MIMETextPlainCharsetUTF8 = MIMETextPlain + "; " + charsetUTF8

	MIMETextHTML            = "text/html"
	MIMETextHTMLCharsetUTF8 = MIMETextHTML + "; " + charsetUTF8
)
