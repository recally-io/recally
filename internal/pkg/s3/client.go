package s3

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"recally/internal/pkg/config"
	"recally/internal/pkg/logger"
	"time"

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
		AllowedMethod: []string{"GET", "PUT", "POST", "DELETE", "HEAD"},
		AllowedOrigin: []string{"*"},
		ExposeHeader:  []string{"ETag"},
	},
}

var DefaultClient *Client

func init() {
	if config.Settings.S3.Enabled {
		if client, err := New(config.Settings.S3); err != nil {
			logger.Default.Error("failed to initialize s3 client", "err", err)
		} else {
			DefaultClient = client
		}
	}
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

	go func() {
		err = c.PutBucketCors(context.Background())
		if err != nil {
			logger.Default.Error("failed to put bucket cors", "err", err)
		} else {
			logger.Default.Info("bucket cors set successfully", "bucket", c.bucketName)
		}
	}()
	return c, nil
}

func (c *Client) PublicURL() string {
	return c.publicURL
}

func (c *Client) Upload(ctx context.Context, objectKey string, reader io.Reader, size int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	info, err := c.PutObject(ctx, c.bucketName, objectKey, reader, size, opts)
	if err != nil {
		return minio.UploadInfo{}, fmt.Errorf("s3: failed to upload file: %w", err)
	}
	return info, nil
}

func (c *Client) PresignedPutObject(ctx context.Context, objectName string, expires time.Duration) (*url.URL, error) {
	return c.Client.PresignedPutObject(ctx, c.bucketName, objectName, expires)
}

func (c *Client) PresignedHeadObject(ctx context.Context, objectName string, expires time.Duration, reqParams url.Values) (u *url.URL, err error) {
	return c.Client.PresignedHeadObject(ctx, c.bucketName, objectName, expires, reqParams)
}

func (c *Client) PresignedGetObject(ctx context.Context, objectName string, expires time.Duration, reqParams url.Values) (u *url.URL, err error) {
	if c.publicURL != "" {
		uri := fmt.Sprintf("%s/%s/%s", c.publicURL, c.bucketName, objectName)
		return url.Parse(uri)
	}

	return c.Client.PresignedGetObject(ctx, c.bucketName, objectName, expires, reqParams)
}

func (c *Client) Delete(ctx context.Context, objectKey string) error {
	err := c.RemoveObject(ctx, c.bucketName, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("s3: failed to delete file: %w", err)
	}
	return nil
}

func (c *Client) GetPublicURL(ctx context.Context, objectKey string) (string, error) {
	u, err := c.PresignedGetObject(ctx, objectKey, time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("s3: failed to get presigned get object: %w", err)
	}

	return u.String(), nil
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
		logger.Default.Debug("bucket cors already set", "bucket", c.bucketName, "cors", bucketCors.CORSRules)
		return nil
	}

	logger.Default.Debug("bucket cors not set, setting it", "bucket", c.bucketName)
	// Create a new CORS configuration using the predefined rules
	corsConfig := cors.NewConfig(corsRules)
	// Set the CORS configuration for the bucket
	err = c.SetBucketCors(ctx, c.bucketName, corsConfig)
	if err != nil {
		return fmt.Errorf("s3: failed to put bucket cors: %w", err)
	}

	return nil
}
