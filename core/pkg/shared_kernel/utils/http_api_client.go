package utils

import (
	"net/http"
	"net/url"
	"time"
)

// APIClient wraps an http.Client with a Base URL
type HttpApiClient struct {
	BaseURL string
	client  *http.Client
}

// NewAPIClient creates a new API client
func NewAPIClient(baseURL string, timeOut time.Duration) *HttpApiClient {
	return &HttpApiClient{
		BaseURL: baseURL,
		client: &http.Client{
			Timeout: timeOut,
		},
	}
}

func (c *HttpApiClient) Do(method, endpoint string, r *http.Request) (*http.Response, error) {
	// 1. Create url (BaseUrl + EndPoint)
	fullURL, err := url.JoinPath(c.BaseURL, endpoint)

	// 2. Create new request to destination
	req, err := http.NewRequest(method, fullURL, r.Body)
	if err != nil {
		return nil, err
	}

	// 3. Copy headers
	for k, v := range r.Header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}

	return c.client.Do(req)
}
