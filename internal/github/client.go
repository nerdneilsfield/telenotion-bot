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
	now        func() time.Time
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
		now:        time.Now,
	}
}

func (c *Client) UploadImage(ctx context.Context, data []byte, filename string) (string, error) {
	if c.repo == "" || c.branch == "" {
		return "", fmt.Errorf("github repo and branch are required")
	}

	hash := sha256.Sum256(data)
	hashHex := hex.EncodeToString(hash[:])
	if len(hashHex) > 8 {
		hashHex = hashHex[:8]
	}

	payload := createFileRequest{
		Message: "upload image",
		Content: base64.StdEncoding.EncodeToString(data),
		Branch:  c.branch,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	var lastErr error
	timestamp := ""
	if c.now != nil {
		timestamp = c.now().UTC().Format("20060102-150405")
	}
	for attempt := 0; attempt < 3; attempt++ {
		namedFile := applyHashSuffix(filename, hashHex, timestamp, attempt)
		path := joinPath(c.pathPrefix, namedFile)
		url := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", c.repo, path)

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
			_, _ = io.ReadAll(resp.Body)
			resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", c.repo, c.branch, path), nil
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

func applyHashSuffix(filename string, hash string, timestamp string, attempt int) string {
	ext := ""
	name := filename
	if idx := strings.LastIndex(filename, "."); idx > 0 {
		name = filename[:idx]
		ext = filename[idx:]
	}

	suffix := hash
	if timestamp != "" {
		suffix = fmt.Sprintf("%s-%s", suffix, timestamp)
	}
	if attempt > 0 {
		suffix = fmt.Sprintf("%s-%d", suffix, attempt)
	}

	if name == "" {
		return suffix + ext
	}

	return fmt.Sprintf("%s-%s%s", name, suffix, ext)
}
