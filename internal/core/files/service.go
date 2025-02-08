package files

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/config"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/s3"
	"recally/internal/pkg/session"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	utls "github.com/refraction-networking/utls"
)

var DefaultService = NewService(s3.DefaultClient)

type Service struct {
	dao DAO
	s3  *s3.Client
}

func NewService(s3 *s3.Client) *Service {
	return &Service{
		dao: db.New(),
		s3:  s3,
	}
}

func (s *Service) CreateFile(ctx context.Context, tx db.DBTX, file *DTO) (*DTO, error) {
	_, err := auth.LoadUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dbo, err := s.dao.CreateFile(ctx, tx, file.Dump())
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	file.Load(&dbo)
	return file, nil
}

func (s *Service) GetFile(ctx context.Context, tx db.DBTX, id uuid.UUID) (*DTO, error) {
	_, err := auth.LoadUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dbo, err := s.dao.GetFileByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	var file DTO
	file.Load(&dbo)
	return &file, nil
}

func (s *Service) GetFileByS3Key(ctx context.Context, tx db.DBTX, s3Key string) (*DTO, error) {
	user, err := auth.LoadUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dbo, err := s.dao.GetFileByS3Key(ctx, tx, db.GetFileByS3KeyParams{
		S3Key:  s3Key,
		UserID: user.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file by S3 key: %w", err)
	}

	var file DTO
	file.Load(&dbo)
	return &file, nil
}

func (s *Service) GetFileByOriginalURL(ctx context.Context, tx db.DBTX, originalUrl string) (*DTO, error) {
	user, err := auth.LoadUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dbo, err := s.dao.GetFileByOriginalURL(ctx, tx, db.GetFileByOriginalURLParams{
		OriginalUrl: originalUrl,
		UserID:      user.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file by original URL: %w", err)
	}

	var file DTO
	file.Load(&dbo)
	return &file, nil
}

func (s *Service) DeleteFile(ctx context.Context, tx db.DBTX, id uuid.UUID) error {
	_, err := auth.LoadUserFromContext(ctx)
	if err != nil {
		return err
	}

	if err := s.dao.DeleteFile(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (s *Service) UploadToS3(ctx context.Context, tx db.DBTX, objectKey string, reader io.Reader, metadata *Metadata, opts minio.PutObjectOptions) (*DTO, error) {
	user, err := auth.LoadUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// check if file exist
	dto, err := s.GetFileByOriginalURL(ctx, tx, metadata.OriginalURL)
	if err != nil && !db.IsNotFoundError(err) {
		return nil, fmt.Errorf("failed to get file by original URL: %w", err)
	}
	if dto != nil {
		return dto, nil
	}

	// upload file to s3 if not exist
	if objectKey == "" {
		objectKey = fmt.Sprintf("%s/%s", user.ID.String(), uuid.New().String())
	}

	if opts.ContentType == "" {
		opts.ContentType = metadata.MIMEType
	}
	info, err := s.s3.Upload(ctx, objectKey, reader, metadata.Size, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to S3: %w", err)
	}
	logger.FromContext(ctx).Info("file uploaded to s3", "key", objectKey, "size", info.Size)
	metadata.IsUploaded = true

	file := NewFile(
		user.ID,
		metadata.OriginalURL,
		objectKey,
		metadata.Type,
		WithFileMetadata(*metadata))
	return s.CreateFile(ctx, tx, file)
}

func (s *Service) UploadToS3FromUrl(ctx context.Context, tx db.DBTX, async bool, host, uri string, opts minio.PutObjectOptions) (*DTO, error) {
	// Validate and parse URL
	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid content url: %s, %w", uri, err)
	}

	// get absolute URL
	if u.Host == "" {
		if host == "" {
			return nil, fmt.Errorf("invalid content url: %s", uri)
		}
		u.Host = host
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}

	uri = u.String()

	file, err := s.GetFileByOriginalURL(ctx, tx, uri)
	if err != nil && !db.IsNotFoundError(err) {
		return nil, fmt.Errorf("failed to get file by original URL: %w", err)
	}
	// if image already exists, return the URL
	if file != nil { // file.S3Key is empty if the file is not uploaded to s3
		logger.FromContext(ctx).Debug("file already exists", "url", uri, "objectKey", file.S3Key)
		return file, nil
	}

	user, err := auth.LoadUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	fileExt := path.Ext(uri)
	objectKey := fmt.Sprintf("%s/%d/%d/%s/%s%s", user.ID.String(), now.Year(), now.Month(), u.Host, uuid.New().String(), fileExt)

	upload := func(ctx context.Context) (*DTO, error) {
		contentReader, metadata, err := s.loadContent(ctx, host, uri, fileExt)
		if err != nil {
			return nil, fmt.Errorf("failed to load content: %w", err)
		}
		return s.UploadToS3(ctx, tx, objectKey, contentReader, metadata, opts)
	}

	if async {
		go func() {
			_, err := upload(auth.SetUserToContext(context.Background(), user))
			if err != nil {
				logger.Default.Error("failed to upload file to s3, save the original url", "err", err, "url", uri)
				file := NewFile(user.ID, uri, objectKey, "unknown", WithFileMetadata(Metadata{
					IsUploaded: false,
				}))
				if _, err = s.CreateFile(ctx, tx, file); err != nil {
					logger.Default.Error("failed to save the original url", "err", err, "url", uri)
				}
			}
		}()
		return &DTO{
			S3Key: objectKey,
		}, nil
	}
	return upload(ctx)
}

func (s *Service) loadContent(ctx context.Context, host, uri, ext string) (io.ReadCloser, *Metadata, error) {
	metadata := &Metadata{
		OriginalURL:  uri,
		OriginalHost: host,
		Ext:          ext,
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add user agent to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; RecallyBot/1.0)")

	// Perform request
	sess := session.New(session.WithClientHelloID(utls.HelloChrome_100_PSK), session.WithTimeout(30*time.Second))
	resp, err := sess.Client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to download content: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("failed to download content: status %s", resp.Status)
	}

	// Read with max size limit (e.g., 50MB)
	const maxSize = 50 * 1024 * 1024
	content, err := io.ReadAll(io.LimitReader(resp.Body, maxSize))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read content: %w", err)
	}

	// Determine content type
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(content)
	}
	metadata.MIMEType = contentType

	if contentType != "" {
		metadata.Type = strings.Split(contentType, "/")[0]
	}

	// Get size
	size := resp.ContentLength
	if size <= 0 {
		if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
			if parsedSize, err := strconv.ParseInt(contentLength, 10, 64); err == nil {
				size = parsedSize
			}
		}
	}
	if size <= 0 {
		size = int64(len(content))
	}
	metadata.Size = size

	return io.NopCloser(bytes.NewReader(content)), metadata, nil
}

