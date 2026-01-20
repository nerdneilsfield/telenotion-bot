package github

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
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
		client:     newHTTPClient(),
	}
}

func (c *Client) UploadImage(ctx context.Context, data []byte, extension string) (string, error) {
	if c.repo == "" || c.branch == "" {
		return "", fmt.Errorf("github repo and branch are required")
	}

	ext := normalizeExtension(extension)
	hash := sha256.Sum256(data)
	hashHex := hex.EncodeToString(hash[:])
	filename := hashHex + ext
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

	rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", c.repo, c.branch, path)
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
		if err != nil {
			return "", err
		}

		req.Header.Set("Authorization", "Bearer "+c.token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = err
		} else {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return rawURL, nil
			}
			if resp.StatusCode == http.StatusUnprocessableEntity {
				return rawURL, nil
			}
			lastErr = fmt.Errorf("upload failed with status %d", resp.StatusCode)
			if resp.StatusCode < 500 && resp.StatusCode != http.StatusTooManyRequests {
				return "", lastErr
			}
		}

		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(time.Duration(attempt+1) * 500 * time.Millisecond):
		}
	}

	return "", lastErr
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
		},
	}
}

func normalizeExtension(extension string) string {
	ext := strings.TrimSpace(extension)
	if ext == "" {
		return ".bin"
	}
	if !strings.HasPrefix(ext, ".") {
		return "." + ext
	}
	return ext
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
