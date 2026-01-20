package session

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/nerdneilsfield/telenotion-bot/internal/config"
)

var maxImageSizeBytes int64 = 20 * 1024 * 1024

var (
	imageTooLargeMessage        = "Image skipped: exceeds 20MB limit."
	imageUnsupportedTypeMessage = "Image skipped: unsupported image type."
)

var (
	ErrImageTooLarge        = errors.New("image exceeds 20MB limit")
	ErrUnsupportedImageType = errors.New("unsupported image type")
)

var mediaHTTPClient = &http.Client{
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

var supportedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

func ConfigureMedia(cfg *config.Config) error {
	if cfg == nil {
		return nil
	}
	maxMB := cfg.Media.MaxImageSizeMB
	if maxMB <= 0 {
		return fmt.Errorf("media.max_image_size_mb must be positive")
	}
	maxImageSizeBytes = maxMB * 1024 * 1024
	imageTooLargeMessage = fmt.Sprintf("Image skipped: exceeds %dMB limit.", maxMB)

	allowed := make(map[string]string)
	allowedTypes := cfg.Media.AllowedImageTypes
	if len(allowedTypes) == 0 {
		allowedTypes = []string{"image/jpeg", "image/png", "image/gif", "image/webp"}
	}
	for _, value := range allowedTypes {
		key := strings.ToLower(strings.TrimSpace(value))
		if key == "" {
			continue
		}
		ext, ok := supportedImageTypes[key]
		if !ok {
			return fmt.Errorf("unsupported media type in config: %s", key)
		}
		allowed[key] = ext
	}
	if len(allowed) == 0 {
		return fmt.Errorf("media.allowed_image_types is required")
	}
	allowedImageTypes = allowed
	return nil
}

func downloadImage(ctx context.Context, url string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := mediaHTTPClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil, "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	if resp.ContentLength > maxImageSizeBytes {
		return nil, "", ErrImageTooLarge
	}

	limited := &io.LimitedReader{R: resp.Body, N: maxImageSizeBytes + 1}
	header := make([]byte, 512)
	n, err := io.ReadFull(limited, header)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		return nil, "", err
	}
	header = header[:n]
	if len(header) == 0 {
		return nil, "", ErrUnsupportedImageType
	}

	contentType := http.DetectContentType(header)
	ext, ok := allowedImageTypes[contentType]
	if !ok {
		return nil, "", ErrUnsupportedImageType
	}

	rest, err := io.ReadAll(limited)
	if err != nil {
		return nil, "", err
	}

	total := int64(len(header) + len(rest))
	if total > maxImageSizeBytes {
		return nil, "", ErrImageTooLarge
	}

	data := make([]byte, 0, len(header)+len(rest))
	data = append(data, header...)
	data = append(data, rest...)
	return data, ext, nil
}