func (s *Service) GetPresignedGetObjectURL(ctx context.Context, tx db.DBTX, objectKey string, expires time.Duration, reqParams url.Values) (string, error) {
	file, err := s.GetFileByS3Key(ctx, tx, objectKey)
	if err != nil {
		return "", err
	}
	if !file.Metadata.IsUploaded {
		return file.OriginalURL, nil
	}
	u, err := s.s3.PresignedGetObject(ctx, objectKey, expires, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to get presigned get URL: %w", err)
	}
	return u.String(), nil
}

func (s *Service) GetPresignedHeadObjectURL(ctx context.Context, tx db.DBTX, objectKey string, expires time.Duration, reqParams url.Values) (string, error) {
	file, err := s.GetFileByS3Key(ctx, tx, objectKey)
	if err != nil {
		return "", err
	}
	if !file.Metadata.IsUploaded {
		return file.OriginalURL, nil
	}
	u, err := s.s3.PresignedHeadObject(ctx, objectKey, expires, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to get presigned head URL: %w", err)
	}
	return u.String(), nil
}

func (s *Service) GetPresignedPutObjectURL(ctx context.Context, objectKey string, expires time.Duration) (string, error) {
	u, err := s.s3.PresignedPutObject(ctx, objectKey, expires)
	if err != nil {
		return "", fmt.Errorf("failed to get presigned head URL: %w", err)
	}
	return u.String(), nil
}

// GetPublicURL returns the public URL of the file
func (s *Service) GetPublicURL(ctx context.Context, objectKey string) (string, error) {
	return s.s3.GetPublicURL(ctx, objectKey)
}

// GetShareURL returns the share URL of the file by proxy API
func (s *Service) GetShareURL(ctx context.Context, objectKey string) string {
	return fmt.Sprintf("%s/api/v1/shared/files/%s", config.Settings.Service.Fqdn, objectKey)
}
