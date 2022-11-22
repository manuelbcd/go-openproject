package openproject

import (
	"context"
	"fmt"
	"github.com/trivago/tgo/tcontainer"
	"math"
	"strings"

	"net/url"
	"time"
)

// WorkPackageService handles workpackages for the OpenProject instance / API.
type WorkPackageService struct {
	client *Client
}

// UnmarshalJSON will transform the OpenProject time into a time.Time
// during the transformation of the OpenProject JSON response
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

// MarshalJSON will transform the time.Time into a OpenProject time
// during the creation of a OpenProject request
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format("\"2006-01-02T15:04:05Z\"")), nil
}

// UnmarshalJSON will transform the OpenProject date into a time.Time
// during the transformation of the OpenProject JSON response
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

// MarshalJSON will transform the Date object into a short
// date string as OpenProject expects during the creation of a
// OpenProject request
func (t Date) MarshalJSON() ([]byte, error) {
	time := time.Time(t)
	return []byte(time.Format("\"2006-01-02\"")), nil
}

// WorkPackage represents an OpenProject ticket or issue
type WorkPackage struct {
	DerivedStartDate     interface{}           `json:"derivedStartDate,omitempty"`
	DerivedDueDate       interface{}           `json:"derivedDueDate,omitempty"`
	SpentTime            string                `json:"spentTime,omitempty"`
	LaborCosts           string                `json:"laborCosts,omitempty"`
	MaterialCosts        string                `json:"materialCosts,omitempty"`
	OverallCosts         string                `json:"overallCosts,omitempty"`
	Subject              string                `json:"subject,omitempty" structs:"subject,omitempty"`
	Description          *WPDescription        `json:"description,omitempty" structs:"description,omitempty"`
	Type                 string                `json:"_type,omitempty" structs:"_type,omitempty"`
	ID                   int                   `json:"id,omitempty" structs:"id,omitempty"`
	CreatedAt            *Time                 `json:"createdAt,omitempty" structs:"createdAt,omitempty"`
	UpdatedAt            *Time                 `json:"updatedAt,omitempty" structs:"updatedAt,omitempty"`
	StartDate            *Date                 `json:"startDate,omitempty" structs:"startDate,omitempty"`
	ScheduleManually     bool                  `json:"scheduleManually,omitempty"`
	DueDate              *Date                 `json:"dueDate,omitempty" structs:"dueDate,omitempty"`
	LockVersion          int                   `json:"lockVersion,omitempty" structs:"lockVersion,omitempty"`
	Position             int                   `json:"position,omitempty" structs:"position,omitempty"`
	Custom               tcontainer.MarshalMap `json:"custom,omitempty" structs:"custom,omitempty"`
	EstimatedTime        string                `json:"estimatedTime,omitempty"`
	DerivedEstimatedTime interface{}           `json:"derivedEstimatedTime,omitempty"`
	PercentageDone       int                   `json:"percentageDone,omitempty"`
	RemainingTime        string                `json:"remainingTime,omitempty"`

	Links *WPLinks `json:"_links,omitempty" _links:"id,omitempty"`
}

// WPDescription type contains description and format
type WPDescription OPGenericDescription

