package core

import "context"

// FileService provides file operations using a pluggable storage backend.
// All frontends (CLI, TUI, web, etc.) should use this service
// rather than calling storage backends directly.
type FileService struct {
	backend Backend
}

// NewFileService creates a new FileService with the given storage backend.
func NewFileService(backend Backend) *FileService {
	return &FileService{backend: backend}
}

// ListFolder lists files and folders within the given folder ID.
func (s *FileService) ListFolder(ctx context.Context, folderID string) ([]FileEntry, error) {
	return s.backend.ListFolder(ctx, folderID)
}
