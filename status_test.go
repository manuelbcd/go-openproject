package openproject

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestStatusService_GetByID_Success(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/api/v3/statuses/2"
	raw, err := ioutil.ReadFile("./mocks/get/get-status.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEdpoint)

		fmt.Fprint(w, string(raw))
	})

	if user, _, err := testClient.Status.Get("2"); err != nil {
		t.Errorf("Error given: %s", err)
	} else if user == nil {
		t.Error("Expected user. User is nil")
	}
}

func TestStatusService_Get_Success(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/api/v3/statuses"
	raw, err := ioutil.ReadFile("./mocks/get/get-statuses-no-filters.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEdpoint)

		fmt.Fprint(w, string(raw))
	})

	if user, _, err := testClient.Status.Get(""); err != nil {
		t.Errorf("Error given: %s", err)
	} else if user == nil {
		t.Error("Expected user. User is nil")
	}
}
