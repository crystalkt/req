package req

import "io"

const (
	hdrUserAgentKey   = "User-Agent"
	hdrUserAgentValue = "req/v3 (https://github.com/imroc/req)"
	hdrContentTypeKey = "Content-Type"
	plainTextType     = "text/plain; charset=utf-8"
	jsonContentType   = "application/json; charset=utf-8"
	xmlContentType    = "text/xml; charset=utf-8"
	formContentType   = "application/x-www-form-urlencoded"
)

type uploadFile struct {
	ParamName string
	FilePath  string
	io.Reader
}