// WPLinks are WorkPackage Links
type WPLinks struct {
	Self                        *OPGenericLink  `json:"self,omitempty"`
	Update                      *OPGenericLink  `json:"update,omitempty"`
	Schema                      *OPGenericLink  `json:"schema,omitempty"`
	UpdateImmediately           *OPGenericLink  `json:"updateImmediately,omitempty"`
	Delete                      *OPGenericLink  `json:"delete,omitempty"`
	LogTime                     *OPGenericLink  `json:"logTime,omitempty"`
	Move                        *OPGenericLink  `json:"move,omitempty"`
	Copy                        *OPGenericLink  `json:"copy,omitempty"`
	Pdf                         *OPGenericLink  `json:"pdf,omitempty"`
	Atom                        *OPGenericLink  `json:"atom,omitempty"`
	AvailableRelationCandidates *OPGenericLink  `json:"availableRelationCandidates,omitempty"`
	CustomFields                *OPGenericLink  `json:"customFields,omitempty"`
	ConfigureForm               *OPGenericLink  `json:"configureForm,omitempty"`
	Activities                  *OPGenericLink  `json:"activities,omitempty"`
	AvailableWatchers           *OPGenericLink  `json:"availableWatchers,omitempty"`
	Relations                   *OPGenericLink  `json:"relations,omitempty"`
	Revisions                   *OPGenericLink  `json:"revisions,omitempty"`
	Watchers                    *OPGenericLink  `json:"watchers,omitempty"`
	AddRelation                 *OPGenericLink  `json:"addRelation,omitempty"`
	AddChild                    *OPGenericLink  `json:"addChild,omitempty"`
	ChangeParent                *OPGenericLink  `json:"changeParent,omitempty"`
	AddComment                  *OPGenericLink  `json:"addComment,omitempty"`
	PreviewMarkup               *OPGenericLink  `json:"previewMarkup,omitempty"`
	TimeEntries                 *OPGenericLink  `json:"timeEntries,omitempty"`
	Category                    *OPGenericLink  `json:"category,omitempty"`
	Type                        *OPGenericLink  `json:"type,omitempty"`
	Priority                    *OPGenericLink  `json:"priority,omitempty"`
	Project                     *OPGenericLink  `json:"project,omitempty"`
	Status                      *OPGenericLink  `json:"status,omitempty"`
	Author                      *OPGenericLink  `json:"author,omitempty"`
	Responsible                 *OPGenericLink  `json:"responsible,omitempty"`
	Assignee                    *OPGenericLink  `json:"assignee,omitempty"`
	Version                     *OPGenericLink  `json:"version,omitempty"`
	LogCosts                    *OPGenericLink  `json:"logCosts,omitempty"`
	ShowCosts                   *OPGenericLink  `json:"showCosts,omitempty"`
	CostsByType                 *OPGenericLink  `json:"costsByType,omitempty"`
	GithubPullRequests          *OPGenericLink  `json:"github_pull_requests,omitempty"`
	Parent                      *OPGenericLink  `json:"parent,omitempty"`
	Ancestors                   []OPGenericLink `json:"ancestors,omitempty"`
	CustomActions               []OPGenericLink `json:"customActions,omitempty"`
	RemoveWatcher               struct {
		Href      string `json:"href,omitempty"`
		Method    string `json:"method,omitempty"`
		Templated bool   `json:"templated,omitempty"`
	} `json:"removeWatcher,omitempty"`
	AddWatcher struct {
		Href    string `json:"href,omitempty"`
		Method  string `json:"method,omitempty"`
		Payload struct {
			User struct {
				Href string `json:"href,omitempty"`
			} `json:"user,omitempty"`
		} `json:"payload,omitempty"`
		Templated bool `json:"templated,omitempty"`
	} `json:"addWatcher,omitempty"`
	CustomField1 []interface{} `json:"customField1,omitempty"`
	Watch        struct {
		Href    string `json:"href,omitempty"`
		Method  string `json:"method,omitempty"`
		Payload struct {
			User struct {
				Href string `json:"href,omitempty"`
			} `json:"user,omitempty"`
		} `json:"payload,omitempty"`
	} `json:"watch,omitempty"`
}

// WPForm represents WorkPackage form
// OpenProject API v3 provides a WorkPackage form to get a template of work-packages dynamically
// A "Form" endpoint is available for that purpose.
type WPForm struct {
	Type     string         `json:"_type,omitempty" structs:"_type,omitempty"`
	Embedded WPFormEmbedded `json:"_embedded,omitempty" structs:"_embedded,omitempty"`
	Links    WPFormLinks    `json:"_links,omitempty" structs:"_links,omitempty"`
}

// WPFormEmbedded represents the 'embedded' struct nested in 'form'
type WPFormEmbedded struct {
	Payload WPPayload `json:"payload,omitempty" structs:"payload,omitempty"`
}

// WPPayload represents the 'payload' struct nested in 'form.embedded'
type WPPayload struct {
	Subject string `json:"subject,omitempty" structs:"subject,omitempty"`

	StartDate string `json:"startDate,omitempty" structs:"startDate,omitempty"`
}

// WPFormLinks represents WorkPackage Form Links
type WPFormLinks struct {
}

// Constants to represent OpenProject standard GET parameters
const paramFilters = "filters"

// FilterOptions allows you to specify search parameters for the get-workpackage action
// When used they will be converted to GET parameters within the URL
// Up to now OpenProject only allows "AND" combinations. "OR" combinations feature is under development,
// tracked by this ticket https://community.openproject.org/projects/openproject/work_packages/26837/activity
// More information about filters https://docs.openproject.org/api/filters/
type FilterOptions struct {
	Fields []OptionsFields
}

// OptionsFields array wraps field, Operator, Value within FilterOptions
type OptionsFields struct {
	Field    string
	Operator SearchOperator
	Value    string
}

