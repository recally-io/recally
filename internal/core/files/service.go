package files

import (
	"context"
	"fmt"
	"io"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/s3"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

var DefaultService = &Service{
	dao: db.New(),
	s3:  s3.DefaultClient,
}

type Service struct {
	dao DAO
	s3  *s3.Client
}

func NewService(dao DAO) *Service {
	return &Service{
		dao: dao,
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

type UploadOptions struct {
	minio.PutObjectOptions
	FileType FileType
}

func (s *Service) UploadToS3(ctx context.Context, tx db.DBTX, originalURL, objectKey string, reader io.Reader, size int64, opts UploadOptions) (*DTO, error) {
	user, err := auth.LoadUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// check if file exist
	dto, err := s.GetFileByOriginalURL(ctx, tx, originalURL)
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
	info, err := s.s3.Upload(ctx, objectKey, reader, size, opts.PutObjectOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to S3: %w", err)
	}
	logger.FromContext(ctx).Info("file uploaded to s3", "key", objectKey, "size", info.Size)

	file := NewFile(
		user.ID,
		originalURL,
		objectKey,
		FileTypeImage,
		WithFileMetadata(Metadata{
			MIMEType: opts.ContentType,
			FileSize: size,
		}))
	return s.CreateFile(ctx, tx, file)
}
