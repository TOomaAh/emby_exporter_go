package request

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type Client interface {
	GetBaseURL() *url.URL
	ApplyAuthentication(r *http.Request) error
	SetHeaders(headers http.Header)
	GetClient() *http.Client
}

var (
	ErrorCannotParsePath = errors.New("cannot parse path")
	ErrorCannotReadBody  = errors.New("cannot read body")
	ErrorInvalidURL      = errors.New("invalid URL")
	ErrorInvalidRequest  = errors.New("invalid request")
	Error400BadRequest   = errors.New("400 bad request")
	Error404NotFound     = errors.New("404 not found")
	Error500Internal     = errors.New("500 internal server error")
	Error401Unauthorized = errors.New("401 unauthorized")
	Error403Forbidden    = errors.New("403 forbidden")
)

type RequestManager struct {
	client Client
}

func NewRequestManager(c Client) *RequestManager {
	return &RequestManager{
		client: c,
	}
}

func (r *RequestManager) NewRequest(method string, path string, body io.Reader) (*http.Request, error) {
	u, err := r.client.GetBaseURL().Parse(path)
	if err != nil {
		return nil, ErrorCannotParsePath
	}
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, ErrorInvalidRequest
	}
	r.client.SetHeaders(req.Header)
	r.client.ApplyAuthentication(req)
	return req, nil
}

func parseStatusCode(code int) error {
	switch code {
	case http.StatusBadRequest:
		return Error400BadRequest
	case http.StatusNotFound:
		return Error404NotFound
	case http.StatusInternalServerError:
		return Error500Internal
	case http.StatusUnauthorized:
		return Error401Unauthorized
	case http.StatusForbidden:
		return Error403Forbidden
	}
	return nil
}

func (r *RequestManager) Do(req *http.Request, v interface{}) error {
	resp, err := r.client.GetClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := parseStatusCode(resp.StatusCode); err != nil {
		return err
	}

	if v != nil {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return ErrorCannotReadBody
		}
		return json.Unmarshal(body, v)
	}
	return nil
}

func (r *RequestManager) DoFile(req *http.Request, out io.Writer) error {
	resp, err := r.client.GetClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := parseStatusCode(resp.StatusCode); err != nil {
		return err
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}
	return nil
}
