package s3

import (
	"context"
	"fmt"
	"io"
	"time"
	"vibrain/internal/pkg/config"
	"vibrain/internal/pkg/logger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/cors"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	*minio.Client
	bucketName string
	publicURL  string
}

var corsRules = []cors.Rule{
	{
		AllowedHeader: []string{"*"},
		AllowedMethod: []string{"GET", "PUT"},
		AllowedOrigin: []string{"*"},
	},
}

func New(cfg config.S3Config) (*Client, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}
	c := &Client{Client: client, bucketName: cfg.BucketName, publicURL: cfg.PublicURL}
	err = c.PutBucketCors(context.Background())
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) Upload(ctx context.Context, objectKey string, reader io.Reader, size int64) (minio.UploadInfo, error) {
	info, err := c.PutObject(ctx, c.bucketName, objectKey, reader, size, minio.PutObjectOptions{})
	if err != nil {
		return minio.UploadInfo{}, fmt.Errorf("s3: failed to upload file: %w", err)
	}
	return info, nil
}

func (c *Client) GetPresignedURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	url, err := c.PresignedPutObject(ctx, c.bucketName, objectKey, expiry)
	if err != nil {
		return "", fmt.Errorf("s3: failed to get presigned put url: %w", err)
	}
	return url.String(), nil
}

func (c *Client) Delete(ctx context.Context, objectKey string) error {
	err := c.RemoveObject(ctx, c.bucketName, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("s3: failed to delete file: %w", err)
	}
	return nil
}

func (c *Client) GetPublicURL(objectKey string) string {
	return fmt.Sprintf("%s/%s/%s", c.publicURL, c.bucketName, objectKey)
}

// PutBucketCors sets the CORS configuration for the bucket if it's not already set.
// It first checks if CORS rules are already configured, and if not, applies the predefined rules.
//
// The function uses the corsRules defined at the package level, which allow all headers,
// GET and PUT methods from any origin.
//
// If an error occurs while getting or setting the CORS configuration, it returns a wrapped error.
func (c *Client) PutBucketCors(ctx context.Context) error {
	// Get the current CORS configuration for the bucket
	bucketCors, err := c.GetBucketCors(context.Background(), c.bucketName)
	if err == nil && bucketCors != nil && len(bucketCors.CORSRules) > 0 {
		return nil
	}

	logger.Default.Info("bucket cors not set, setting it", "bucket", c.bucketName)
	// Create a new CORS configuration using the predefined rules
	corsConfig := cors.NewConfig(corsRules)
	// Set the CORS configuration for the bucket
	err = c.SetBucketCors(ctx, c.bucketName, corsConfig)
	if err != nil {
		return fmt.Errorf("s3: failed to put bucket cors: %w", err)
	}

	return nil
}
