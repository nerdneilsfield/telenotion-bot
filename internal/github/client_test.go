package github

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestNewClient(t *testing.T) {
	client := NewClient("token", "owner/repo", "main", "images/")

	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.token != "token" {
		t.Errorf("token = %q, want %q", client.token, "token")
	}
	if client.repo != "owner/repo" {
		t.Errorf("repo = %q, want %q", client.repo, "owner/repo")
	}
	if client.branch != "main" {
		t.Errorf("branch = %q, want %q", client.branch, "main")
	}
	if client.pathPrefix != "images/" {
		t.Errorf("pathPrefix = %q, want %q", client.pathPrefix, "images/")
	}
}

func TestUploadImage_Success(t *testing.T) {
	data := []byte("test image data")
	hash := sha256.Sum256(data)
	hashHex := hex.EncodeToString(hash[:])
	if len(hashHex) > 8 {
		hashHex = hashHex[:8]
	}
	expectedPath := "images/test-" + hashHex + ".jpg"
	expectedURL := "https://raw.githubusercontent.com/owner/repo/main/" + expectedPath

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %s, want %s", r.Method, http.MethodPut)
		}
		if r.URL.Path != "/repos/owner/repo/contents/"+expectedPath {
			t.Errorf("Path = %s, want %s", r.URL.Path, "/repos/owner/repo/contents/"+expectedPath)
		}

		var req createFileRequest
		json.NewDecoder(r.Body).Decode(&req)

		auth := r.Header.Get("Authorization")
		if auth != "Bearer token" {
			t.Errorf("Authorization = %q, want %q", auth, "Bearer token")
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"content": map[string]interface{}{
				"download_url": expectedURL,
			},
		})
	}))
	defer server.Close()

	client := NewClient("token", "owner/repo", "main", "images/")
	client.client = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			req.URL.Scheme = "http"
			req.URL.Host = server.Listener.Addr().String()
			return http.DefaultTransport.RoundTrip(req)
		}),
	}

	url, err := client.UploadImage(context.Background(), data, "test.jpg")
	if err != nil {
		t.Fatalf("UploadImage() error = %v", err)
	}

	if url != expectedURL {
		t.Errorf("UploadImage() = %q, want %q", url, expectedURL)
	}
}

func TestUploadImage_RetryOnServerError(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"content": map[string]interface{}{
				"download_url": "https://raw.githubusercontent.com/owner/repo/main/test.jpg",
			},
		})
	}))
	defer server.Close()

	client := NewClient("token", "owner/repo", "main", "")
	client.client = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			req.URL.Scheme = "http"
			req.URL.Host = server.Listener.Addr().String()
			return http.DefaultTransport.RoundTrip(req)
		}),
	}

	data := []byte("test image data")
	url, err := client.UploadImage(context.Background(), data, "test.jpg")
	if err != nil {
		t.Fatalf("UploadImage() error = %v", err)
	}

	if attempts != 3 {
		t.Errorf("attempts = %d, want 3", attempts)
	}

	if url == "" {
		t.Error("UploadImage() returned empty URL")
	}
}

func TestUploadImage_NoRetryOnClientError(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient("token", "owner/repo", "main", "")
	client.client = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			req.URL.Scheme = "http"
			req.URL.Host = server.Listener.Addr().String()
			return http.DefaultTransport.RoundTrip(req)
		}),
	}

	data := []byte("test image data")
	_, err := client.UploadImage(context.Background(), data, "test.jpg")

	if err == nil {
		t.Error("UploadImage() should return error on 404")
	}

	if attempts != 1 {
		t.Errorf("attempts = %d, want 1 (no retry on 4xx)", attempts)
	}
}

func TestUploadImage_MissingRepoOrBranch(t *testing.T) {
	client := NewClient("token", "", "main", "")

	data := []byte("test image data")
	_, err := client.UploadImage(context.Background(), data, "test.jpg")

	if err == nil {
		t.Error("UploadImage() should return error when repo is empty")
	}
	if err.Error() != "github repo and branch are required" {
		t.Errorf("error = %q, want %q", err.Error(), "github repo and branch are required")
	}
}

func TestUploadImage_ContextCancelled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Long delay to ensure context is cancelled
		select {}
	}))
	defer server.Close()

	client := NewClient("token", "owner/repo", "main", "")
	client.client = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			req.URL.Scheme = "http"
			req.URL.Host = server.Listener.Addr().String()
			return http.DefaultTransport.RoundTrip(req)
		}),
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	data := []byte("test image data")
	_, err := client.UploadImage(ctx, data, "test.jpg")

	if err == nil {
		t.Error("UploadImage() should return error when context is cancelled")
	}
}

func TestJoinPath(t *testing.T) {
	tests := []struct {
		prefix string
		name   string
		want   string
	}{
		{"images/", "test.jpg", "images/test.jpg"},
		{"images", "test.jpg", "images/test.jpg"},
		{"", "test.jpg", "test.jpg"},
		{"images/", "", "images/"},
		{"  images/  ", "test.jpg", "images/test.jpg"},
	}

	for _, tt := range tests {
		t.Run(tt.prefix+"+"+tt.name, func(t *testing.T) {
			got := joinPath(tt.prefix, tt.name)
			if got != tt.want {
				t.Errorf("joinPath(%q, %q) = %q, want %q", tt.prefix, tt.name, got, tt.want)
			}
		})
	}
}

func TestUploadImage_ProperBase64Encoding(t *testing.T) {
	var receivedContent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req createFileRequest
		json.NewDecoder(r.Body).Decode(&req)
		receivedContent = req.Content
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}))
	defer server.Close()

	client := NewClient("token", "owner/repo", "main", "")
	client.client = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			req.URL.Scheme = "http"
			req.URL.Host = server.Listener.Addr().String()
			return http.DefaultTransport.RoundTrip(req)
		}),
	}

	// Use known content: "hello" base64 encoded is "aGVsbG8="
	data := []byte("hello")
	_, _ = client.UploadImage(context.Background(), data, "test.jpg")

	expectedContent := base64.StdEncoding.EncodeToString([]byte("hello"))
	if receivedContent != expectedContent {
		t.Errorf("UploadImage() sent content = %q, want %q", receivedContent, expectedContent)
	}
}
