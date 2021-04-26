package binding

import "net/http"

const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
	MIMEYAML              = "application/x-yaml"
)

var (
	JSON          = jsonBinding{}
)

type Binding interface {
	Name() string
	Bind(*http.Request, interface{}) error
}

func Default(contentType string) Binding {
	switch contentType {
	case MIMEJSON:
		return JSON
	default:
		return nil
	}
}
