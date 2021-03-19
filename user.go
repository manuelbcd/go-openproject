package openproject

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
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
	Status    string `json:"status,omitempty" structs:"status,omitempty"`
	Language  string `json:"language,omitempty" structs:"language,omitempty"`
	Password  string `json:"password,omitempty" structs:"password,omitempty"`
}

/**
searchResult is only a small wrapper around the Search
*/
type searchResultUser struct {
	Embedded searchEmbeddedUser `json:"_embedded" structs:"_embedded"`
	Total    int                `json:"total" structs:"total"`
	Count    int                `json:"count" structs:"count"`
	PageSize int                `json:"pageSize" structs:"pageSize"`
	Offset   int                `json:"offset" structs:"offset"`
}

type searchEmbeddedUser struct {
	Elements []User `json:"elements" structs:"elements"`
}

/**
GetWithContext gets user info from OpenProject using its Account Id
// TODO: Implement GetList and adapt tests
*/
func (s *UserService) GetWithContext(ctx context.Context, accountId string) (*User, *Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/users?id=%s", accountId)
	Obj, Resp, err := GetWithContext(s, ctx, apiEndpoint)
	return Obj.(*User), Resp, err
}

/**
Get wraps GetWithContext using the background context.
*/
func (s *UserService) Get(accountId string) (*User, *Response, error) {
	return s.GetWithContext(context.Background(), accountId)
}

/**
GetListWithContext will retrieve a list of users using filters
*/
func (s *UserService) GetListWithContext(ctx context.Context, options *FilterOptions) ([]User, *Response, error) {
	u := url.URL{
		Path: "api/v3/users",
	}

	req, err := s.client.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return []User{}, nil, err
	}

	if options != nil {
		values := options.prepareFilters()
		req.URL.RawQuery = values.Encode()
	}

	v := new(searchResultUser)
	resp, err := s.client.Do(req, v)
	if err != nil {
		err = NewOpenProjectError(resp, err)
	}
	return v.Embedded.Elements, resp, err
}

/**
GetList wraps GetListWithContext using the background context.
*/
func (s *UserService) GetList(options *FilterOptions) ([]User, *Response, error) {
	return s.GetListWithContext(context.Background(), options)
}

/**
	CreateWithContext creates a user from a JSON representation.
**/
func (s *UserService) CreateWithContext(ctx context.Context, user *User) (*User, *Response, error) {
	apiEndpoint := "api/v3/users"
	req, err := s.client.NewRequestWithContext(ctx, "POST", apiEndpoint, user)
	if err != nil {
		return nil, nil, err
	}
	resp, err := s.client.Do(req, nil)
	if err != nil {
		// incase of error return the resp for further inspection
		return nil, resp, err
	}

	userResponse := new(User)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("could not read the returned data")
	}
	err = json.Unmarshal(data, userResponse)
	if err != nil {
		return nil, resp, fmt.Errorf("could not unmarshall the data into struct")
	}
	return userResponse, resp, nil
}

/**
Create wraps CreateWithContext using the background context.
*/
func (s *UserService) Create(user *User) (*User, *Response, error) {
	return s.CreateWithContext(context.Background(), user)
}

/**
DeleteWithContext will delete a single user.
*/
func (s *UserService) DeleteWithContext(ctx context.Context, userID string) (*Response, error) {
	apiEndpoint := fmt.Sprintf("api/v3/users/%s", userID)

	req, err := s.client.NewRequestWithContext(ctx, "DELETE", apiEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	return resp, err
}

/**
Delete wraps DeleteWithContext using the background context.
*/
func (s *UserService) Delete(userID string) (*Response, error) {
	return s.DeleteWithContext(context.Background(), userID)
}
