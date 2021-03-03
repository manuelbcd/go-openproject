package openproject

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/trivago/tgo/tcontainer"
	"io/ioutil"
	"time"
)

/**
WorkPackageService handles workpackages for the OpenProject instance / API.
*/
type WorkPackageService struct {
	client *Client
}

// Time represents the Time definition of OpenProject as a time.Time of go
type Time time.Time

// Date represents the Date definition of OpenProject as a time.Time of go
type Date time.Time

func (t Time) Equal(u Time) bool {
	return time.Time(t).Equal(time.Time(u))
}

/**
UnmarshalJSON will transform the OpenProject time into a time.Time
during the transformation of the OpenProject JSON response
*/
func (t *Time) UnmarshalJSON(b []byte) error {
	// Ignore null, like in the main JSON package.
	if string(b) == "null" {
		return nil
	}
	ti, err := time.Parse("\"2006-01-02T15:04:05Z\"", string(b))
	if err != nil {
		return err
	}
	*t = Time(ti)
	return nil
}

/**
MarshalJSON will transform the time.Time into a OpenProject time
during the creation of a OpenProject request
*/
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format("\"2006-01-02T15:04:05Z\"")), nil
}

/**
UnmarshalJSON will transform the OpenProject date into a time.Time
during the transformation of the OpenProject JSON response
*/
func (t *Date) UnmarshalJSON(b []byte) error {
	// Ignore null, like in the main JSON package.
	if string(b) == "null" {
		return nil
	}
	ti, err := time.Parse("\"2006-01-02\"", string(b))
	if err != nil {
		return err
	}
	*t = Date(ti)
	return nil
}

/**
MarshalJSON will transform the Date object into a short
date string as OpenProject expects during the creation of a
OpenProject request
*/
func (t Date) MarshalJSON() ([]byte, error) {
	time := time.Time(t)
	return []byte(time.Format("\"2006-01-02\"")), nil
}

/**
Issue represents an OpenProject WorkPackage.

Please note: Time and Date fields are pointers in order to avoid rendering them when not initialized
*/
type WorkPackage struct {
	Subject     string         `json:"subject,omitempty" structs:"subject,omitempty"`
	Description *WPDescription `json:"description,omitempty" structs:"description,omitempty"`
	Type        string         `json:"_type,omitempty" structs:"_type,omitempty"`
	Id          int            `json:"id,omitempty" structs:"id,omitempty"`
	CreatedAt   *Time          `json:"createdAt,omitempty" structs:"createdAt,omitempty"`
	UpdatedAt   *Time          `json:"updatedAt,omitempty" structs:"updatedAt,omitempty"`
	StartDate   *Date          `json:"startDate,omitempty" structs:"startDate,omitempty"`
	DueDate     *Date          `json:"dueDate,omitempty" structs:"dueDate,omitempty"`
	LockVersion int            `json:"lockVersion,omitempty" structs:"lockVersion,omitempty"`
	Position    int            `json:"position,omitempty" structs:"position,omitempty"`

	Custom tcontainer.MarshalMap
}

/**
WorkPackage type
*/
type WPDescription struct {
	Format string `json:"format,omitempty" structs:"format,omitempty"`
	Raw    string `json:"raw,omitempty" structs:"raw,omitempty"`
	Html   string `json:"html,omitempty" structs:"html,omitempty"`
}

/**
WorkPackage form
OpenProject API v3 provides a WorkPackage form to get a template of work-packages dynamically
A "Form" endpoint is available for that purpose.
*/
type WPForm struct {
	Type     string         `json:"_type,omitempty" structs:"_type,omitempty"`
	Embedded WPFormEmbedded `json:"_embedded,omitempty" structs:"_embedded,omitempty"`
	Links    WPFormLinks    `json:"_links,omitempty" structs:"_links,omitempty"`
}

/**
WPFormEmbedded represents the 'embedded' struct nested in 'form'
*/
type WPFormEmbedded struct {
	Payload WPPayload `json:"payload,omitempty" structs:"payload,omitempty"`
}

/**
WPPayload represents the 'payload' struct nested in 'form.embedded'
*/
type WPPayload struct {
	Subject string `json:"subject,omitempty" structs:"subject,omitempty"`

	StartDate string `json:"startDate,omitempty" structs:"startDate,omitempty"`
}

type WPFormLinks struct {
}

/**
Search operators
Doc. https://docs.openproject.org/api/filters/#header-available-filters-1
*/
type SearchOperator int32

const (
	EQUAL          SearchOperator = 0 // =
	NOTEQUAL       SearchOperator = 1 // <>
	GREATERTHAN    SearchOperator = 2 // >
	LOWERTHAN      SearchOperator = 3 // <
	SEARCHSTRING   SearchOperator = 4 // **
	LIKE           SearchOperator = 5 // ~
	GREATEROREQUAL SearchOperator = 6 // >=
	LOWEROREQUAL   SearchOperator = 7 // <=
)

/**
SearchOptions allows you to specify search parameters.
When used they will be converted to GET parameters within the URL
*/
type SearchOptions struct {
	Fields []struct {
		Field    string
		Operator SearchOperator
		Value    string
	}
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
