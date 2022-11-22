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

// QueryResult represents a query. Use for single query result.
type QueryResult struct {
	Type               string        `json:"_type,omitempty" structs:"_type,omitempty"`
	Starred            bool          `json:"starred,omitempty" structs:"starred,omitempty"`
	ID                 int           `json:"id,omitempty" structs:"id,omitempty"`
	Name               string        `json:"name,omitempty" structs:"name,omitempty"`
	Project            string        `json:"project,omitempty" structs:"columns,omitempty"`
	CreatedAt          *Time         `json:"createdAt,omitempty" structs:"createdAt,omitempty"`
	UpdatedAt          *Time         `json:"updatedAt,omitempty" structs:"updatedAt,omitempty"`
	Filters            []QueryFilter `json:"filters,omitempty" structs:"filters,omitempty"`
	IncludeSubprojects bool          `json:"includeSubprojects,omitempty" structs:"includeSubprojects,omitempty"`
	Sums               bool          `json:"sums,omitempty" structs:"sums,omitempty"`
	Public             bool          `json:"public,omitempty" structs:"public,omitempty"`
	Hidden             bool          `json:"hidden,omitempty" structs:"hidden,omitempty"`
	TimelineVisible    bool          `json:"timelineVisible,omitempty" structs:"timelineVisible,omitempty"`
	HighlightingMode   string        `json:"highlightingMode,omitempty" structs:"highlightingMode,omitempty"`
	ShowHierarchies    bool          `json:"showHierarchies,omitempty" structs:"showHierarchies,omitempty"`
	TimelineZoomLevel  string        `json:"timelineZoomLevel,omitempty" structs:"timelineZoomLevel,omitempty"`
	Embedded           struct {
		Project               *Project              `json:"project,omitempty" structs:"project,omitempty"`
		User                  *User                 `json:"user,omitempty" structs:"user,omitempty"`
		SortBy                []QuerySortBy         `json:"sortBy,omitempty" structs:"sortBy,omitempty"`
		Columns               []QueryEmbeddedColumn `json:"columns,omitempty" structs:"columns,omitempty"`
		HighlightedAttributes []QueryEmbeddedColumn `json:"highlightedAttributes,omitempty" structs:"highlightedAttributes,omitempty"`
		Results               struct {
			PaginationParam
			Type     string `json:"_type,omitempty" structs:"_type,omitempty"`
			Embedded struct {
				Elements []WorkPackage `json:"elements,omitempty" structs:"elements,omitempty"`
				Schemas  Schemas       `json:"schemas,omitempty" structs:"schemas,omitempty"`
			}
			Links struct {
				Self                       *OPGenericLink  `json:"self,omitempty" structs:"self,omitempty"`
				JumpTo                     *OPGenericLink  `json:"jumpToWorkPackage,omitempty" structs:"jumpToWorkPackage,omitempty"`
				ChangeSize                 *OPGenericLink  `json:"changePageSize,omitempty" structs:"changePageSize,omitempty"`
				NextByOffset               *OPGenericLink  `json:"nextByOffset,omitempty" structs:"nextByOffset,omitempty"`
				EditWorkPackage            *OPGenericLink  `json:"editWorkPackage,omitempty" structs:"editWorkPackage,omitempty"`
				CreateWorkPackage          *OPGenericLink  `json:"createWorkPackage,omitempty" structs:"createWorkPackage,omitempty"`
				CreateWorkPackageImmediate *OPGenericLink  `json:"createWorkPackageImmediate,omitempty" structs:"createWorkPackageImmediate,omitempty"`
				Schemas                    *OPGenericLink  `json:"schemas,omitempty" structs:"schemas,omitempty"`
				CustomFields               *OPGenericLink  `json:"customFields,omitempty" structs:"customFields,omitempty"`
				Representations            []OPGenericLink `json:"representations,omitempty" structs:"representations,omitempty"`
			} `json:"_links,omitempty" structs:"_links,omitempty"`
		} `json:"results,omitempty" structs:"results,omitempty"`
	} `json:"_embedded,omitempty" structs:"_embedded,omitempty"`
	Links struct {
		Self                      *OPGenericLink  `json:"self,omitempty" structs:"self,omitempty"`
		Project                   *OPGenericLink  `json:"project,omitempty" structs:"project,omitempty"`
		Results                   *OPGenericLink  `json:"results,omitempty" structs:"results,omitempty"`
		Star                      *OPGenericLink  `json:"star,omitempty" structs:"star,omitempty"`
		Schema                    *OPGenericLink  `json:"schema,omitempty" structs:"schema,omitempty"`
		Update                    *OPGenericLink  `json:"update,omitempty" structs:"update,omitempty"`
		UpdateImmediately         *OPGenericLink  `json:"updateImmediately,omitempty" structs:"updateImmediately,omitempty"`
		UpdateOrderedWorkPackages *OPGenericLink  `json:"updateOrderedWorkPackages,omitempty" structs:"updateOrderedWorkPackages,omitempty"`
		Delete                    *OPGenericLink  `json:"delete,omitempty" structs:"delete,omitempty"`
		User                      *OPGenericLink  `json:"user,omitempty" structs:"user,omitempty"`
		SortBy                    []OPGenericLink `json:"sortBy,omitempty" structs:"sortBy,omitempty"`
		GroupBy                   *OPGenericLink  `json:"groupBy,omitempty" structs:"groupBy,omitempty"`
		Columns                   []OPGenericLink `json:"columns,omitempty" structs:"columns,omitempty"`
		HighlightedAttributes     []OPGenericLink `json:"highlightedAttributes,omitempty" structs:"highlightedAttributes,omitempty"`
	} `json:"_links,omitempty" structs:"_links,omitempty"`
	TimelineLabels interface{}     `json:"timelineLabels,omitempty"`
	Columns        []OPGenericLink `json:"columns,omitempty" structs:"columns,omitempty"`
	GroupBy        []OPGenericLink `json:"groupBy,omitempty" structs:"groupBy,omitempty"`
	SortBy         []OPGenericLink `json:"sortBy,omitempty" structs:"sortBy,omitempty"`
}

