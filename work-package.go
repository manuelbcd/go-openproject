package openproject

import (
	"context"
	"encoding/json"
	"fmt"
)

/**
	WorkPackageService handles workpackages for the OpenProject instance / API.
 */
type WorkPackageService struct {
	client *Client
}

/**
	Issue represents an OpenProject WorkPackage.
 */
type WorkPackage struct {
	ID             string               `json:"id,omitempty" structs:"id,omitempty"`
	Self           string               `json:"self,omitempty" structs:"self,omitempty"`
	Key            string               `json:"key,omitempty" structs:"key,omitempty"`
	Fields         *IssueFields         `json:"fields,omitempty" structs:"fields,omitempty"`
	RenderedFields *IssueRenderedFields `json:"renderedFields,omitempty" structs:"renderedFields,omitempty"`
	Changelog      *Changelog           `json:"changelog,omitempty" structs:"changelog,omitempty"`
	Transitions    []Transition         `json:"transitions,omitempty" structs:"transitions,omitempty"`
	Names          map[string]string    `json:"names,omitempty" structs:"names,omitempty"`
}

/**
	GetWithContext returns a full representation of the issue for the given OpenProject key.
 	The given options will be appended to the query string
 */
func (s *WorkPackageService) GetWithContext(ctx context.Context, issueID string, options *GetQueryOptions) (*Issue, *Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/issue/%s", issueID)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	if options != nil {
		q, err := query.Values(options)
		if err != nil {
			return nil, nil, err
		}
		req.URL.RawQuery = q.Encode()
	}

	issue := new(Issue)
	resp, err := s.client.Do(req, issue)
	if err != nil {
		jerr := NewJiraError(resp, err)
		return nil, resp, jerr
	}

	return issue, resp, nil
}

// Get wraps GetWithContext using the background context.
func (s *WorkPackageService) Get(issueID string, options *GetQueryOptions) (*Issue, *Response, error) {
	return s.GetWithContext(context.Background(), issueID, options)
}