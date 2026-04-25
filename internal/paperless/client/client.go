package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
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
			Timeout: 60 * time.Second,
		},
	}
}

func (c *Client) do(ctx context.Context, method, path string, body io.Reader, contentType string) (*http.Response, error) {
	url := c.BaseURL + strings.TrimPrefix(path, "/")
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Token "+c.APIToken)
	req.Header.Set("Accept", "application/json")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	} else if body != nil {
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

func (c *Client) Post(ctx context.Context, path string, body interface{}, target interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = strings.NewReader(string(data))
	}

	resp, err := c.do(ctx, http.MethodPost, path, bodyReader, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if target != nil {
		return json.NewDecoder(resp.Body).Decode(target)
	}
	return nil
}

func (c *Client) Patch(ctx context.Context, path string, body interface{}, target interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = strings.NewReader(string(data))
	}

	resp, err := c.do(ctx, http.MethodPatch, path, bodyReader, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if target != nil {
		return json.NewDecoder(resp.Body).Decode(target)
	}
	return nil
}

func (c *Client) Delete(ctx context.Context, path string) error {
	resp, err := c.do(ctx, http.MethodDelete, path, nil, "")
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) UploadDocument(ctx context.Context, fileName string, fileContent io.Reader, metadata map[string]string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("document", fileName)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, fileContent); err != nil {
		return "", err
	}

	for key, val := range metadata {
		if err := writer.WriteField(key, val); err != nil {
			return "", err
		}
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	resp, err := c.do(ctx, http.MethodPost, "documents/post_document/", body, writer.FormDataContentType())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (c *Client) DownloadDocument(ctx context.Context, id int, mode string) ([]byte, string, error) {
	var path string
	switch mode {
	case "preview":
		path = fmt.Sprintf("documents/%d/preview/", id)
	case "thumbnail":
		path = fmt.Sprintf("documents/%d/thumbnail/", id)
	default:
		path = fmt.Sprintf("documents/%d/download/", id)
	}

	resp, err := c.do(ctx, http.MethodGet, path, nil, "")
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return content, resp.Header.Get("Content-Type"), nil
}

func (c *Client) DeleteDocument(ctx context.Context, id int) error {
	return c.Delete(ctx, fmt.Sprintf("documents/%d/", id))
}

func (c *Client) BulkEditDocuments(ctx context.Context, documentIDs []int, method string, args map[string]interface{}) error {
	body := map[string]interface{}{
		"documents": documentIDs,
		"method":    method,
		"parameters":   args,
	}
	return c.Post(ctx, "documents/bulk_edit/", body, nil)
}

func (c *Client) CreateTag(ctx context.Context, name string) (interface{}, error) {
	body := map[string]interface{}{
		"name":               name,
		"matching_algorithm": 0, // None/Manual
		"is_inbox_tag":       false,
	}
	var result interface{}
	err := c.Post(ctx, "tags/", body, &result)
	return result, err
}

func (c *Client) UpdateTag(ctx context.Context, id int, updates map[string]interface{}) (interface{}, error) {
	path := fmt.Sprintf("tags/%d/", id)
	var result interface{}
	err := c.Patch(ctx, path, updates, &result)
	return result, err
}

func (c *Client) DeleteTag(ctx context.Context, id int) error {
	return c.Delete(ctx, fmt.Sprintf("tags/%d/", id))
}

func (c *Client) CreateCorrespondent(ctx context.Context, name string) (interface{}, error) {
	body := map[string]interface{}{
		"name":               name,
		"matching_algorithm": 0, // None/Manual
	}
	var result interface{}
	err := c.Post(ctx, "correspondents/", body, &result)
	return result, err
}

func (c *Client) UpdateCorrespondent(ctx context.Context, id int, updates map[string]interface{}) (interface{}, error) {
	path := fmt.Sprintf("correspondents/%d/", id)
	var result interface{}
	err := c.Patch(ctx, path, updates, &result)
	return result, err
}

func (c *Client) DeleteCorrespondent(ctx context.Context, id int) error {
	return c.Delete(ctx, fmt.Sprintf("correspondents/%d/", id))
}

func (c *Client) CreateDocumentType(ctx context.Context, name string) (interface{}, error) {
	body := map[string]interface{}{
		"name":               name,
		"matching_algorithm": 0, // None/Manual
	}
	var result interface{}
	err := c.Post(ctx, "document_types/", body, &result)
	return result, err
}

func (c *Client) UpdateDocumentType(ctx context.Context, id int, updates map[string]interface{}) (interface{}, error) {
	path := fmt.Sprintf("document_types/%d/", id)
	var result interface{}
	err := c.Patch(ctx, path, updates, &result)
	return result, err
}

func (c *Client) DeleteDocumentType(ctx context.Context, id int) error {
	return c.Delete(ctx, fmt.Sprintf("document_types/%d/", id))
}

func (c *Client) ListStoragePaths(ctx context.Context, page, pageSize int) (interface{}, error) {
	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("page_size", fmt.Sprintf("%d", pageSize))

	path := "storage_paths/?" + params.Encode()
	var result interface{}
	err := c.Get(ctx, path, &result)
	return result, err
}

func (c *Client) CreateStoragePath(ctx context.Context, name, pathStr string) (interface{}, error) {
	body := map[string]interface{}{
		"name":               name,
		"path":               pathStr,
		"matching_algorithm": 0, // None/Manual
	}
	var result interface{}
	err := c.Post(ctx, "storage_paths/", body, &result)
	return result, err
}

func (c *Client) UpdateStoragePath(ctx context.Context, id int, updates map[string]interface{}) (interface{}, error) {
	path := fmt.Sprintf("storage_paths/%d/", id)
	var result interface{}
	err := c.Patch(ctx, path, updates, &result)
	return result, err
}

func (c *Client) DeleteStoragePath(ctx context.Context, id int) error {
	return c.Delete(ctx, fmt.Sprintf("storage_paths/%d/", id))
}

func (c *Client) UpdateDocument(ctx context.Context, id int, updates map[string]interface{}) (interface{}, error) {
	path := fmt.Sprintf("documents/%d/", id)
	var result interface{}
	err := c.Patch(ctx, path, updates, &result)
	return result, err
}

func (c *Client) Get(ctx context.Context, path string, target interface{}) error {
	resp, err := c.do(ctx, http.MethodGet, path, nil, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func (c *Client) SearchDocuments(ctx context.Context, query string, page, pageSize int) (interface{}, error) {
	params := url.Values{}
	params.Add("query", query)
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("page_size", fmt.Sprintf("%d", pageSize))

	path := "documents/?" + params.Encode()
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
	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("page_size", fmt.Sprintf("%d", pageSize))

	path := "tags/?" + params.Encode()
	var result interface{}
	err := c.Get(ctx, path, &result)
	return result, err
}

func (c *Client) ListCorrespondents(ctx context.Context, page, pageSize int) (interface{}, error) {
	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("page_size", fmt.Sprintf("%d", pageSize))

	path := "correspondents/?" + params.Encode()
	var result interface{}
	err := c.Get(ctx, path, &result)
	return result, err
}

func (c *Client) ListDocumentTypes(ctx context.Context, page, pageSize int) (interface{}, error) {
	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("page_size", fmt.Sprintf("%d", pageSize))

	path := "document_types/?" + params.Encode()
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

func (c *Client) ListTasks(ctx context.Context) (interface{}, error) {
	path := "tasks/"
	var result interface{}
	err := c.Get(ctx, path, &result)
	return result, err
}

func (c *Client) ListCustomFields(ctx context.Context, page, pageSize int) (interface{}, error) {
	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("page_size", fmt.Sprintf("%d", pageSize))

	path := "custom_fields/?" + params.Encode()
	var result interface{}
	err := c.Get(ctx, path, &result)
	return result, err
}

func (c *Client) CreateCustomField(ctx context.Context, name, dataType string) (interface{}, error) {
	body := map[string]interface{}{
		"name":      name,
		"data_type": dataType,
	}
	var result interface{}
	err := c.Post(ctx, "custom_fields/", body, &result)
	return result, err
}

func (c *Client) UpdateCustomField(ctx context.Context, id int, updates map[string]interface{}) (interface{}, error) {
	path := fmt.Sprintf("custom_fields/%d/", id)
	var result interface{}
	err := c.Patch(ctx, path, updates, &result)
	return result, err
}

func (c *Client) DeleteCustomField(ctx context.Context, id int) error {
	return c.Delete(ctx, fmt.Sprintf("custom_fields/%d/", id))
}

func (c *Client) ListLogs(ctx context.Context) (interface{}, error) {
	path := "logs/"
	var result interface{}
	err := c.Get(ctx, path, &result)
	return result, err
}

func (c *Client) GetLog(ctx context.Context, logName string) (string, error) {
	path := fmt.Sprintf("logs/%s/", logName)
	resp, err := c.do(ctx, http.MethodGet, path, nil, "")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c *Client) ListShareLinks(ctx context.Context, page, pageSize int) (interface{}, error) {
	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("page_size", fmt.Sprintf("%d", pageSize))

	path := "share_links/?" + params.Encode()
	var result interface{}
	err := c.Get(ctx, path, &result)
	return result, err
}

func (c *Client) CreateShareLink(ctx context.Context, documentID int, expires string) (interface{}, error) {
	body := map[string]interface{}{
		"document": documentID,
	}
	if expires != "" {
		body["expiration"] = expires
	}
	var result interface{}
	err := c.Post(ctx, "share_links/", body, &result)
	return result, err
}

func (c *Client) DeleteShareLink(ctx context.Context, id int) error {
	return c.Delete(ctx, fmt.Sprintf("share_links/%d/", id))
}
