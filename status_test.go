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

func TestStatusService_GetList_Success(t *testing.T) {
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

	statuses, _, err := testClient.Status.GetList()
	if statuses == nil {
		t.Error("Expected status list but received nil")
	}
	if statuses.Embedded.Elements[3].Name != "needs clarification" {
		errString := "Expected status name \"needs clarifiaction\" in pos 3 of received list"
		errString += fmt.Sprintf("\n (got \"%s\"", statuses.Embedded.Elements[5].Name)
		t.Error(errString)
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}
