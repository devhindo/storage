package core

import "context"

// FileEntry represents a file or folder in a storage backend.
type FileEntry struct {
	ID       string
	Name     string
	MimeType string
	Size     int64
	IsFolder bool
}

// Backend defines the interface for a storage provider.
// Implementations include Google Drive, S3, local filesystem, etc.
type Backend interface {
	// ListFolder lists files and folders within the given folder ID.
	// Use "root" for the top-level folder.
	ListFolder(ctx context.Context, folderID string) ([]FileEntry, error)
}
