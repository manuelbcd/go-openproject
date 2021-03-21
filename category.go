package openproject

import (
	"context"
	"fmt"
)

/**
CategoryService handles categories for the OpenProject instance / API.
*/
type CategoryService struct {
	client *Client
}

// CategoryList represent a list of Projects
type CategoryList struct {
	Embedded CategoryElements `json:"_embedded,omitempty" structs:"_embedded,omitempty"`
}

type CategoryElements struct {
	Elements []Category `json:"elements,omitempty" structs:"elements,omitempty"`
}

/**
Work-package categories
*/
type Category struct {
	Type string `json:"_type,omitempty" structs:"_type,omitempty"`
	ID   int    `json:"id,omitempty" structs:"id,omitempty"`
	Name string `json:"name,omitempty" structs:"name,omitempty"`
}

/**
GetWithContext returns a single category for the given category ID.
*/
func (s *CategoryService) GetWithContext(ctx context.Context, categoryID string) (*Category, *Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/categories/%s", categoryID)
	Obj, Resp, err := GetWithContext(s, ctx, apiEndpoint)
	return Obj.(*Category), Resp, err
}

// Get wraps GetWithContext using the background context.
func (s *CategoryService) Get(categoryID string) (*Category, *Response, error) {
	return s.GetWithContext(context.Background(), categoryID)
}

/**
Retrieve category list from project with context
*/
func (s *CategoryService) GetListWithContext(ctx context.Context, projectID string) (*CategoryList, *Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/projects/%s/categories", projectID)
	Obj, Resp, err := GetListWithContext(s, ctx, apiEndpoint)
	return Obj.(*CategoryList), Resp, err
}

/**
GetList wraps GetListWithContext using the background context.
*/
func (s *CategoryService) GetList(projectID string) (*CategoryList, *Response, error) {
	return s.GetListWithContext(context.Background(), projectID)
}
