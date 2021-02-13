package openproject

import (
	"net/http"
	"net/url"
	"strings"
)

// httpClient defines an interface for an http.Client implementation
type httpClient interface {
	Do(request *http.Request) (response *http.Response, err error)
}

// A Client manages communication with the OpenProject API.
type Client struct {
	// HTTP client used to communicate with the API.
	client httpClient

	// Base URL for API requests.
	baseURL *url.URL

	// Session storage if the user authenticates with Session cookies
	session *Session

	// Services used for talking to different parts of OpenProject API.
	Authentication   *AuthenticationService
	WorkPackage      *WPService
	Project          *ProjectService
	User             *UserService
}

// NewClient returns a new OpenProject API client.
// If a nil httpClient is provided, http.DefaultClient will be used.
func NewClient(httpClient httpClient, baseURL string) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	// ensure the baseURL contains a trailing slash so that all paths are preserved in later calls
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		client:  httpClient,
		baseURL: parsedBaseURL,
	}
	c.Authentication = &AuthenticationService{client: c}
	c.WorkPackage = &IssueService{client: c}
	c.Project = &ProjectService{client: c}
	c.User = &UserService{client: c}

	return c, nil
}