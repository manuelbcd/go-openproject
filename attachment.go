package openproject

import (
	"context"
	"fmt"
)

/**
Attachment service handles attachments for the OpenProject instance / API.
*/
type AttachmentService struct {
	client *Client
}

/**
Attachment is the object representing OpenProject attachments.
// TODO: Complete fields and complex fields (user, links, downloadlocation, container...)
*/
type Attachment struct {
	Type        string               `json:"_type,omitempty" structs:"_type,omitempty"`
	Id          int                  `json:"id,omitempty" structs:"id,omitempty"`
	FileName    string               `json:"filename,omitempty" structs:"filename,omitempty"`
	FileSize    int                  `json:"filesize,omitempty" structs:"filesize,omitempty"`
	Description OPGenericDescription `json:"description,omitempty" structs:"description,omitempty"`
	ContentType string               `json:"contentType,omitempty" structs:"contentType,omitempty"`
	Digest      AttachmentDigest     `json:"digest,omitempty" structs:"digest,omitempty"`
}

type AttachmentDigest struct {
	Algorithm string `json:"algorithm,omitempty" structs:"algorithm,omitempty"`
	Hash      string `json:"hash,omitempty" structs:"hash,omitempty"`
}

/**
GetWithContext gets a wiki page from OpenProject using its ID
*/
func (s *AttachmentService) GetWithContext(ctx context.Context, attachmentID string) (*Attachment, *Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/attachments/%s", attachmentID)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	attachment := new(Attachment)
	resp, err := s.client.Do(req, attachment)
	if err != nil {
		return nil, resp, NewOpenProjectError(resp, err)
	}
	return attachment, resp, nil
}

/**
Get wraps GetWithContext using the background context.
*/
func (s *AttachmentService) Get(attachmentID string) (*Attachment, *Response, error) {
	return s.GetWithContext(context.Background(), attachmentID)
}
