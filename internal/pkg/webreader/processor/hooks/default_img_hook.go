package hooks

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/s3"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type ImageHook struct {
	s3Client *s3.Client
	pool     *db.Pool
}

type ImageHookOption func(h *ImageHook)

func WithImageHookS3ClientOption(cli *s3.Client) ImageHookOption {
	return func(h *ImageHook) {
		h.s3Client = cli
	}
}

func WithImageHookDBPoolOption(pool *db.Pool) ImageHookOption {
	return func(h *ImageHook) {
		h.pool = pool
	}
}

func NewImageHook(opts ...ImageHookOption) *ImageHook {
	h := &ImageHook{
		s3Client: s3.DefaultClient,
		pool:     db.DefaultPool,
	}

	for _, opt := range opts {
		opt(h)
	}
	return h

}

// func (h *ImageHook) ConvertToBase64(selec *goquery.Selection) {
// 	selec.Find("img").Each(func(i int, s *goquery.Selection) {
// 		src := s.AttrOr("src", "")
// 		if src == "" || strings.HasPrefix(src, "data:image") {
// 			return
// 		}
// 		img, _, contentType, err := h.loadImage(src)
// 		if err != nil {
// 			return
// 		}
// 		// Convert image to base64
// 		base64Str := fmt.Sprintf("data:%s;base64,%s", contentType, base64.StdEncoding.EncodeToString(img))
// 		s.SetAttr("src", base64Str)
// 	})
// }

func (h *ImageHook) Process(selec *goquery.Selection) {
	selec.Find("img").Each(func(i int, s *goquery.Selection) {
		src := s.AttrOr("src", "")

		if src == "" || strings.HasPrefix(src, "data:image") {
			return
		}
		// Upload image to S3 and get the public URL
		objectKey, err := h.UploadToS3(src)
		if err != nil {
			logger.Default.Error("failed to upload image to s3", "err", err)
			return
		}

		s.SetAttr("src", s3.DefaultClient.GetPublicURL(objectKey))
	})
}

func (h *ImageHook) UploadToS3(src string) (string, error) {
	// Validate and parse URL
	u, err := url.Parse(src)
	if err != nil {
		return "", err
	}
	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("invalid URL")
	}
	host := u.Host
	objectKey := fmt.Sprintf("images/%s/%s", host, uuid.New().String())

	// Asynchronously load and upload the image for better performance
	go func() {
		// Load image
		img, contentType, size, err := h.LoadImage(src)
		if err != nil {
			return
		}

		_, _ = h.UploadImage(img, objectKey, contentType, size)
	}()
	return s3.DefaultClient.GetPublicURL(objectKey), nil
}

func (h *ImageHook) UploadImage(img []byte, objectKey, contentType string, size int64) (string, error) {
	// Upload image to S3
	info, err := s3.DefaultClient.Upload(context.Background(), objectKey, bytes.NewReader(img), size, minio.PutObjectOptions{
		ContentType:  contentType,
		CacheControl: "max-age=31536000, public",
	})
	if err != nil {
		logger.Default.Error("failed to upload image to s3", "err", err, "objectKey", objectKey, "info", info)
		return "", err
	}
	logger.Default.Info("image uploaded to s3", "objectKey", objectKey)
	return s3.DefaultClient.GetPublicURL(objectKey), nil
}

func (h *ImageHook) LoadImage(uri string) (img []byte, contentType string, size int64, err error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Add user agent to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; RecallyBot/1.0)")

	// Perform request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", 0, fmt.Errorf("failed to download image: status %s", resp.Status)
	}

	// Read with max size limit (e.g., 10MB)
	const maxSize = 10 * 1024 * 1024
	img, err = io.ReadAll(io.LimitReader(resp.Body, maxSize))
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to read image: %w", err)
	}

	// Determine content type
	contentType = resp.Header.Get("Content-Type")
	if contentType == "" || !strings.HasPrefix(contentType, "image/") {
		contentType = http.DetectContentType(img)
		if !strings.HasPrefix(contentType, "image/") {
			return nil, "", 0, fmt.Errorf("invalid content type: %s", contentType)
		}
	}

	// Get size
	size = resp.ContentLength
	if size <= 0 {
		if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
			if parsedSize, err := strconv.ParseInt(contentLength, 10, 64); err == nil {
				size = parsedSize
			}
		}
	}
	if size <= 0 {
		size = int64(len(img))
	}
	logger.Default.Debug("image loaded", "contentType", contentType, "size", size, "src", uri)
	return img, contentType, size, nil
}
