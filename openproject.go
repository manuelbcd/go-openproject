package openproject

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// SearchOperator represents Search operators by custom type const
// Doc. https://docs.openproject.org/api/filters/#header-available-filters-1
type SearchOperator string

const (
	// Equal operator
	Equal SearchOperator = "="
	// Different operator
	Different SearchOperator = "<>"
	// GreaterThan operator
	GreaterThan SearchOperator = ">"
	// LowerThan operator
	LowerThan SearchOperator = "<"
	// SearchString operator
	SearchString SearchOperator = "**"
	// Like operator
	Like SearchOperator = "~"
	// GreaterOrEqual operator
	GreaterOrEqual SearchOperator = ">="
	// LowerOrEqual operator
	LowerOrEqual SearchOperator = "<="
)

type IPaginationResponse interface {
	TotalPage() int
	ConcatEmbed(interface{})
}

// Pagination parameters
const kOffset = "offset"
const kPageSize = "pageSize"

// Time represents the Time definition of OpenProject as a time.Time of go
type Time time.Time

// Equal compares time
func (t Time) Equal(u Time) bool {
	return time.Time(t).Equal(time.Time(u))
}

// Date represents the Date definition of OpenProject as a time.Time of go
type Date time.Time

// httpClient defines an interface for an http.Client implementation
type httpClient interface {
	Do(request *http.Request) (response *http.Response, err error)
}

// Client manages communication with the OpenProject API.
type Client struct {
	// HTTP client used to communicate with the API.
	client httpClient

	// Base URL for API requests.
	baseURL *url.URL

	// Session storage if the user authenticates with Session cookies
	session *Session

	// Services used for talking to different parts of OpenProject API.
	Authentication *AuthenticationService
	WorkPackage    *WorkPackageService
	Project        *ProjectService
	User           *UserService
	Status         *StatusService
	WikiPage       *WikiPageService
	Attachment     *AttachmentService
	Category       *CategoryService
	Query          *QueryService
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
	c.WorkPackage = &WorkPackageService{client: c}
	c.Project = &ProjectService{client: c}
	c.User = &UserService{client: c}
	c.Status = &StatusService{client: c}
	c.WikiPage = &WikiPageService{client: c}
	c.Attachment = &AttachmentService{client: c}
	c.Category = &CategoryService{client: c}
	c.Query = &QueryService{client: c}

	return c, nil
}

// NewRequestWithContext creates an API request.
// A relative URL can be provided in urlStr, in which case it is resolved relative to the baseURL of the Client.
// If specified, the value pointed to by body is JSON encoded and included as the request body.
func (c *Client) NewRequestWithContext(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	// Relative URLs should be specified without a preceding slash since baseURL will have the trailing slash
	rel.Path = strings.TrimLeft(rel.Path, "/")

	u := c.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := newRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Set authentication information
	if c.Authentication.authType == authTypeSession {
		// Set session cookie if there is one
		if c.session != nil {
			for _, cookie := range c.session.Cookies {
				req.AddCookie(cookie)
			}
		}
	} else if c.Authentication.authType == authTypeBasic {
		// Set basic auth information
		if c.Authentication.username != "" {
			req.SetBasicAuth(c.Authentication.username, c.Authentication.password)
		}
	}

	return req, nil
}

// NewRequest wraps NewRequestWithContext using the background context.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	return c.NewRequestWithContext(context.Background(), method, urlStr, body)
}

// NewMultiPartRequestWithContext creates an API request including a multi-part file.
// A relative URL can be provided in urlStr, in which case it is resolved relative to the baseURL of the Client.
// If specified, the value pointed to by buf is a multipart form.
func (c *Client) NewMultiPartRequestWithContext(ctx context.Context, method, urlStr string, buf *bytes.Buffer) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	// Relative URLs should be specified without a preceding slash since baseURL will have the trailing slash
	rel.Path = strings.TrimLeft(rel.Path, "/")

	u := c.baseURL.ResolveReference(rel)

	req, err := newRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// Set authentication information
	if c.Authentication.authType == authTypeSession {
		// Set session cookie if there is one
		if c.session != nil {
			for _, cookie := range c.session.Cookies {
				req.AddCookie(cookie)
			}
		}
	} else if c.Authentication.authType == authTypeBasic {
		// Set basic auth information
		if c.Authentication.username != "" {
			req.SetBasicAuth(c.Authentication.username, c.Authentication.password)
		}
	}

	return req, nil
}

// NewMultiPartRequest wraps NewMultiPartRequestWithContext using the background context.
func (c *Client) NewMultiPartRequest(method, urlStr string, buf *bytes.Buffer) (*http.Request, error) {
	return c.NewMultiPartRequestWithContext(context.Background(), method, urlStr, buf)
}

