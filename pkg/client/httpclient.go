package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ranna-go/ranna/pkg/models"
)

const (
	defaultVersion   = "v1"
	defaultUserAgent = "ranna/pkg/client"
)

// Options for the HTTP client.
type Options struct {
	Endpoint      string `json:"endpoint"`
	Version       string `json:"version"`
	Authorization string `json:"authorization"`
	UserAgent     string `json:"useragent"`
}

type httpClient struct {
	options *Options
	client  *http.Client
}

// New returns a new HTTP API Client.
//
// An error is returned when the passed options
// are invalid.
func New(options Options) (c Client, err error) {
	if err = checkAndDefaultOptions(&options); err != nil {
		return
	}
	c = &httpClient{
		options: &options,
		client: &http.Client{
			Timeout: 120 * time.Second, // 2 Minute timeout
		},
	}
	return
}

func checkAndDefaultOptions(options *Options) error {
	if options.Endpoint == "" {
		return errors.New("option endpoint must be provided")
	}
	if options.Version == "" {
		options.Version = defaultVersion
	}
	if options.UserAgent == "" {
		options.UserAgent = defaultUserAgent
	}
	return nil
}

func (t *httpClient) request(method, path string, body any, resData any) (err error) {
	url := fmt.Sprintf("%s/%s/%s", t.options.Endpoint, t.options.Version, path)

	var bodyReader io.Reader
	if body != nil {
		buff := bytes.NewBuffer([]byte{})
		if err = json.NewEncoder(buff).Encode(body); err != nil {
			return err
		}
		bodyReader = buff
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", t.options.UserAgent)
	if t.options.Authorization != "" {
		req.Header.Add("Authorization", t.options.Authorization)
	}

	res, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		resErr := &ResponseError{
			ErrorModel: &models.ErrorModel{
				Code:  res.StatusCode,
				Error: "unknown",
			},
		}
		if res.ContentLength > 0 && strings.HasPrefix(res.Header.Get("Content-Type"), "application/json") {
			if err = json.NewDecoder(res.Body).Decode(resErr.ErrorModel); err != nil {
				return err
			}
		}
		return resErr
	}

	return json.NewDecoder(res.Body).Decode(resData)
}

func (t *httpClient) Spec() (spec models.SpecMap, err error) {
	err = t.request("GET", "spec", nil, &spec)
	return spec, err
}

func (t *httpClient) Exec(req models.ExecutionRequest) (res models.ExecutionResponse, err error) {
	err = t.request("POST", "exec", req, &res)
	return res, err
}
