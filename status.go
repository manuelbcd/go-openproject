package openproject

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"math"
	"sync"
)

// StatusService handles statuses from the OpenProject instance / API.
type StatusService struct {
	client *Client
}

// SearchResultStatus represent a list of Statuses
type SearchResultStatus struct {
	Embedded statusElements `json:"_embedded,omitempty" structs:"_embedded,omitempty"`
	PaginationParam
	lock *sync.RWMutex
}

func (s *SearchResultStatus) TotalPage() int {
	return int(math.Ceil(float64(s.Total) / float64(s.PageSize)))
}

func (s *SearchResultStatus) ConcatEmbed(status interface{}) {
	s.lock.Lock()
	s.Embedded.Elements = append(s.Embedded.Elements, status.(*SearchResultStatus).Embedded.Elements...)
}

// statusElements array wraps elemets within searchResultStatus
type statusElements struct {
	Elements []Status `json:"elements,omitempty" structs:"elements,omitempty"`
}

// Status is the object representing OpenProject statuses.
// TODO: Complete fields (add defaultDoneRatio, _links, ...)
type Status struct {
	Type       string `json:"_type,omitempty" structs:"_type,omitempty"`
	ID         int    `json:"id,omitempty" structs:"id,omitempty"`
	Name       string `json:"name,omitempty" structs:"name,omitempty"`
	IsClosed   bool   `json:"isClosed,omitempty" structs:"isClosed,omitempty"`
	Color      string `json:"color,omitempty" structs:"color,omitempty"`
	IsDefault  bool   `json:"isDefault,omitempty" structs:"isDefault,omitempty"`
	IsReadOnly bool   `json:"isReadOnly,omitempty" structs:"isReadOnly,omitempty"`
	Position   int    `json:"position,omitempty" structs:"position,omitempty"`
}

// GetWithContext gets statuses info from OpenProject using its status ID
// TODO: Implement GetList and adapt tests
func (s *StatusService) GetWithContext(ctx context.Context, statusID string) (*Status, *Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/statuses/%s", statusID)
	Obj, Resp, err := GetWithContext(ctx, s, apiEndpoint)
	if err != nil {
		return nil, Resp, err
	}
	return Obj.(*Status), Resp, err
}

// Get wraps GetWithContext using the background context.
func (s *StatusService) Get(statusID string) (*Status, *Response, error) {
	return s.GetWithContext(context.Background(), statusID)
}

// GetList wraps GetListWithContext using the background context.
func (s *StatusService) GetList(offset int, pageSize int) (*SearchResultStatus, *Response, error) {
	return s.GetListWithContext(context.Background(), offset, pageSize)
}

// GetListWithContext Retrieve status list with context
func (s *StatusService) GetListWithContext(ctx context.Context, offset int, pageSize int) (*SearchResultStatus, *Response, error) {
	apiEndpoint := "api/v3/statuses"
	Obj, Resp, err := GetListWithContext(ctx, s, apiEndpoint, nil, offset, pageSize)
	if err != nil {
		return nil, Resp, err
	}
	status, success := Obj.(*SearchResultStatus)
	if !success {
		return nil, nil, errors.New("convert to SearchResultStatus failed")
	}
	status.lock = &sync.RWMutex{}
	return Obj.(*SearchResultStatus), Resp, err
}
