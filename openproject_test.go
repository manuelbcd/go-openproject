package openproject

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

const (
	testOpenProjectInstanceURL = "https://community.openproject.org/"
)

var (
	// testMux is the HTTP request multiplexer used with the test server.
	testMux *http.ServeMux

	// testClient is the OpenProject client being tested.
	testClient *Client

	// testServer is a test HTTP server used to provide mock API responses.
	testServer *httptest.Server
)

/**
setup sets up a test HTTP server along with a openproject.Client that is configured to talk to that test server.
Tests should register handlers on mux which provide mock responses for the API method being tested.

*/
func setup() {
	// Test server
	testMux = http.NewServeMux()
	testServer = httptest.NewServer(testMux)

	// OpenProject client configured to use test server
	testClient, _ = NewClient(nil, testServer.URL)
}

// teardown closes the test HTTP server.
func teardown() {
	testServer.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testRequestURL(t *testing.T, r *http.Request, want string) {
	if got := r.URL.String(); !strings.HasPrefix(got, want) {
		t.Errorf("Request URL: %v, want %v", got, want)
	}
}

func testRequestParams(t *testing.T, r *http.Request, want map[string]string) {
	params := r.URL.Query()

	if len(params) != len(want) {
		t.Errorf("Request params: %d, want %d", len(params), len(want))
	}

	for key, val := range want {
		if got := params.Get(key); val != got {
			t.Errorf("Request params: %s, want %s", got, val)
		}

	}

}

func TestNewClient_WrongUrl(t *testing.T) {
	c, err := NewClient(nil, "://community.openproject.org")

	if err == nil {
		t.Error("Expected an error. Got none")
	}
	if c != nil {
		t.Errorf("Expected no client. Got %+v", c)
	}
}

func TestNewClient_WithHttpClient(t *testing.T) {
	httpClient := http.DefaultClient
	httpClient.Timeout = 10 * time.Minute

	c, err := NewClient(httpClient, testOpenProjectInstanceURL)
	if err != nil {
		t.Errorf("Got an error: %s", err)
	}
	if c == nil {
		t.Error("Expected a client. Got none")
		return
	}
	if !reflect.DeepEqual(c.client, httpClient) {
		t.Errorf("HTTP clients are not equal. Injected %+v, got %+v", httpClient, c.client)
	}
}

func TestNewClient_WithServices(t *testing.T) {
	c, err := NewClient(nil, testOpenProjectInstanceURL)

	if err != nil {
		t.Errorf("Got an error: %s", err)
	}
	if c.Authentication == nil {
		t.Error("No AuthenticationService provided")
	}
	if c.WorkPackage == nil {
		t.Error("No WorkpackageService provided")
	}
	if c.Project == nil {
		t.Error("No ProjectService provided")
	}
	if c.User == nil {
		t.Error("No UserService provided")
	}
}

func TestCheckResponse(t *testing.T) {
	codes := []int{
		http.StatusOK, http.StatusPartialContent, 299,
	}

	for _, c := range codes {
		r := &http.Response{
			StatusCode: c,
		}
		if err := CheckResponse(r); err != nil {
			t.Errorf("CheckResponse throws an error: %s", err)
		}
	}
}

func TestClient_NewRequest(t *testing.T) {
	c, err := NewClient(nil, testOpenProjectInstanceURL)
	if err != nil {
		t.Errorf("An error occurred. Expected nil. Got %+v.", err)
	}

	inURL, outURL := "api/v3/workpackages/", testOpenProjectInstanceURL+"api/v3/workpackages/"
	inBody, outBody := &WorkPackage{ID: 1}, `{"id":1}`+"\n"
	req, _ := c.NewRequest("GET", inURL, inBody)

	// Test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// Test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%v) Body is %v, want %v", inBody, got, want)
	}
}

func TestClient_NewRawRequest(t *testing.T) {
	c, err := NewClient(nil, testOpenProjectInstanceURL)
	if err != nil {
		t.Errorf("An error occurred. Expected nil. Got %+v.", err)
	}

	inURL, outURL := "api/v3/workpackages/", testOpenProjectInstanceURL+"api/v3/workpackages/"

	outBody := `{"id":1}` + "\n"
	inBody := outBody
	req, _ := c.NewRawRequest("GET", inURL, strings.NewReader(outBody))

	// Test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRawRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// Test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRawRequest(%v) Body is %v, want %v", inBody, got, want)
	}
}

func testURLParseError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func TestClient_NewRequest_BadURL(t *testing.T) {
	c, err := NewClient(nil, testOpenProjectInstanceURL)
	if err != nil {
		t.Errorf("An error occurred. Expected nil. Got %+v.", err)
	}
	_, err = c.NewRequest("GET", ":", nil)
	testURLParseError(t, err)
}

func TestClient_NewRequest_BasicAuth(t *testing.T) {
	c, err := NewClient(nil, testOpenProjectInstanceURL)
	if err != nil {
		t.Errorf("An error occurred. Expected nil. Got %+v.", err)
	}

	c.Authentication.SetBasicAuth("test-user", "test-password")

	inURL := "api/v3/"
	inBody := &WorkPackage{ID: 1}
	req, err := c.NewRequest("GET", inURL, inBody)

	if err != nil {
		t.Errorf("An error occurred. Expected nil. Got %+v.", err)
	}

	username, password, ok := req.BasicAuth()
	if !ok || username != "test-user" || password != "test-password" {
		t.Errorf("An error occurred. Expected basic auth username %s and password %s. Got username %s and password %s.", "test-user", "test-password", username, password)
	}
}

