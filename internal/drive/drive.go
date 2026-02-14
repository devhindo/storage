package drive

import (
	"fmt"
	"net/http"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// Service wraps the Google Drive API service.
type Service struct {
	srv *drive.Service
}

// FileEntry represents a file or folder in Google Drive.
type FileEntry struct {
	ID       string
	Name     string
	MimeType string
	Size     int64
	IsFolder bool
}

// NewService creates a new Drive service from an authenticated HTTP client.
func NewService(client *http.Client) (*Service, error) {
	srv, err := drive.NewService(nil, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create drive service: %w", err)
	}
	return &Service{srv: srv}, nil
}

// ListFolder lists files and folders within the given folder ID.
// Use "root" for the top-level Drive folder.
func (s *Service) ListFolder(folderID string) ([]FileEntry, error) {
	query := fmt.Sprintf("'%s' in parents and trashed = false", folderID)

	var entries []FileEntry
	pageToken := ""

	for {
		call := s.srv.Files.List().
			Q(query).
			Fields("nextPageToken, files(id, name, mimeType, size)").
			OrderBy("folder, name").
			PageSize(100)

		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		result, err := call.Do()
		if err != nil {
			return nil, fmt.Errorf("unable to list files: %w", err)
		}

		for _, f := range result.Files {
			entries = append(entries, FileEntry{
				ID:       f.Id,
				Name:     f.Name,
				MimeType: f.MimeType,
				Size:     f.Size,
				IsFolder: f.MimeType == "application/vnd.google-apps.folder",
			})
		}

		pageToken = result.NextPageToken
		if pageToken == "" {
			break
		}
	}

	return entries, nil
}
