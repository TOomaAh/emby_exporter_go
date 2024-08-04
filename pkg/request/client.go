package request

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseURL    *url.URL
	APIKey     string
	headers    http.Header
	httpClient *http.Client
}

var (
	ErrorCannotParsePath = errors.New("cannot parse path")
	ErrorCannotReadBody  = errors.New("cannot read body")
	ErrorInvalidURL      = errors.New("invalid URL")
	ErrorInvalidRequest  = errors.New("invalid request")
	Error404NotFound     = errors.New("404 not found")
	Error500Internal     = errors.New("500 internal server error")
	Error401Unauthorized = errors.New("401 unauthorized")
	Error403Forbidden    = errors.New("403 forbidden")
)

func NewClient(baseURL string, apiKey string) (*Client, error) {

	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, ErrorInvalidURL
	}

	headers := http.Header{}
	headers.Set("Accept", "application/json")
	headers.Set("Content-Type", "application/json")
	headers.Set("User-Agent", "emby_exporter_go")
	headers.Set("X-Emby-Token", apiKey)

	return &Client{
		BaseURL: parsed,
		APIKey:  apiKey,
		headers: headers,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil

}

func (c *Client) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, ErrorCannotParsePath
	}
	request := c.BaseURL.ResolveReference(u)
	req, err := http.NewRequest(method, request.String(), body)
	if err != nil {
		return nil, ErrorInvalidRequest
	}
	req.Header = c.headers
	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		return Error404NotFound
	case http.StatusInternalServerError:
		return Error500Internal
	case http.StatusUnauthorized:
		return Error401Unauthorized
	case http.StatusForbidden:
		return Error403Forbidden
	}

	defer resp.Body.Close()

	if v != nil {
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			return ErrorCannotReadBody
		}

		return json.Unmarshal(body, v)
	}
	return nil
}
