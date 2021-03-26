package openproject

import (
	"context"
	"fmt"
)

/**
StatusService handles statuses from the OpenProject instance / API.
*/
type QueryService struct {
	client *Client
}

// StatusList represent a list of Projects
type searchResultQuery struct {
	Embedded statusElements `json:"_embedded,omitempty" structs:"_embedded,omitempty"`
	Total    int            `json:"total" structs:"total"`
	Count    int            `json:"count" structs:"count"`
	PageSize int            `json:"pageSize" structs:"pageSize"`
	Offset   int            `json:"offset" structs:"offset"`
}

type queryElements struct {
	Elements []Status `json:"elements,omitempty" structs:"elements,omitempty"`
}

/**
Query is the object representing OpenProject queries.
TODO: Complete fields (i.e. timelineVisible, showHierarchies, timelineZoomLevel, showHierarchies, etc...
*/
type Query struct {
	Type      string        `json:"_type,omitempty" structs:"_type,omitempty"`
	Starred   bool          `json:"starred,omitempty" structs:"starred,omitempty"`
	Id        int           `json:"id,omitempty" structs:"id,omitempty"`
	Name      string        `json:"name,omitempty" structs:"name,omitempty"`
	CreatedAt *Time         `json:"createdAt,omitempty" structs:"createdAt,omitempty"`
	UpdatedAt *Time         `json:"updatedAt,omitempty" structs:"updatedAt,omitempty"`
	Filters   []QueryFilter `json:"filters,omitempty" structs:"filters,omitempty"`
	Sums      bool          `json:"sums,omitempty" structs:"sums,omitempty"`
	Public    bool          `json:"public,omitempty" structs:"public,omitempty"`
	Hidden    bool          `json:"hidden,omitempty" structs:"hidden,omitempty"`
}

// TODO:Complete fields ( i.e. _links)
type QueryFilter struct {
	Type string `json:"_type,omitempty" structs:"_type,omitempty"`
	Name string `json:"name,omitempty" structs:"name,omitempty"`
}

/*
GetWithContext gets query info from OpenProject using its query Id
*/
func (s *QueryService) GetWithContext(ctx context.Context, queryID string) (*Query, *Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/queries/%s", queryID)
	Obj, Resp, err := GetWithContext(s, ctx, apiEndpoint)
	return Obj.(*Query), Resp, err
}

/**
Get wraps GetWithContext using the background context.
*/
func (s *QueryService) Get(queryID string) (*Query, *Response, error) {
	return s.GetWithContext(context.Background(), queryID)
}

/**
GetList wraps GetListWithContext using the background context.
*/
func (s *QueryService) GetList() (*searchResultQuery, *Response, error) {
	return s.GetListWithContext(context.Background())
}

/**
Retrieve status list with context
TODO: Implement search parameters-options
*/
func (s *QueryService) GetListWithContext(ctx context.Context) (*searchResultQuery, *Response, error) {
	apiEndpoint := "api/v3/queries"
	Obj, Resp, err := GetListWithContext(s, ctx, apiEndpoint, nil)
	return Obj.(*searchResultQuery), Resp, err
}
