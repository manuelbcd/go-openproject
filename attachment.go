package openproject

import (
	"context"
	"fmt"
	"io/ioutil"
)

// AttachmentService handles attachments for the OpenProject instance / API.
type AttachmentService struct {
	client *Client
}

// Attachment is the object representing OpenProject attachments.
// TODO: Complete fields and complex fields (user, links, downloadlocation, container...)
type Attachment struct {
	Type        string               `json:"_type,omitempty" structs:"_type,omitempty"`
	ID          int                  `json:"id,omitempty" structs:"id,omitempty"`
	FileName    string               `json:"filename,omitempty" structs:"filename,omitempty"`
	FileSize    int                  `json:"filesize,omitempty" structs:"filesize,omitempty"`
	Description OPGenericDescription `json:"description,omitempty" structs:"description,omitempty"`
	ContentType string               `json:"contentType,omitempty" structs:"contentType,omitempty"`
	Digest      AttachmentDigest     `json:"digest,omitempty" structs:"digest,omitempty"`
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