// SearchResultWP is only a small wrapper around the Search
type SearchResultWP struct {
	Embedded SearchEmbeddedWP `json:"_embedded" structs:"_embedded"`
	PaginationParam
}

func (s *SearchResultWP) TotalPage() int {
	return int(math.Ceil(float64(s.Total) / float64(s.PageSize)))
}

func (s *SearchResultWP) ConcatEmbed(wp interface{}) {
	s.Embedded.Elements = append(s.Embedded.Elements, wp.(*SearchResultWP).Embedded.Elements...)
}

// SearchEmbeddedWP represent elements within WorkPackage list
type SearchEmbeddedWP struct {
	Elements []WorkPackage `json:"elements" structs:"elements"`
}

// GetWithContext returns a full representation of the issue for the given OpenProject key.
// The given options will be appended to the query string
func (s *WorkPackageService) GetWithContext(ctx context.Context, workpackageID string) (*WorkPackage, *Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/work_packages/%s", workpackageID)
	Obj, Resp, err := GetWithContext(ctx, s, apiEndpoint)
	if err != nil {
		return nil, Resp, err
	}
	return Obj.(*WorkPackage), Resp, err
}

// Get wraps GetWithContext using the background context.
func (s *WorkPackageService) Get(workpackageID string) (*WorkPackage, *Response, error) {
	return s.GetWithContext(context.Background(), workpackageID)
}

//	prepareFilters convert FilterOptions to single URL-Encoded string to be inserted into GET request
//
// as parameter.
func (fops *FilterOptions) prepareFilters(oldValues url.Values) url.Values {
	values := oldValues
	if oldValues == nil {
		values = make(url.Values)
	}

	filterTemplate := "["
	for idx, field := range fops.Fields {
		value := ""
		if strings.Contains(field.Value, ",") {
			parts := strings.Split(field.Value, ",")
			for i, part := range parts {
				value += fmt.Sprintf(`"%s"`, part)
				if i < len(parts)-1 {
					value += ","
				}
			}
		} else {
			value = `"` + field.Value + `"`
		}

		s := fmt.Sprintf(
			"{\"%[1]v\":{\"operator\":\"%[2]v\",\"values\":[%[3]v]}}",
			field.Field, field.Operator, value)

		filterTemplate += s

		if idx != len(fops.Fields)-1 {
			filterTemplate += ","
		}
	}
	filterTemplate += "]"

	values.Add(paramFilters, filterTemplate)

	return values
}

// CreateWithContext creates a work-package or a sub-task from a JSON representation.
func (s *WorkPackageService) CreateWithContext(ctx context.Context, wpObject *WorkPackage, projectName string) (*WorkPackage, *Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/projects/%s/work_packages", projectName)
	wpResponse, resp, err := CreateWithContext(ctx, wpObject, s, apiEndpoint)
	if err != nil {
		return nil, resp, err
	}
	return wpResponse.(*WorkPackage), resp, err
}

// Create wraps CreateWithContext using the background context.
func (s *WorkPackageService) Create(wpObject *WorkPackage, projectName string) (*WorkPackage, *Response, error) {
	return s.CreateWithContext(context.Background(), wpObject, projectName)
}

// GetListWithContext will retrieve a list of work-packages using filters
func (s *WorkPackageService) GetListWithContext(ctx context.Context, options *FilterOptions, offset int, pageSize int) (*SearchResultWP, *Response, error) {
	u := url.URL{
		Path: "api/v3/work_packages",
	}

	objList, resp, err := GetListWithContext(ctx, s, u.String(), options, offset, pageSize)
	if err != nil {
		return nil, resp, err
	}
	return objList.(*SearchResultWP), resp, err
}

// GetList wraps GetListWithContext using the background context.
func (s *WorkPackageService) GetList(options *FilterOptions, offset int, pageSize int) (*SearchResultWP, *Response, error) {
	return s.GetListWithContext(context.Background(), options, offset, pageSize)
}

// DeleteWithContext will delete a single work-package.
func (s *WorkPackageService) DeleteWithContext(ctx context.Context, workpackageID string) (*Response, error) {
	apiEndPoint := fmt.Sprintf("api/v3/work_packages/%s", workpackageID)
	resp, err := DeleteWithContext(ctx, s, apiEndPoint)
	return resp, err
}

// Delete wraps DeleteWithContext using the background context.
func (s *WorkPackageService) Delete(workpackageID string) (*Response, error) {
	return s.DeleteWithContext(context.Background(), workpackageID)
}
