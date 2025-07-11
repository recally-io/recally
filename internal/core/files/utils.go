package files

import (
	"mime"
	"path/filepath"
)

// GetFileExtensionWithDefault returns the file extension of the given fileName.
// If the fileName has no extension, it returns the defaultExt.
func GetFileExtensionWithDefault(fileName, defaultExt string) string {
	ext := filepath.Ext(fileName)
	if ext == "" {
		return defaultExt
	}

	return ext
}

// GetFileMIMEWithDefault returns the MIME type of the given fileName based on its extension.
// If the MIME type cannot be determined, it returns the defaultMIME.
func GetFileMIMEWithDefault(fileName, defaultMIME string) string {
	mimeType := mime.TypeByExtension(filepath.Ext(fileName))
	if mimeType == "" {
		return defaultMIME
	}

	return mimeType
}
