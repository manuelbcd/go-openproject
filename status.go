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
GetWithContext gets statuses info from OpenProject using its Account Id
// TODO: Implement GetList and adapt tests
*/
func (s *StatusService) GetWithContext(ctx context.Context, accountId string) (*Status, *Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/statuses/%s", accountId)
	Obj, Resp, err := GetWithContext(s, ctx, apiEndpoint)
	return Obj.(*Status), Resp, err
}

/**
Get wraps GetWithContext using the background context.
*/
func (s *StatusService) Get(accountId string) (*Status, *Response, error) {
	return s.GetWithContext(context.Background(), accountId)
}
