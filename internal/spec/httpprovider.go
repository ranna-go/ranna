package spec

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"
	"strings"
)

// HttpProvider implements Provider for fetching
// spec definitions over an HTTP endpoint.
type HttpProvider struct {
	*baseProvider
	url string
}

// NewHttpProvider initializes a new HttpProvider
// fetching from the given resource url.
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

	buf, err := io.ReadAll(res.Body)
	if err != nil {
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
