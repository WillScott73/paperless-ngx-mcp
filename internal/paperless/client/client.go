package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	BaseURL    string
	APIToken   string
	httpClient *http.Client
}

func New(baseURL, apiToken string) *Client {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return &Client{
		BaseURL:  baseURL,
		APIToken: apiToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) do(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	url := c.BaseURL + strings.TrimPrefix(path, "/")
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Token "+c.APIToken)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("api error (status %d): %s", resp.StatusCode, string(data))
	}

	return resp, nil
}

func (c *Client) Get(ctx context.Context, path string, target interface{}) error {
	resp, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func (c *Client) SearchDocuments(ctx context.Context, query string, page, pageSize int) (interface{}, error) {
	path := fmt.Sprintf("documents/?query=%s&page=%d&page_size=%d", query, page, pageSize)
	var result interface{}
	err := c.Get(ctx, path, &result)
	return result, err
}

func (c *Client) GetDocument(ctx context.Context, id int) (interface{}, error) {
	path := fmt.Sprintf("documents/%d/", id)
	var result interface{}
	err := c.Get(ctx, path, &result)
	return result, err
}

func (c *Client) ListTags(ctx context.Context, page, pageSize int) (interface{}, error) {
	path := fmt.Sprintf("tags/?page=%d&page_size=%d", page, pageSize)
	var result interface{}
	err := c.Get(ctx, path, &result)
	return result, err
}

func (c *Client) ListCorrespondents(ctx context.Context, page, pageSize int) (interface{}, error) {
	path := fmt.Sprintf("correspondents/?page=%d&page_size=%d", page, pageSize)
	var result interface{}
	err := c.Get(ctx, path, &result)
	return result, err
}

func (c *Client) ListDocumentTypes(ctx context.Context, page, pageSize int) (interface{}, error) {
	path := fmt.Sprintf("document_types/?page=%d&page_size=%d", page, pageSize)
	var result interface{}
	err := c.Get(ctx, path, &result)
	return result, err
}

func (c *Client) GetStatistics(ctx context.Context) (interface{}, error) {
	path := "statistics/"
	var result interface{}
	err := c.Get(ctx, path, &result)
	return result, err
}

