package gdrive

import (
	"context"
	"fmt"
	"net/http"

	"github.com/devhindo/storage/pkg/core"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// Backend implements core.Backend using Google Drive.
type Backend struct {
	srv *drive.Service
}

// New creates a new Google Drive backend from an authenticated HTTP client.
func New(client *http.Client) (*Backend, error) {
	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create drive service: %w", err)
	}
	return &Backend{srv: srv}, nil
}

// ListFolder lists files and folders within the given folder ID.
func (b *Backend) ListFolder(ctx context.Context, folderID string) ([]core.FileEntry, error) {
	query := fmt.Sprintf("'%s' in parents and trashed = false", folderID)

	var entries []core.FileEntry
	pageToken := ""

	for {
		call := b.srv.Files.List().
			Context(ctx).
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
			entries = append(entries, core.FileEntry{
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
