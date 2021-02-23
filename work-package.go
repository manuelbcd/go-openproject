package openproject

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

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
	WorkPackages []WorkPackage `json:"workpackages" structs:"workpackages"`
	StartAt      int           `json:"startAt" structs:"startAt"`
	MaxResults   int           `json:"maxResults" structs:"maxResults"`
	Total        int           `json:"total" structs:"total"`
}

/**
	Issue represents an OpenProject WorkPackage.
*/
type WorkPackage struct {
	Subject 		string 			`json:"subject,omitempty" structs:"subject,omitempty"`
	Description		*WPDescription	`json:"description,omitempty" structs:"description,omitempty"`
}

/**
	WorkPackage type
*/
type WPDescription struct {
	Format      string `json:"format,omitempty" structs:"format,omitempty"`
	Raw         string `json:"raw,omitempty" structs:"raw,omitempty"`
	Html 		string `json:"html,omitempty" structs:"html,omitempty"`
}

/**
	WorkPackage form
	OpenProject API v3 provides a WorkPackage form to get a template of work-packages dynamically
	A "Form" endpoint is available for that purpose.
 */
type WPForm struct {
	Type		string 			`json:"_type,omitempty" structs:"_type,omitempty"`
	Embedded	WPFormEmbedded  `json:"_embedded,omitempty" structs:"_embedded,omitempty"`
	Links		WPFormLinks		`json:"_links,omitempty" structs:"_links,omitempty"`
}

/**
	WPFormEmbedded represents the 'embedded' struct nested in 'form'
 */
type WPFormEmbedded struct {
	Payload		WPPayload		`json:"payload,omitempty" structs:"payload,omitempty"`
}

/**
	WPPayload represents the 'payload' struct nested in 'form.embedded'
 */
type WPPayload struct {
	Subject		string			`json:"subject,omitempty" structs:"subject,omitempty"`

	StartDate	string			`json:"startDate,omitempty" structs:"startDate,omitempty"`
}

type WPFormLinks struct {

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

/**
	CreateWithContext creates a work-package or a sub-task from a JSON representation.
**/
func (s *WorkPackageService) CreateWithContext(ctx context.Context, projectName string, issue *WorkPackage) (*WorkPackage, *Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/projects/%s/work_packages", projectName)
	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, issue)
	if err != nil {
		return nil, nil, err
	}
	resp, err := s.client.Do(req, nil)
	if err != nil {
		// incase of error return the resp for further inspection
		return nil, resp, err
	}

	wpResponse := new(WorkPackage)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("could not read the returned data")
	}
	err = json.Unmarshal(data, wpResponse)
	if err != nil {
		return nil, resp, fmt.Errorf("could not unmarshall the data into struct")
	}
	return wpResponse, resp, nil
}

// Create wraps CreateWithContext using the background context.
func (s *WorkPackageService) Create(issue *WorkPackage, projectName string) (*WorkPackage, *Response, error) {
	return s.CreateWithContext(context.Background(), projectName, issue)
}