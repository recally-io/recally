package hooks

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/s3"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type ImageHook struct{}

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

func (h *ImageHook) UploadToS3(selec *goquery.Selection) {
	selec.Find("img").Each(func(i int, s *goquery.Selection) {
		src := s.AttrOr("src", "")

		// Validate and parse URL
		u, err := url.Parse(src)
		if err != nil {
			return
		}
		if u.Scheme == "" || u.Host == "" {
			return
		}
		host := u.Host
		objectKey := fmt.Sprintf("images/%s/%s", host, uuid.New().String())

		// Asynchronously load and upload the image for better performance
		go func() {
			// Load image
			img, contentType, size, err := h.loadImage(src)
			if err != nil {
				return
			}
			logger.Default.Debug("image loaded", "host", host, "contentType", contentType, "size", size)
			// Upload image to S3
			info, err := s3.DefaultClient.Upload(context.Background(), objectKey, bytes.NewReader(img), size, minio.PutObjectOptions{
				ContentType:  contentType,
				CacheControl: "max-age=31536000, public",
			})
			if err != nil {
				logger.Default.Error("failed to upload image to s3", "err", err, "objectKey", objectKey, "info", info)
				return
			}
			logger.Default.Info("image uploaded to s3", "objectKey", objectKey)
		}()

		s.SetAttr("src", s3.DefaultClient.GetPublicURL(objectKey))
	})
}

func (h *ImageHook) loadImage(uri string) (img []byte, contentType string, size int64, err error) {
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

	return img, contentType, size, nil
}