func TestClient_NewRequest_EmptyBody(t *testing.T) {
	c, err := NewClient(nil, testOpenProjectInstanceURL)
	if err != nil {
		t.Errorf("An error occurred. Expected nil. Got %+v.", err)
	}
	req, err := c.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}
	if req.Body != nil {
		t.Fatalf("constructed request contains a non-nil Body")
	}
}

func TestClient_NewMultiPartRequest(t *testing.T) {
	c, err := NewClient(nil, testOpenProjectInstanceURL)
	if err != nil {
		t.Errorf("An error occurred. Expected nil. Got %+v.", err)
	}

	cookie := &http.Cookie{Name: "testcookie", Value: "testvalue"}
	c.session = &Session{Cookies: []*http.Cookie{cookie}}
	c.Authentication.authType = authTypeSession

	inURL := "api/v3/workpackage/"
	inBuf := bytes.NewBufferString("teststring")
	req, err := c.NewMultiPartRequest("GET", inURL, inBuf)

	if err != nil {
		t.Errorf("An error occurred. Expected nil. Got %+v.", err)
	}

	if len(req.Cookies()) != len(c.session.Cookies) {
		t.Errorf("An error occurred. Expected %d cookie(s). Got %d.", len(c.session.Cookies), len(req.Cookies()))
	}

	for i, v := range req.Cookies() {
		if v.String() != c.session.Cookies[i].String() {
			t.Errorf("An error occurred. Unexpected cookie. Expected %s, actual %s.", v.String(), c.session.Cookies[i].String())
		}
	}

	if req.Header.Get("X-Atlassian-Token") != "nocheck" {
		t.Errorf("An error occurred. Unexpected X-Atlassian-Token header value. Expected nocheck, actual %s.", req.Header.Get("X-Atlassian-Token"))
	}
}

func TestClient_NewMultiPartRequest_BasicAuth(t *testing.T) {
	c, err := NewClient(nil, testOpenProjectInstanceURL)
	if err != nil {
		t.Errorf("An error occurred. Expected nil. Got %+v.", err)
	}

	c.Authentication.SetBasicAuth("test-user", "test-password")

	inURL := "api/v3/workpackages/"
	inBuf := bytes.NewBufferString("teststring")
	req, err := c.NewMultiPartRequest("GET", inURL, inBuf)

	if err != nil {
		t.Errorf("An error occurred. Expected nil. Got %+v.", err)
	}

	username, password, ok := req.BasicAuth()
	if !ok || username != "test-user" || password != "test-password" {
		t.Errorf("An error occurred. Expected basic auth username %s and password %s. Got username %s and password %s.", "test-user", "test-password", username, password)
	}
}

func TestClient_Do(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := testClient.NewRequest("GET", "/", nil)
	body := new(foo)
	testClient.Do(req, body)

	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestClient_Do_HTTPResponse(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := testClient.NewRequest("GET", "/", nil)
	res, _ := testClient.Do(req, nil)
	_, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Errorf("Error on parsing HTTP Response = %v", err.Error())
	} else if res.StatusCode != 200 {
		t.Errorf("Response code = %v, want %v", res.StatusCode, 200)
	}
}

func TestClient_Do_HTTPError(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := testClient.NewRequest("GET", "/", nil)
	_, err := testClient.Do(req, nil)

	if err == nil {
		t.Error("Expected HTTP 400 error.")
	}
}

func TestClient_Do_RedirectLoop(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	})

	req, _ := testClient.NewRequest("GET", "/", nil)
	_, err := testClient.Do(req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected a URL error; got %+v.", err)
	}
}

func TestClient_GetBaseURL_WithURL(t *testing.T) {
	u, err := url.Parse(testOpenProjectInstanceURL)
	if err != nil {
		t.Errorf("URL parsing -> Got an error: %s", err)
	}

	c, err := NewClient(nil, testOpenProjectInstanceURL)
	if err != nil {
		t.Errorf("Client creation -> Got an error: %s", err)
	}
	if c == nil {
		t.Error("Expected a client. Got none")
	}

	if b := c.GetBaseURL(); !reflect.DeepEqual(b, *u) {
		t.Errorf("Base URLs are not equal. Expected %+v, got %+v", *u, b)
	}
}

func TestBasicAuthTransport(t *testing.T) {
	setup()
	defer teardown()

	username, password := "username", "password"

	testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok {
			t.Errorf("request does not contain basic auth credentials")
		}
		if u != username {
			t.Errorf("request contained basic auth username %q, want %q", u, username)
		}
		if p != password {
			t.Errorf("request contained basic auth password %q, want %q", p, password)
		}
	})

	tp := &BasicAuthTransport{
		Username: username,
		Password: password,
	}

	basicAuthClient, _ := NewClient(tp.Client(), testServer.URL)
	req, _ := basicAuthClient.NewRequest("GET", ".", nil)
	basicAuthClient.Do(req, nil)
}

func TestBasicAuthTransport_transport(t *testing.T) {
	// default transport
	tp := &BasicAuthTransport{}
	if tp.transport() != http.DefaultTransport {
		t.Errorf("Expected http.DefaultTransport to be used.")
	}

	// custom transport
	tp = &BasicAuthTransport{
		Transport: &http.Transport{},
	}
	if tp.transport() == http.DefaultTransport {
		t.Errorf("Expected custom transport to be used.")
	}
}
