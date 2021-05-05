package spec

import (
	"fmt"
	"mime"
	"net/http"
	"path"
	"strings"
)

type HttpProvider struct {
	*baseProvider
	url string
}

func NewHttpProvider(url string) *HttpProvider {
	return &HttpProvider{baseProvider: newBaseProvider(), url: url}
}

func (hp *HttpProvider) Load() (err error) {
	res, err := http.Get(hp.url)
	if err != nil {
		return
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("request failed: %d", res.StatusCode)
	}

	var buf []byte
	if _, err = res.Body.Read(buf); err != nil {
		return
	}

	ext := strings.ToLower(path.Ext(hp.url))
	if ext == "" {
		ext, _, err = mime.ParseMediaType(res.Header.Get("content-type"))
		if err != nil {
			return
		}
	}

	err = hp.parseAndSet(buf, ext)
	return
}