// QueryFilter filters within a query
type QueryFilter struct {
	Type  string `json:"_type,omitempty" structs:"_type,omitempty"`
	Name  string `json:"name,omitempty" structs:"name,omitempty"`
	Links struct {
		Schema   *OPGenericLink `json:"schema,omitempty" structs:"schema,omitempty"`
		Filter   *OPGenericLink `json:"filter,omitempty" structs:"filter,omitempty"`
		Operator *OPGenericLink `json:"operator,omitempty" structs:"operator,omitempty"`
		// TODO: Complete fields
		Values []interface{} `json:"values,omitempty"`
	} `json:"_links,omitempty" structs:"_links,omitempty"`
}

type QueryEmbeddedColumn struct {
	Type  string `json:"_type,omitempty"`
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Links struct {
		Self *OPGenericLink `json:"self,omitempty"`
	} `json:"_links,omitempty"`
}

type QuerySortBy struct {
	QueryEmbeddedColumn
	Links struct {
		Self      *OPGenericLink `json:"self,omitempty"`
		Column    *OPGenericLink `json:"column,omitempty"`
		Direction *OPGenericLink `json:"direction,omitempty"`
	} `json:"_links,omitempty"`
}

// GetWithContext gets query info from OpenProject using its query ID
func (s *QueryService) GetWithContext(ctx context.Context, queryID string) (*QueryResult, *Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/queries/%s", queryID)
	Obj, Resp, err := GetWithContext(ctx, s, apiEndpoint)
	if err != nil {
		return nil, Resp, err
	}
	return Obj.(*QueryResult), Resp, err
}

// Get wraps GetWithContext using the background context.
func (s *QueryService) Get(queryID string) (*QueryResult, *Response, error) {
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
func (s *QueryService) CreateWithContext(ctx context.Context, queryObj *QueryResult) (*QueryResult, *Response, error) {
	apiEndpoint := "api/v3/queries"
	wpResponse, resp, err := CreateWithContext(ctx, queryObj, s, apiEndpoint)
	if err != nil {
		return nil, resp, err
	}
	return wpResponse.(*QueryResult), resp, err
}

// Create wraps CreateWithContext using the background context.
func (s *QueryService) Create(queryObj *QueryResult) (*QueryResult, *Response, error) {
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