// Do sends an API request and returns the API response.
// The API response is JSON decoded and stored in the value pointed to by v, or returned as an error if an API error has occurred.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	httpResp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// requestDump, err := httputil.DumpResponse(httpResp, true)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(requestDump)
	// }

	err = CheckResponse(httpResp)
	if err != nil {
		// In case of error we still return the response
		return newResponse(httpResp, nil), err
	}

	if v != nil {
		// Open a NewDecoder and defer closing the reader only if there is a provided interface to decode to
		defer httpResp.Body.Close()
		err = json.NewDecoder(httpResp.Body).Decode(v)
	}

	resp := newResponse(httpResp, v)
	return resp, err
}

// Download request a file download
func (c *Client) Download(req *http.Request) (*http.Response, error) {
	httpResp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// requestDump, err := httputil.DumpResponse(httpResp, true)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(requestDump)
	// }

	err = CheckResponse(httpResp)

	return httpResp, err
}

// CheckResponse checks the API response for errors, and returns them if present.
// A response is considered an error if it has a status code outside the 200 range.
// The caller is responsible to analyze the response body.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	err := fmt.Errorf("request failed. Please analyze the request body for more details. Status code: %d", r.StatusCode)
	return err
}

// GetBaseURL will return you the Base URL.
// This is the same URL as in the NewClient constructor
func (c *Client) GetBaseURL() url.URL {
	return *c.baseURL
}

// Response represents OpenProject API response. It wraps http.Response returned from
// API and provides information about paging.
type Response struct {
	*http.Response
	Total    int
	Count    int
	PageSize int
	Offset   int
}

// New response
func newResponse(r *http.Response, v interface{}) *Response {
	resp := &Response{Response: r}
	resp.populatePageValues(v)
	return resp
}

// Sets paging values if response json was parsed to searchResult type
// (can be extended with other types if they also need paging info)
// TODO: Improve implementation to avoid redundancy without losing efficiency (reflect alternative is not efficient)
func (r *Response) populatePageValues(v interface{}) {
	switch value := v.(type) {
	case *SearchResultWP:
		r.Total = value.Total
		r.Count = value.Count
		r.PageSize = value.PageSize
		r.Offset = value.Offset
	case *SearchResultUser:
		r.Total = value.Total
		r.Count = value.Count
		r.PageSize = value.PageSize
		r.Offset = value.Offset
	case *SearchResultQuery:
		r.Total = value.Total
		r.Count = value.Count
		r.PageSize = value.PageSize
		r.Offset = value.Offset
	}
}

// BasicAuthTransport is an http.RoundTripper that authenticates all requests
// using HTTP Basic Authentication with the provided username and password.
type BasicAuthTransport struct {
	Username string
	Password string

	// Transport is the underlying HTTP transport to use when making requests.
	// It will default to http.DefaultTransport if nil.
	Transport http.RoundTripper
}

// RoundTrip implements the RoundTripper interface.  We just add the
// basic auth and return the RoundTripper for this transport type.
func (t *BasicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := cloneRequest(req) // per RoundTripper contract

	req2.SetBasicAuth(t.Username, t.Password)
	return t.transport().RoundTrip(req2)
}

// Client returns an *http.Client that makes requests that are authenticated
// using HTTP Basic Authentication.  This is a nice little bit of sugar
// so we can just get the client instead of creating the client in the calling code.
// If it's necessary to send more information on client init, the calling code can
// always skip this and set the transport itself.
func (t *BasicAuthTransport) Client() *http.Client {
	return &http.Client{Transport: t}
}

// Transport
func (t *BasicAuthTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}

// getObjectAndClient gets an inputObject (inputObject is an OpenProject object like WorkPackage, WikiPage, Status, etc.)
// and return a pointer to its Client from its service and an instance of the object itself
func getObjectAndClient(inputObj interface{}) (client *Client, resultObj interface{}) {
	switch inputObj.(type) {
	case *AttachmentService:
		client = inputObj.(*AttachmentService).client
		resultObj = new(Attachment)
	case *CategoryService:
		client = inputObj.(*CategoryService).client
		resultObj = new(Category)
	case *ProjectService:
		client = inputObj.(*ProjectService).client
		resultObj = new(Project)
	case *QueryService:
		client = inputObj.(*QueryService).client
		resultObj = new(Query)
	case *StatusService:
		client = inputObj.(*StatusService).client
		resultObj = new(Status)
	case *UserService:
		client = inputObj.(*UserService).client
		resultObj = new(User)
	case *WikiPageService:
		client = inputObj.(*WikiPageService).client
		resultObj = new(WikiPage)
	case *WorkPackageService:
		client = inputObj.(*WorkPackageService).client
		resultObj = new(WorkPackage)
	}

	return client, resultObj
}

