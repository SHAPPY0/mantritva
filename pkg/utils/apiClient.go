package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type APIClient struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

func NewAPIClient(baseURL *string) *APIClient {
	ac := &APIClient{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
	if baseURL != nil {
		ac.baseURL = *baseURL
	}
	return ac
}

func (c *APIClient) Call(ctx context.Context, method, path string, body any) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, NewError("marshal: %v", err)
		}
		bodyReader = bytes.NewReader(data)
	}
	url := path
	if c.baseURL != "" {
		url = c.baseURL + path
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, NewError("build request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, NewError("Execute: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		return nil, NewError("API error %d: %s", resp.StatusCode, b)
	}

	return resp, nil
}

// func (c *APIClient) GetPost(ctx context.Context, id int) (*Post, error) {
// 	resp, err := c.do(ctx, http.MethodGet, fmt.Sprintf("/posts/%d", id), nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var post Post
// 	return &post, json.NewDecoder(resp.Body).Decode(&post)
// }
