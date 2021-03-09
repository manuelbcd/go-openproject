package openproject

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/trivago/tgo/tcontainer"

	"io/ioutil"
	"net/url"
	"time"
)

/**
WorkPackageService handles workpackages for the OpenProject instance / API.
*/
type WorkPackageService struct {
	client *Client
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
	Custom      tcontainer.MarshalMap

	Links *WPLinks `json:"_links,omitempty" _links:"id,omitempty"`
}

/**
WorkPackageDescription type contains description and format
*/
type WPDescription OPGenericDescription

/**
WorkPackage Links
*/
type WPLinks struct {
	Self     WPLinksField `json:"self,omitempty" structs:"self,omitempty"`
	Type     WPLinksField `json:"type,omitempty" structs:"type,omitempty"`
	Priority WPLinksField `json:"priority,omitempty" structs:"priority,omitempty"`
	Status   WPLinksField `json:"status,omitempty" structs:"status,omitempty"`
}

type WPLinksField struct {
	Href  string
	Title string
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
	DIFFERENT      SearchOperator = 1 // <>
	GREATERTHAN    SearchOperator = 2 // >
	LOWERTHAN      SearchOperator = 3 // <
	SEARCHSTRING   SearchOperator = 4 // '**'
	LIKE           SearchOperator = 5 // ~
	GREATEROREQUAL SearchOperator = 6 // >=
	LOWEROREQUAL   SearchOperator = 7 // <=
)

/**
Constants to represent OpenProject standard GET parameters
*/
const PARAM_FILTERS = "filters"

/**
FilterOptions allows you to specify search parameters for the get-workpackage action
When used they will be converted to GET parameters within the URL
Up to now OpenProject only allows "AND" combinations. "OR" combinations feature is under development,
tracked by this ticket https://community.openproject.org/projects/openproject/work_packages/26837/activity

More information about filters https://docs.openproject.org/api/filters/
*/
type FilterOptions struct {
	Fields []OptionsFields
}

type OptionsFields struct {
	Field    string
	Operator SearchOperator
	Value    string
}

/**
searchResult is only a small wrapper around the Search
*/
type searchResultWP struct {
	Embedded searchEmbeddedWP `json:"_embedded" structs:"_embedded"`
	Total    int              `json:"total" structs:"total"`
	Count    int              `json:"count" structs:"count"`
	PageSize int              `json:"pageSize" structs:"pageSize"`
	Offset   int              `json:"offset" structs:"offset"`
}

type searchEmbeddedWP struct {
	Elements []WorkPackage `json:"elements" structs:"elements"`
}

/**
	GetWithContext returns a full representation of the issue for the given OpenProject key.
 	The given options will be appended to the query string
*/
func (s *WorkPackageService) GetWithContext(ctx context.Context, workpackageID string) (*WorkPackage, *Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/work_packages/%s", workpackageID)

	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	workPackage := new(WorkPackage)
	resp, err := s.client.Do(req, workPackage)
	if err != nil {
		jerr := NewOpenProjectError(resp, err)
		return nil, resp, jerr
	}

	return workPackage, resp, nil
}

/**
Get wraps GetWithContext using the background context.
*/
func (s *WorkPackageService) Get(issueID string) (*WorkPackage, *Response, error) {
	return s.GetWithContext(context.Background(), issueID)
}

/**
	prepareFilters convert FilterOptions to single URL-Encoded string to be inserted into GET request
    as parameter.
*/
func (fops *FilterOptions) prepareFilters() url.Values {
	values := make(url.Values)

	filterTemplate := "["
	for _, field := range fops.Fields {
		s := fmt.Sprintf(
			"{\"%[1]v\":{\"operator\":\"%[2]v\",\"values\":[\"%[3]v\"]}}",
			field.Field, interpretOperator(field.Operator), field.Value)

		filterTemplate += s
	}
	filterTemplate += "]"

	values.Add(PARAM_FILTERS, filterTemplate)

	return values
}

/**
Interpret Operator collection and return its string
*/
func interpretOperator(operator SearchOperator) string {
	result := "="

	switch operator {
	case GREATERTHAN:
		result = ">"
	case LOWERTHAN:
		result = "<"
	case DIFFERENT:
		result = "<>"
	}
	//TODO: Complete list of operators

	return result
}

/**
	CreateWithContext creates a work-package or a sub-task from a JSON representation.
**/
func (s *WorkPackageService) CreateWithContext(ctx context.Context, projectName string, workPackage *WorkPackage) (*WorkPackage, *Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/projects/%s/work_packages", projectName)
	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, workPackage)
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

/**
Create wraps CreateWithContext using the background context.
*/
func (s *WorkPackageService) Create(workPackage *WorkPackage, projectName string) (*WorkPackage, *Response, error) {
	return s.CreateWithContext(context.Background(), projectName, workPackage)
}

/**
GetListWithContext will retrieve a list of work-packages using filters
*/
func (s *WorkPackageService) GetListWithContext(ctx context.Context, options *FilterOptions) ([]WorkPackage, *Response, error) {
	u := url.URL{
		Path: "api/v3/work_packages",
	}

	req, err := s.client.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return []WorkPackage{}, nil, err
	}

	if options != nil {
		values := options.prepareFilters()
		req.URL.RawQuery = values.Encode()
	}

	v := new(searchResultWP)
	resp, err := s.client.Do(req, v)
	if err != nil {
		err = NewOpenProjectError(resp, err)
	}
	return v.Embedded.Elements, resp, err
}

/**
GetList wraps GetListWithContext using the background context.
*/
func (s *WorkPackageService) GetList(options *FilterOptions) ([]WorkPackage, *Response, error) {
	return s.GetListWithContext(context.Background(), options)
}