// getObjectAndClient gets an inputObject (inputObject is an OpenProject object like WorkPackage, WikiPage, Status, etc.)
// and return a pointer to its Client from its service and an instance of the ObjectList
func getObjectListAndClient(inputObj interface{}) (client *Client, resultObjList interface{}) {
	switch inputObj.(type) {
	case *AttachmentService:
		client = inputObj.(*AttachmentService).client
		// TODO implement
	case *CategoryService:
		client = inputObj.(*CategoryService).client
		resultObjList = new(CategoryList)
	case *ProjectService:
		client = inputObj.(*ProjectService).client
		resultObjList = new(SearchResultProject)
	case *QueryService:
		client = inputObj.(*QueryService).client
		resultObjList = new(SearchResultQuery)
	case *StatusService:
		client = inputObj.(*StatusService).client
		resultObjList = new(SearchResultStatus)
	case *UserService:
		client = inputObj.(*UserService).client
		resultObjList = new(SearchResultUser)
	// WikiPage endpoint does not support POST action
	// case *WikiPageService:
	case *WorkPackageService:
		client = inputObj.(*WorkPackageService).client
		resultObjList = new(SearchResultWP)
	}

	return client, resultObjList
}

// GetWithContext (generic) retrieves object (HTTP GET verb)
// obj can be any main object (attachment, user, project, work-package, etc...) as well as response interface{}
func GetWithContext(ctx context.Context, objService interface{}, apiEndPoint string) (interface{}, *Response, error) {
	client, resultObj := getObjectAndClient(objService)
	apiEndPoint = strings.TrimRight(apiEndPoint, "/")
	if client == nil {
		return nil, nil, errors.New("Null client, object not identified")
	}

	req, err := client.NewRequestWithContext(ctx, "GET", apiEndPoint, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.Do(req, resultObj)

	if err != nil {
		return nil, resp, NewOpenProjectError(resp, err)
	}
	return resultObj, resp, nil
}

// GetListWithContext (generic) retrieves list of objects (HTTP GET verb)
// obj list is a collection of any main object (attachment, user, project, work-package, etc...) as well as response interface{}
func GetListWithContext(ctx context.Context, objService interface{}, apiEndPoint string,
	options *FilterOptions, offset int, pageSize int) (interface{}, *Response, error) {
	client, resultObjList := getObjectListAndClient(objService)
	apiEndPoint = strings.TrimRight(apiEndPoint, "/")
	req, err := client.NewRequestWithContext(ctx, "GET", apiEndPoint, nil)
	if err != nil {
		return nil, nil, err
	}

	values := make(url.Values)
	values.Add(kOffset, strconv.Itoa(offset))
	values.Add(kPageSize, strconv.Itoa(pageSize))

	if options != nil {
		values = options.prepareFilters(values)
	}
	req.URL.RawQuery = values.Encode()

	resp, err := client.Do(req, resultObjList)
	if err != nil {
		oerr := NewOpenProjectError(resp, err)
		return nil, resp, oerr
	}

	return resultObjList, resp, nil
}

// CreateWithContext (generic) creates an instance af an object (HTTP POST verb)
// Return the instance of the object rendered into proper struct as interface{} to be cast in the caller
func CreateWithContext(ctx context.Context, object interface{}, objService interface{}, apiEndPoint string) (interface{}, *Response, error) {
	client, resultObj := getObjectAndClient(objService)
	req, err := client.NewRequestWithContext(ctx, "POST", apiEndPoint, object)
	if err != nil {
		return nil, nil, err
	}
	resp, err := client.Do(req, nil)
	if err != nil {
		// incase of error return the resp for further inspection
		return nil, resp, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("could not read the returned data")
	}
	err = json.Unmarshal(data, resultObj)
	if err != nil {
		return nil, resp, fmt.Errorf("could not unmarshall the data into struct")
	}
	return resultObj, resp, nil
}

// DeleteWithContext (generic) retrieves object (HTTP DELETE verb)
// obj can be any main object (attachment, user, project, work-package, etc...)
func DeleteWithContext(ctx context.Context, objService interface{}, apiEndPoint string) (*Response, error) {
	client, _ := getObjectAndClient(objService)
	apiEndPoint = strings.TrimRight(apiEndPoint, "/")
	req, err := client.NewRequestWithContext(ctx, "DELETE", apiEndPoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req, nil)
	return resp, err
}
