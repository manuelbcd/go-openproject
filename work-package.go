package openproject

import (
	"context"
	"fmt"

	"github.com/google/go-querystring/query"
)

/**
WorkPackageService handles workpackages for the OpenProject instance / API.
*/
type WorkPackageService struct {
	client *Client
}

/**
	searchResult is only a small wrapper around the Search
 */
type searchResult struct {
	WorkPackages  	[]WorkPackage 	`json:"workpackages" structs:"workpackages"`
	StartAt    		int     		`json:"startAt" structs:"startAt"`
	MaxResults 		int     		`json:"maxResults" structs:"maxResults"`
	Total      		int     		`json:"total" structs:"total"`
}

/**
Issue represents an OpenProject WorkPackage.
*/
type WorkPackage struct {
	Type    string `json:"_type,omitempty" structs:"_type,omitempty"`
	ID      string `json:"id,omitempty" structs:"id,omitempty"`
	Subject string `json:"subject,omitempty" structs:"subject,omitempty"`
}

/**
	GetWithContext returns a full representation of the issue for the given OpenProject key.
 	The given options will be appended to the query string
*/
func (s *WorkPackageService) GetWithContext(ctx context.Context, workpackageID string, options *GetQueryOptions) (*WorkPackage, *Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/work_packages/%s", workpackageID)
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

	issue := new(WorkPackage)
	resp, err := s.client.Do(req, issue)
	if err != nil {
		jerr := NewOpenProjectError(resp, err)
		return nil, resp, jerr
	}

	return issue, resp, nil
}

/**
	Get wraps GetWithContext using the background context.
*/
func (s *WorkPackageService) Get(issueID string, options *GetQueryOptions) (*WorkPackage, *Response, error) {
	return s.GetWithContext(context.Background(), issueID, options)
}
