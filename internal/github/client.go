package github

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	token      string
	repo       string
	branch     string
	pathPrefix string
	client     *http.Client
}

type createFileRequest struct {
	Message string `json:"message"`
	Content string `json:"content"`
	Branch  string `json:"branch,omitempty"`
}

func NewClient(token, repo, branch, pathPrefix string) *Client {
	return &Client{
		token:      token,
		repo:       repo,
		branch:     branch,
		pathPrefix: pathPrefix,
		client:     http.DefaultClient,
	}
}

func (c *Client) UploadImage(ctx context.Context, data []byte, filename string) (string, error) {
	if c.repo == "" || c.branch == "" {
		return "", fmt.Errorf("github repo and branch are required")
	}

	path := joinPath(c.pathPrefix, filename)
	url := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", c.repo, path)

	payload := createFileRequest{
		Message: "upload image",
		Content: base64.StdEncoding.EncodeToString(data),
		Branch:  c.branch,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	_, _ = io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("upload failed with status %d", resp.StatusCode)
	}

	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", c.repo, c.branch, path), nil
}

func joinPath(prefix, name string) string {
	trimmed := strings.TrimSpace(prefix)
	if trimmed == "" {
		return name
	}
	if strings.HasSuffix(trimmed, "/") {
		return trimmed + name
	}
	return trimmed + "/" + name
}
