package openproject

import (
	"context"
	"fmt"
)

/**
ProjectService handles wiki-pages for the OpenProject instance / API.
*/
type WikiPageService struct {
	client *Client
}

/**
WikiPage is the object representing OpenProject users.
// TODO: Complete fields and complex fields (_embedded.attachments, links, project, ...)
*/
type WikiPage struct {
	Type     string       `json:"_type,omitempty" structs:"_type,omitempty"`
	Id       int          `json:"id,omitempty" structs:"id,omitempty"`
	Title    string       `json:"title,omitempty" structs:"title,omitempty"`
	Embedded WikiEmbedded `json:"_embedded,omitempty" structs:"_embedded,omitempty"`
}

type WikiEmbedded struct {
	Project WikiProject `json:"project,omitempty" structs:"project,omitempty"`
}

type WikiProject struct {
	Type       string `json:"_type,omitempty" structs:"_type,omitempty"`
	Id         int    `json:"id,omitempty" structs:"id,omitempty"`
	Identifier string `json:"identifier,omitempty" structs:"identifier,omitempty"`
	CreatedAt  *Time  `json:"createdAt,omitempty" structs:"createdAt,omitempty"`
	UpdatedAt  *Time  `json:"updatedAt,omitempty" structs:"updatedAt,omitempty"`
	Status     string `json:"status,omitempty" structs:"status,omitempty"`
}

/**
GetWithContext gets a wiki page from OpenProject using its ID
*/
func (s *WikiPageService) GetWithContext(ctx context.Context, wikiID string) (*WikiPage, *Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/wiki_pages/%s", wikiID)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	wikiPage := new(WikiPage)
	resp, err := s.client.Do(req, wikiPage)
	if err != nil {
		return nil, resp, NewOpenProjectError(resp, err)
	}
	return wikiPage, resp, nil
}

/**
Get wraps GetWithContext using the background context.
*/
func (s *WikiPageService) Get(wikiID string) (*WikiPage, *Response, error) {
	return s.GetWithContext(context.Background(), wikiID)
}
