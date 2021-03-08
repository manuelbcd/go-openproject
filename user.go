package openproject

import (
	"context"
	"fmt"
)

/**
UserService handles users for the OpenProject instance / API.
*/
type UserService struct {
	client *Client
}

/**
User is the object representing OpenProject users.
TODO: Complete object with fields identityUrl, language, _links
*/
type User struct {
	Type      string `json:"_type,omitempty" structs:"_type,omitempty"`
	Id        int    `json:"id,omitempty" structs:"id,omitempty"`
	Name      string `json:"name,omitempty" structs:"name,omitempty"`
	CreatedAt *Time  `json:"createdAt,omitempty" structs:"createdAt,omitempty"`
	UpdatedAt *Time  `json:"updatedAt,omitempty" structs:"updatedAt,omitempty"`
	Login     string `json:"login,omitempty" structs:"login,omitempty"`
	Admin     bool   `json:"admin,omitempty" structs:"admin,omitempty"`
	FirstName string `json:"firstName,omitempty" structs:"firstName,omitempty"`
	lastName  string `json:"lastName,omitempty" structs:"lastName,omitempty"`
	Email     string `json:"email,omitempty" structs:"email,omitempty"`
	Avatar    string `json:"avatar,omitempty" structs:"avatar,omitempty"`
	Status    bool   `json:"status,omitempty" structs:"status,omitempty"`
}

/**
GetWithContext gets user info from OpenProject using its Account Id
// TODO: Implement GetList and adapt tests
*/
func (s *UserService) GetWithContext(ctx context.Context, accountId string) (*User, *Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/users?id=%s", accountId)
	req, err := s.client.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, NewOpenProjectError(resp, err)
	}
	return user, resp, nil
}

/**
Get wraps GetWithContext using the background context.
*/
func (s *UserService) Get(accountId string) (*User, *Response, error) {
	return s.GetWithContext(context.Background(), accountId)
}
