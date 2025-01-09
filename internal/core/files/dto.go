package files

import (
	"encoding/json"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Metadata struct {
	IsUploaded bool `json:"is_uploaded,omitempty"` // true if the file is uploaded to S3

	OriginalURL  string `json:"url,omitempty"`
	OriginalHost string `json:"host,omitempty"`

	Name     string `json:"name,omitempty"`
	Type     string `json:"type"`
	Ext      string `json:"ext,omitempty"`
	Size     int64  `json:"size,omitempty"`
	Hash     string `json:"hash,omitempty"`
	MIMEType string `json:"mime_type,omitempty"`
}

// DTO represents the domain model for a file
type DTO struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	OriginalURL string    `json:"original_url"`
	S3Key       string    `json:"s3_key"`
	S3URL       string    `json:"s3_url,omitempty"`
	FileName    string    `json:"file_name,omitempty"`
	FileType    string    `json:"file_type"`
	FileSize    int64     `json:"file_size,omitempty"`
	FileHash    string    `json:"file_hash,omitempty"`
	Metadata    Metadata  `json:"metadata,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Load converts a database object to a domain object
func (f *DTO) Load(dbo *db.File) {
	f.ID = dbo.ID
	f.UserID = dbo.UserID
	f.OriginalURL = dbo.OriginalUrl
	f.S3Key = dbo.S3Key
	f.S3URL = dbo.S3Url.String
	f.FileName = dbo.FileName.String
	f.FileType = dbo.FileType
	f.FileSize = dbo.FileSize.Int64
	f.FileHash = dbo.FileHash.String
	f.CreatedAt = dbo.CreatedAt.Time
	f.UpdatedAt = dbo.UpdatedAt.Time

	if dbo.Metadata != nil {
		if err := json.Unmarshal(dbo.Metadata, &f.Metadata); err != nil {
			logger.Default.Warn("failed to unmarshal File metadata",
				"err", err, "metadata", string(dbo.Metadata))
		}
	}
}

// Dump converts a domain object to database parameters for creation
func (f *DTO) Dump() db.CreateFileParams {
	metadata, _ := json.Marshal(f.Metadata)
	return db.CreateFileParams{
		UserID:      f.UserID,
		OriginalUrl: f.OriginalURL,
		S3Key:       f.S3Key,
		S3Url:       pgtype.Text{String: f.S3URL, Valid: f.S3URL != ""},
		FileName:    pgtype.Text{String: f.FileName, Valid: f.FileName != ""},
		FileType:    string(f.FileType),
		FileSize:    pgtype.Int8{Int64: f.FileSize, Valid: f.FileSize != 0},
		FileHash:    pgtype.Text{String: f.FileHash, Valid: f.FileHash != ""},
		Metadata:    metadata,
	}
}

// DumpToUpdateParams converts a domain object to database parameters for updating
func (f *DTO) DumpToUpdateParams() db.UpdateFileParams {
	metadata, _ := json.Marshal(f.Metadata)
	return db.UpdateFileParams{
		ID:       f.ID,
		S3Url:    pgtype.Text{String: f.S3URL, Valid: f.S3URL != ""},
		FileName: pgtype.Text{String: f.FileName, Valid: f.FileName != ""},
		FileType: string(f.FileType),
		FileSize: pgtype.Int8{Int64: f.FileSize, Valid: f.FileSize != 0},
		Metadata: metadata,
	}
}

// FileOption defines a function type for configuring FileDTO
type FileOption func(*DTO)

// NewFile creates a new FileDTO with the given options
func NewFile(userID uuid.UUID, originalURL string, s3Key string, fileType string, opts ...FileOption) *DTO {
	f := &DTO{
		ID:          uuid.New(),
		UserID:      userID,
		OriginalURL: originalURL,
		S3Key:       s3Key,
		FileType:    fileType,
	}

	for _, opt := range opts {
		opt(f)
	}
	return f
}

// Option functions for configuring a new file
func WithFileName(fileName string) FileOption {
	return func(f *DTO) {
		f.FileName = fileName
	}
}

func WithS3URL(s3URL string) FileOption {
	return func(f *DTO) {
		f.S3URL = s3URL
	}
}

func WithFileSize(size int64) FileOption {
	return func(f *DTO) {
		f.FileSize = size
	}
}

func WithFileHash(hash string) FileOption {
	return func(f *DTO) {
		f.FileHash = hash
	}
}

func WithFileMetadata(metadata Metadata) FileOption {
	return func(f *DTO) {
		f.Metadata = metadata
	}
}
