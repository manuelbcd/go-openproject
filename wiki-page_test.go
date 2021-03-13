package openproject

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestWikiPageService_Get(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/api/v3/wiki_pages/1"

	raw, err := ioutil.ReadFile("./mocks/get/get-wikipage.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEdpoint)
		fmt.Fprint(w, string(raw))
	})

	wiki, _, err := testClient.WikiPage.Get("1")
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
	if wiki == nil {
		t.Error("Expected wiki page. Wiki page is nil")
		return
	}
	if wiki.Title != "This is a wiki page" {
		t.Errorf("Unexpected wiki title %s", wiki.Title)
	}
}