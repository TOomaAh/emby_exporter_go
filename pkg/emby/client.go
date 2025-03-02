package emby

import (
	"TOomaAh/emby_exporter_go/pkg/request"
	"net/http"
	"net/url"
	"time"
)

type EmbyClient struct {
	BaseURL    *url.URL
	APIKey     string
	headers    http.Header
	httpClient *http.Client
}

func NewEmbyClient(baseURL string, apiKey string) (request.Client, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, request.ErrorInvalidURL
	}

	headers := http.Header{}
	headers.Set("Accept", "application/json")
	headers.Set("Content-Type", "application/json")
	headers.Set("User-Agent", "emby_exporter_go")
	headers.Set("X-Emby-Token", apiKey)

	return &EmbyClient{
		BaseURL: parsed,
		APIKey:  apiKey,
		headers: headers,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (c *EmbyClient) ApplyAuthentication(r *http.Request) error {
	r.Header.Set("X-Emby-Token", c.APIKey)
	return nil
}

func (c *EmbyClient) SetHeaders(headers http.Header) {
	headers.Set("Accept", "application/json")
	headers.Set("Content-Type", "application/json")
	headers.Set("User-Agent", "emby_exporter_go")
}

func (c *EmbyClient) GetBaseURL() *url.URL {
	return c.BaseURL
}

func (c *EmbyClient) GetClient() *http.Client {
	return c.httpClient
}
