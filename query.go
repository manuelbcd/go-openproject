package openproject

import (
	"context"
	"fmt"
	"math"
)

// QueryService handles statuses from the OpenProject instance / API.
type QueryService struct {
	client *Client
}

// SearchResultQuery represent a list of Projects
type SearchResultQuery struct {
	Embedded statusElements `json:"_embedded,omitempty" structs:"_embedded,omitempty"`
	PaginationParam
}

func (s *SearchResultQuery) TotalPage() int {
	return int(math.Ceil(float64(s.Total) / float64(s.PageSize)))
}

func (s *SearchResultQuery) ConcatEmbed(states interface{}) {
	s.Embedded.Elements = append(s.Embedded.Elements, states.(*SearchResultQuery).Embedded.Elements...)
}

// QueryElements array of elements within a query
type QueryElements struct {
	Elements []Status `json:"elements,omitempty" structs:"elements,omitempty"`
}

// Query is the object representing OpenProject queries.
// TODO: Complete fields (i.e. timelineVisible, showHierarchies, timelineZoomLevel, showHierarchies, etc...
type Query struct {
	Type             string          `json:"_type,omitempty" structs:"_type,omitempty"`
	Starred          bool            `json:"starred,omitempty" structs:"starred,omitempty"`
	ID               int             `json:"id,omitempty" structs:"id,omitempty"`
	Name             string          `json:"name,omitempty" structs:"name,omitempty"`
	Project          string          `json:"project,omitempty" structs:"columns,omitempty"`
	CreatedAt        *Time           `json:"createdAt,omitempty" structs:"createdAt,omitempty"`
	UpdatedAt        *Time           `json:"updatedAt,omitempty" structs:"updatedAt,omitempty"`
	Filters          []QueryFilter   `json:"filters,omitempty" structs:"filters,omitempty"`
	Sums             bool            `json:"sums,omitempty" structs:"sums,omitempty"`
	Public           bool            `json:"public,omitempty" structs:"public,omitempty"`
	Hidden           bool            `json:"hidden,omitempty" structs:"hidden,omitempty"`
	TimelineVisible  bool            `json:"timelineVisible,omitempty" structs:"timelineVisible,omitempty"`
	HighlightingMode string          `json:"highlightingMode,omitempty" structs:"highlightingMode,omitempty"`
	ShowHierarchies  bool            `json:"showHierarchies,omitempty" structs:"showHierarchies,omitempty"`
	Columns          []OPGenericLink `json:"columns,omitempty" structs:"columns,omitempty"`
	GroupBy          []OPGenericLink `json:"groupBy,omitempty" structs:"groupBy,omitempty"`
	SortBy           []OPGenericLink `json:"sortBy,omitempty" structs:"sortBy,omitempty"`
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
	if err != nil {
		return nil, Resp, err
	}
	return Obj.(*Query), Resp, err
}

// Get wraps GetWithContext using the background context.
func (s *QueryService) Get(queryID string) (*Query, *Response, error) {
	return s.GetWithContext(context.Background(), queryID)
}

// GetListWithContext Retrieve status list with context
// TODO: Implement search parameters-options
func (s *QueryService) GetListWithContext(ctx context.Context, offset int, pageSize int) (*SearchResultQuery, *Response, error) {
	apiEndpoint := "api/v3/queries"
	Obj, Resp, err := GetListWithContext(ctx, s, apiEndpoint, nil, offset, pageSize)
	if err != nil {
		return nil, Resp, err
	}
	return Obj.(*SearchResultQuery), Resp, err
}

// GetList wraps GetListWithContext using the background context.
func (s *QueryService) GetList(offset int, pageSize int) (*SearchResultQuery, *Response, error) {
	return s.GetListWithContext(context.Background(), offset, pageSize)
}

// CreateWithContext creates a query from a JSON representation.
func (s *QueryService) CreateWithContext(ctx context.Context, queryObj *Query) (*Query, *Response, error) {
	apiEndpoint := "api/v3/queries"
	wpResponse, resp, err := CreateWithContext(ctx, queryObj, s, apiEndpoint)
	if err != nil {
		return nil, resp, err
	}
	return wpResponse.(*Query), resp, err
}

// Create wraps CreateWithContext using the background context.
func (s *QueryService) Create(queryObj *Query) (*Query, *Response, error) {
	return s.CreateWithContext(context.Background(), queryObj)
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
