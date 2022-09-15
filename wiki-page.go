package openproject

import (
	"context"
	"fmt"
)

// WikiPageService handles wiki-pages for the OpenProject instance / API.
type WikiPageService struct {
	client *Client
}

// WikiPage is the object representing OpenProject users.
// TODO: Complete fields and complex fields (_embedded.attachments, links, project, ...)
type WikiPage struct {
	Type     string       `json:"_type,omitempty" structs:"_type,omitempty"`
	ID       int          `json:"id,omitempty" structs:"id,omitempty"`
	Title    string       `json:"title,omitempty" structs:"title,omitempty"`
	Embedded WikiEmbedded `json:"_embedded,omitempty" structs:"_embedded,omitempty"`
}

// WikiEmbedded wraps embedded field of WikiPage
type WikiEmbedded struct {
	Project WikiProject `json:"project,omitempty" structs:"project,omitempty"`
}

// WikiProject wraps WikiEmbedded data
type WikiProject struct {
	Type       string `json:"_type,omitempty" structs:"_type,omitempty"`
	ID         int    `json:"id,omitempty" structs:"id,omitempty"`
	Identifier string `json:"identifier,omitempty" structs:"identifier,omitempty"`
	CreatedAt  *Time  `json:"createdAt,omitempty" structs:"createdAt,omitempty"`
	UpdatedAt  *Time  `json:"updatedAt,omitempty" structs:"updatedAt,omitempty"`
	Status     string `json:"status,omitempty" structs:"status,omitempty"`
}

// GetWithContext gets a wiki page from OpenProject using its ID
func (s *WikiPageService) GetWithContext(ctx context.Context, wikiID string) (*WikiPage, *Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/wiki_pages/%s", wikiID)
	Obj, Resp, err := GetWithContext(ctx, s, apiEndpoint)
	if err != nil {
		return nil, Resp, err
	}
	return Obj.(*WikiPage), Resp, err
}

// Get wraps GetWithContext using the background context.
func (s *WikiPageService) Get(wikiID string) (*WikiPage, *Response, error) {
	return s.GetWithContext(context.Background(), wikiID)
}
