package webrw

const (
	charsetUTF8 = "charset=UTF-8"

	MIMETextPlain            = "text/plain"
	MIMETextPlainCharsetUTF8 = MIMETextPlain + "; " + charsetUTF8

	MIMETextHTML            = "text/html"
	MIMETextHTMLCharsetUTF8 = MIMETextHTML + "; " + charsetUTF8

	MIMEApplicationJSON            = "application/json"
	MIMEApplicationJSONCharsetUTF8 = MIMEApplicationJSON + "; " + charsetUTF8

	MIMEApplicationXML            = "application/xml"
	MIMEApplicationXMLCharsetUTF8 = MIMEApplicationXML + "; " + charsetUTF8

	MIMEOctetStream = "application/octet-stream"

	MIMEMultipartFormData         = "multipart/form-data"
	MIMEApplicationFormURLEncoded = "application/x-www-form-urlencoded"
)
