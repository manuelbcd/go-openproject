package openproject

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ActivitiesService handles activities for the OpenProject instance / API.
type ActivitiesService struct {
	client *Client
}

type Activity struct {
	Type      string                 `json:"_type"`
	Id        int                    `json:"id"`
	Comment   OPGenericDescription   `json:"comment"`
	Details   []OPGenericDescription `json:"details"`
	Version   int                    `json:"version"`
	CreatedAt time.Time              `json:"createdAt"`
	Links     struct {
		Self        OPGenericLink `json:"self"`
		WorkPackage OPGenericLink `json:"workPackage"`
		User        OPGenericLink `json:"user"`
		Update      OPGenericLink `json:"update"`
	} `json:"_links"`
}

type Activities struct {
	Type     string `json:"_type"`
	Total    int    `json:"total"`
	Count    int    `json:"count"`
	Embedded struct {
		Elements []Activity `json:"elements"`
	} `json:"_embedded"`
	Links struct {
		Self OPGenericLink `json:"self"`
	} `json:"_links"`
}

// GetWithContext gets activity from OpenProject using its ID
func (s *ActivitiesService) GetWithContext(ctx context.Context, activitiesID string) (*Activity, *Response, error) {
	apiEndPoint := fmt.Sprintf("api/v3/activities/%s", activitiesID)
	obj, resp, err := GetWithContext(ctx, s, apiEndPoint)
	if err != nil {
		return nil, resp, err
	}
	return obj.(*Activity), resp, err
}

// Get wraps GetWithContext using the background context.
func (s *ActivitiesService) Get(activitiesID string) (*Activity, *Response, error) {
	return s.GetWithContext(context.Background(), activitiesID)
}

// GetFromWPHrefWithContext gets activities from OpenProject using work package href string, like '/api/v3/work_packages/36353/activities'
func (s *ActivitiesService) GetFromWPHrefWithContext(ctx context.Context, href string) (*Activities, *Response, error) {
	if strings.HasPrefix(href, "/") && len(href) > 1 {
		href = href[1:]
	}
	if href == "" {
		return nil, nil, fmt.Errorf("href is empty")
	}
	apiEndPoint := href
	obj, resp, err := GetListWithContext(ctx, s, apiEndPoint, nil, 0, 0)
	if err != nil {
		return nil, resp, err
	}
	return obj.(*Activities), resp, err
}

// GetFromWPHref wraps GetFromWPHrefWithContext using the background context.
func (s *ActivitiesService) GetFromWPHref(href string) (*Activities, *Response, error) {
	return s.GetFromWPHrefWithContext(context.Background(), href)
}
