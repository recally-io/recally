package hooks

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"recally/internal/core/files"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/s3"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/minio/minio-go/v7"
)

type ImageHook struct {
	host     string
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

func NewImageHook(host string, opts ...ImageHookOption) *ImageHook {
	h := &ImageHook{
		host:     host,
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
		return "", fmt.Errorf("invalid image original source url: %s, %w", src, err)
	}

	// get absolute URL
	if u.Scheme == "" {
		u.Scheme = "http"
	}

	if u.Host == "" {
		u.Host = h.host
	}

	ctx, err := auth.GetContextWithDummyUser(context.Background())
	if err != nil {
		return "", err
	}
	file, err := files.DefaultService.GetFileByOriginalURL(ctx, h.pool.Pool, src)
	if err != nil && !db.IsNotFoundError(err) {
		return "", fmt.Errorf("failed to get file by original URL: %w", err)
	}

	// if image already exists, return the URL
	if file != nil { // file.S3Key is empty if the file is not uploaded to s3
		logger.Default.Info("file already exists", "url", src, "objectKey", file.S3Key)
		return file.S3Key, nil
	}

	// Generate a unique object key
	dummyUser, _ := auth.LoadUserFromContext(ctx)
	objectKey := fmt.Sprintf("%s/images/%s/%s", dummyUser.ID.String(), u.Host, uuid.New().String())
	imgType := path.Ext(u.Path)
	if imgType != "" {
		objectKey += imgType
	}

	// Asynchronously load and upload the image for better performance
	go func() {
		// Load image
		img, contentType, size, err := h.LoadImage(src)
		if err != nil {
			return
		}
		// save file metadata to database
		if err := db.RunInTransaction(ctx, h.pool.Pool, func(ctx context.Context, tx pgx.Tx) error {
			_, err = files.DefaultService.UploadToS3(ctx, tx, src, objectKey, bytes.NewReader(img), size, files.UploadOptions{
				PutObjectOptions: minio.PutObjectOptions{
					ContentType:  contentType,
					CacheControl: "max-age=31536000, public",
				},
				FileType: files.FileTypeImage,
			})
			return err
		}); err != nil {
			logger.FromContext(ctx).Error("failed to save file metadata to database", "err", err)
			return
		}
		logger.FromContext(ctx).Info("file metadata saved to database", "objectKey", objectKey)
	}()
	return objectKey, nil
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
