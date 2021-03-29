package openproject

import (
	"context"
	"fmt"
)

// QueryService handles statuses from the OpenProject instance / API.
type QueryService struct {
	client *Client
}

// SearchResultQuery represent a list of Projects
type SearchResultQuery struct {
	Embedded statusElements `json:"_embedded,omitempty" structs:"_embedded,omitempty"`
	Total    int            `json:"total" structs:"total"`
	Count    int            `json:"count" structs:"count"`
	PageSize int            `json:"pageSize" structs:"pageSize"`
	Offset   int            `json:"offset" structs:"offset"`
}

// QueryElements array of elements within a query
type QueryElements struct {
	Elements []Status `json:"elements,omitempty" structs:"elements,omitempty"`
}

// Query is the object representing OpenProject queries.
// TODO: Complete fields (i.e. timelineVisible, showHierarchies, timelineZoomLevel, showHierarchies, etc...
type Query struct {
	Type      string        `json:"_type,omitempty" structs:"_type,omitempty"`
	Starred   bool          `json:"starred,omitempty" structs:"starred,omitempty"`
	ID        int           `json:"id,omitempty" structs:"id,omitempty"`
	Name      string        `json:"name,omitempty" structs:"name,omitempty"`
	CreatedAt *Time         `json:"createdAt,omitempty" structs:"createdAt,omitempty"`
	UpdatedAt *Time         `json:"updatedAt,omitempty" structs:"updatedAt,omitempty"`
	Filters   []QueryFilter `json:"filters,omitempty" structs:"filters,omitempty"`
	Sums      bool          `json:"sums,omitempty" structs:"sums,omitempty"`
	Public    bool          `json:"public,omitempty" structs:"public,omitempty"`
	Hidden    bool          `json:"hidden,omitempty" structs:"hidden,omitempty"`
}

// QueryFilter filters within a query
// TODO:Complete fields ( i.e. _links)
type QueryFilter struct {
	Type string `json:"_type,omitempty" structs:"_type,omitempty"`
	Name string `json:"name,omitempty" structs:"name,omitempty"`
}

// GetWithContext gets query info from OpenProject using its query ID
func (s *QueryService) GetWithContext(ctx context.Context, queryID string) (*Query, *Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/queries/%s", queryID)
	Obj, Resp, err := GetWithContext(ctx, s, apiEndpoint)
	return Obj.(*Query), Resp, err
}

// Get wraps GetWithContext using the background context.
func (s *QueryService) Get(queryID string) (*Query, *Response, error) {
	return s.GetWithContext(context.Background(), queryID)
}

// GetList wraps GetListWithContext using the background context.
func (s *QueryService) GetList() (*SearchResultQuery, *Response, error) {
	return s.GetListWithContext(context.Background())
}

// GetListWithContext Retrieve status list with context
// TODO: Implement search parameters-options
func (s *QueryService) GetListWithContext(ctx context.Context) (*SearchResultQuery, *Response, error) {
	apiEndpoint := "api/v3/queries"
	Obj, Resp, err := GetListWithContext(ctx, s, apiEndpoint, nil)
	return Obj.(*SearchResultQuery), Resp, err
}

// DeleteWithContext will delete a single query object
func (s *QueryService) DeleteWithContext(ctx context.Context, queryID string) (*Response, error) {
	apiEndPoint := fmt.Sprintf("api/v3/queries/%s", queryID)
	resp, err := DeleteWithContext(ctx, s, apiEndPoint)
	return resp, err
}

// Delete wraps DeleteWithContext using the background context.
func (s *QueryService) Delete(queryID string) (*Response, error) {
	return s.DeleteWithContext(context.Background(), queryID)
}
