package openproject

import (
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
