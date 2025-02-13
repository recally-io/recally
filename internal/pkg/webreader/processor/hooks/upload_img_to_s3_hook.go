package hooks

import (
	"context"
	"fmt"
	"recally/internal/core/files"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/s3"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

		s.SetAttr("src", files.DefaultService.GetShareURL(context.Background(), objectKey))
	})
}

func (h *ImageHook) UploadToS3(src string) (string, error) {
	ctx, user, err := auth.GetContextWithDummyUser(context.Background())
	if err != nil {
		return "", err
	}
	file, err := files.DefaultService.CreateFileAndUploadToS3FromUrl(ctx, h.pool.Pool, user.ID, true, h.host, src)
	if err != nil {
		return "", fmt.Errorf("failed to upload image to s3: %w", err)
	}
	return file.S3Key, nil
}
