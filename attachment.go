package openproject

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"
)

// AttachmentService handles attachments for the OpenProject instance / API.
type AttachmentService struct {
	client *Client
}

// Attachment is the object representing OpenProject attachments.
type Attachment struct {
	Embedded struct {
		Author struct {
			Type   string `json:"_type,omitempty"`
			ID     int    `json:"id,omitempty"`
			Name   string `json:"name,omitempty"`
			Avatar string `json:"avatar,omitempty"`
			Links  struct {
				Self OPGenericLink `json:"self,omitempty"`
			} `json:"_links,omitempty"`
		} `json:"author,omitempty"`
		Container struct {
			Type        string `json:"_type,omitempty"`
			ID          int    `json:"id,omitempty"`
			RowCount    int    `json:"rowCount,omitempty"`
			ColumnCount int    `json:"columnCount,omitempty"`
			Options     struct {
			} `json:"options,omitempty"`
			Widgets []struct {
				Type        string `json:"_type,omitempty"`
				ID          int    `json:"id,omitempty"`
				Identifier  string `json:"identifier,omitempty"`
				StartRow    int    `json:"startRow,omitempty"`
				EndRow      int    `json:"endRow,omitempty"`
				StartColumn int    `json:"startColumn,omitempty"`
				EndColumn   int    `json:"endColumn,omitempty"`
				Options     struct {
					Name string `json:"name,omitempty"`
					Text struct {
						Format string `json:"format,omitempty"`
						Raw    string `json:"raw,omitempty"`
						HTML   string `json:"html,omitempty"`
					} `json:"text,omitempty"`
					QueryID   string `json:"queryId,omitempty"`
					ChartType string `json:"chartType,omitempty"`
				} `json:"options,omitempty"`
			} `json:"widgets,omitempty"`
			CreatedAt time.Time `json:"createdAt,omitempty"`
			UpdatedAt time.Time `json:"updatedAt,omitempty"`
			Links     struct {
				Attachments OPGenericLink `json:"attachments,omitempty"`
				Scope       OPGenericLink `json:"scope,omitempty"`
				Self        OPGenericLink `json:"self,omitempty"`
			} `json:"_links,omitempty"`
		} `json:"container,omitempty"`
	} `json:"_embedded,omitempty"`
	Type        string `json:"_type,omitempty"`
	ID          int    `json:"id,omitempty"`
	FileName    string `json:"fileName,omitempty"`
	FileSize    int    `json:"fileSize,omitempty"`
	Description struct {
		Format string `json:"format,omitempty"`
		Raw    string `json:"raw,omitempty"`
		HTML   string `json:"html,omitempty"`
	} `json:"description,omitempty"`
	ContentType string           `json:"contentType,omitempty"`
	Digest      AttachmentDigest `json:"digest,omitempty"`
	CreatedAt   time.Time        `json:"createdAt,omitempty"`
	Links       struct {
		Self                   OPGenericLink `json:"self,omitempty"`
		Author                 OPGenericLink `json:"author,omitempty"`
		Container              OPGenericLink `json:"container,omitempty"`
		StaticDownloadLocation OPGenericLink `json:"staticDownloadLocation,omitempty"`
		DownloadLocation       OPGenericLink `json:"downloadLocation,omitempty"`
	} `json:"_links,omitempty"`
}

// AttachmentDigest wraps algorithm and hash
type AttachmentDigest struct {
	Algorithm string `json:"algorithm,omitempty" structs:"algorithm,omitempty"`
	Hash      string `json:"hash,omitempty" structs:"hash,omitempty"`
}

// GetWithContext gets a wiki page from OpenProject using its ID
func (s *AttachmentService) GetWithContext(ctx context.Context, attachmentID string) (*Attachment, *Response, error) {
	apiEndPoint := fmt.Sprintf("api/v3/attachments/%s", attachmentID)
	Obj, Resp, err := GetWithContext(ctx, s, apiEndPoint)
	return Obj.(*Attachment), Resp, err
}

// Get wraps GetWithContext using the background context.
func (s *AttachmentService) Get(attachmentID string) (*Attachment, *Response, error) {
	return s.GetWithContext(context.Background(), attachmentID)
}

// DownloadWithContext downloads a file from attachment using attachment ID
func (s *AttachmentService) DownloadWithContext(ctx context.Context, attachmentID string) (*[]byte, error) {
	apiEndpoint := fmt.Sprintf("api/v3/attachments/%s/content", attachmentID)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Download(req)

	if err != nil {
		return nil, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)

	return &respBytes, err
}

// Download wraps DownloadWithContext using the background context.
func (s *AttachmentService) Download(attachmentID string) (*[]byte, error) {
	return s.DownloadWithContext(context.Background(), attachmentID)
}
