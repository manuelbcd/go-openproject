package openproject

import (
	"context"
	"fmt"
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
func (s *UserService) GetListWithContext(ctx context.Context, options *FilterOptions) (*searchResultUser, *Response, error) {
	u := url.URL{
		Path: "api/v3/users",
	}

	objList, resp, err := GetListWithContext(s, ctx, u.String(), options)
	return objList.(*searchResultUser), resp, err
}

/**
GetList wraps GetListWithContext using the background context.
*/
func (s *UserService) GetList(options *FilterOptions) (*searchResultUser, *Response, error) {
	return s.GetListWithContext(context.Background(), options)
}

/**
	CreateWithContext creates a user from a JSON representation.
**/
func (s *UserService) CreateWithContext(ctx context.Context, user *User) (*User, *Response, error) {
	apiEndpoint := "api/v3/users"
	userResponse, resp, err := CreateWithContext(s, ctx, apiEndpoint)
	return userResponse.(*User), resp, err
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
	apiEndPoint := fmt.Sprintf("api/v3/users/%s", userID)
	resp, err := DeleteWithContext(s, ctx, apiEndPoint)
	return resp, err
}

/**
Delete wraps DeleteWithContext using the background context.
*/
func (s *UserService) Delete(userID string) (*Response, error) {
	return s.DeleteWithContext(context.Background(), userID)
}
