package openproject

import (
	"context"
	"fmt"
)

/**
StatusService handles statuses from the OpenProject instance / API.
*/
type StatusService struct {
	client *Client
}

// StatusList represent a list of Projects
type searchResultStatus struct {
	Embedded statusElements `json:"_embedded,omitempty" structs:"_embedded,omitempty"`
	Total    int            `json:"total" structs:"total"`
	Count    int            `json:"count" structs:"count"`
	PageSize int            `json:"pageSize" structs:"pageSize"`
	Offset   int            `json:"offset" structs:"offset"`
}

type statusElements struct {
	Elements []Status `json:"elements,omitempty" structs:"elements,omitempty"`
}

/**
Status is the object representing OpenProject statuses.
TODO: Complete fields (add defaultDoneRatio, _links, ...)
*/
type Status struct {
	Type       string `json:"_type,omitempty" structs:"_type,omitempty"`
	Id         int    `json:"id,omitempty" structs:"id,omitempty"`
	Name       string `json:"name,omitempty" structs:"name,omitempty"`
	IsClosed   bool   `json:"isClosed,omitempty" structs:"isClosed,omitempty"`
	Color      string `json:"color,omitempty" structs:"color,omitempty"`
	IsDefault  bool   `json:"isDefault,omitempty" structs:"isDefault,omitempty"`
	IsReadOnly bool   `json:"isReadOnly,omitempty" structs:"isReadOnly,omitempty"`
	Position   int    `json:"position,omitempty" structs:"position,omitempty"`
}

/**
GetWithContext gets statuses info from OpenProject using its status Id
// TODO: Implement GetList and adapt tests
*/
func (s *StatusService) GetWithContext(ctx context.Context, statusID string) (*Status, *Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/statuses/%s", statusID)
	Obj, Resp, err := GetWithContext(s, ctx, apiEndpoint)
	return Obj.(*Status), Resp, err
}

/**
Get wraps GetWithContext using the background context.
*/
func (s *StatusService) Get(statusID string) (*Status, *Response, error) {
	return s.GetWithContext(context.Background(), statusID)
}

/**
GetList wraps GetListWithContext using the background context.
*/
func (s *StatusService) GetList() (*searchResultStatus, *Response, error) {
	return s.GetListWithContext(context.Background())
}

/**
Retrieve status list with context
TODO: Implement search parameters-options
*/
func (s *StatusService) GetListWithContext(ctx context.Context) (*searchResultStatus, *Response, error) {
	apiEndpoint := "api/v3/statuses"
	Obj, Resp, err := GetListWithContext(s, ctx, apiEndpoint, nil)
	return Obj.(*searchResultStatus), Resp, err
}
