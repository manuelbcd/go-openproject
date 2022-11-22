package openproject

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
)

// ProjectService handles projects for the OpenProject instance / API.
type ProjectService struct {
	client *Client
}

// SearchResultProject represent a list of Projects
type SearchResultProject struct {
	Embedded projectElements `json:"_embedded,omitempty" structs:"_embedded,omitempty"`
	PaginationParam
}

func (s *SearchResultProject) TotalPage() int {
	return int(math.Ceil(float64(s.Total) / float64(s.PageSize)))
}

func (s *SearchResultProject) ConcatEmbed(proj interface{}) {
	s.Embedded.Elements = append(s.Embedded.Elements, proj.(*SearchResultProject).Embedded.Elements...)
}

// ProjectElements represent elements within SearchResultProject
type projectElements struct {
	Elements []Project `json:"elements,omitempty" structs:"elements,omitempty"`
}

// Project structure representing OpenProject project
type Project struct {
	Type              string           `json:"_type,omitempty" structs:"_type,omitempty"`
	ID                int              `json:"id,omitempty" structs:"id,omitempty"`
	Identifier        string           `json:"identifier,omitempty" structs:"identifier,omitempty"`
	Name              string           `json:"name,omitempty" structs:"name,omitempty"`
	Active            bool             `json:"active,omitempty" structs:"active,omitempty"`
	Public            bool             `json:"public,omitempty" structs:"public,omitempty"`
	Description       *ProjDescription `json:"description,omitempty" structs:"description,omitempty"`
	CreatedAt         *Time            `json:"createdAt,omitempty" structs:"createdAt,omitempty"`
	UpdatedAt         *Time            `json:"updatedAt,omitempty" structs:"updatedAt,omitempty"`
	Status            string           `json:"status,omitempty" structs:"status,omitempty"`
	StatusExplanation *ProjDescription `json:"statusExplanation,omitempty" structs:"statusExplanation,omitempty"`
	Links             struct {
		Self                         *OPGenericLink `json:"self,omitempty" structs:"self,omitempty"`
		CreateWorkPackage            *OPGenericLink `json:"createWorkPackage,omitempty" structs:"createWorkPackage,omitempty"`
		CreateWorkPackageImmediately *OPGenericLink `json:"createWorkPackageImmediately,omitempty" structs:"createWorkPackageImmediately,omitempty"`
		WorkPackages                 *OPGenericLink `json:"workPackages,omitempty" structs:"workPackages,omitempty"`
		Categories                   *OPGenericLink `json:"categories,omitempty" structs:"categories,omitempty"`
		Versions                     *OPGenericLink `json:"versions,omitempty" structs:"versions,omitempty"`
		Memberships                  *OPGenericLink `json:"memberships,omitempty" structs:"memberships,omitempty"`
		Types                        *OPGenericLink `json:"types,omitempty" structs:"types,omitempty"`
		Update                       *OPGenericLink `json:"update,omitempty" structs:"update,omitempty"`
		UpdateImmediately            *OPGenericLink `json:"updateImmediately,omitempty" structs:"updateImmediately,omitempty"`
		Delete                       *OPGenericLink `json:"delete,omitempty" structs:"delete,omitempty"`
		Schema                       *OPGenericLink `json:"schema,omitempty" structs:"schema,omitempty"`
		Status                       *OPGenericLink `json:"status,omitempty" structs:"status,omitempty"`
		Ancestors                    []interface{}  `json:"ancestors,omitempty"`
		Parent                       *OPGenericLink `json:"parent,omitempty" structs:"parent,omitempty"`
	} `json:"_links,omitempty" structs:"_links,omitempty"`
}

// ProjDescription type contains description and format
type ProjDescription OPGenericDescription

// GetWithContext returns a single project for the given project key.
func (s *ProjectService) GetWithContext(ctx context.Context, projectID string) (*Project, *Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/projects/%s", projectID)
	Obj, Resp, err := GetWithContext(ctx, s, apiEndpoint)
	if err != nil {
		return nil, Resp, err
	}
	return Obj.(*Project), Resp, err
}

// Get wraps GetWithContext using the background context.
func (s *ProjectService) Get(projectID string) (*Project, *Response, error) {
	return s.GetWithContext(context.Background(), projectID)
}

// GetList wraps GetListWithContext using the background context.
func (s *ProjectService) GetList(offset int, pageSize int) (*SearchResultProject, *Response, error) {
	return s.GetListWithContext(context.Background(), offset, pageSize)
}

// GetListWithContext retrieve project list with context
// TODO: Implement search options
func (s *ProjectService) GetListWithContext(ctx context.Context, offset int, pageSize int) (*SearchResultProject, *Response, error) {
	apiEndpoint := "api/v3/projects"
	Obj, Resp, err := GetListWithContext(ctx, s, apiEndpoint, nil, offset, pageSize)
	if err != nil {
		return nil, Resp, err
	}
	return Obj.(*SearchResultProject), Resp, err
}

// CreateWithContext creates a project from a JSON representation.
func (s *ProjectService) CreateWithContext(ctx context.Context, project *Project) (*Project, *Response, error) {
	apiEndpoint := "api/v3/projects"
	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, project)
	if err != nil {
		return nil, nil, err
	}
	resp, err := s.client.Do(req, nil)
	if err != nil {
		// incase of error return the resp for further inspection
		return nil, resp, err
	}

	projResponse := new(Project)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("could not read the returned data")
	}
	err = json.Unmarshal(data, projResponse)
	if err != nil {
		return nil, resp, fmt.Errorf("could not unmarshall the data into struct")
	}
	return projResponse, resp, nil
}

// Create wraps CreateWithContext using the background context.
func (s *ProjectService) Create(project *Project) (*Project, *Response, error) {
	return s.CreateWithContext(context.Background(), project)
}
