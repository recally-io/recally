package files

import "github.com/minio/minio-go/v7"

type PutObjectOption func(*minio.PutObjectOptions)

func WithPutObjectOptionCacheControl(cc string) PutObjectOption {
	return func(opts *minio.PutObjectOptions) {
		opts.CacheControl = cc
	}
}

func WithPutObjectOptionContentType(ct string) PutObjectOption {
	return func(opts *minio.PutObjectOptions) {
		opts.ContentType = ct
	}
}

func NewPutObjectOptions(opts ...PutObjectOption) minio.PutObjectOptions {
	res := minio.PutObjectOptions{
		CacheControl: "max-age=31536000, public",
	}
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
