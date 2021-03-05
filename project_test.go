package openproject

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestProjectService_GetList(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/api/v3/projects"

	raw, err := ioutil.ReadFile("./mocks/get/get-projects-no-filters.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEdpoint)
		fmt.Fprint(w, string(raw))
	})

	projects, _, err := testClient.Project.GetList()
	if projects == nil {
		t.Error("Expected project list but received nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}

func TestProjectService_Get(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/api/v3/projects/2"

	raw, err := ioutil.ReadFile("./mocks/get/get-project.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEdpoint)
		fmt.Fprint(w, string(raw))
	})

	project, _, err := testClient.Project.Get("2")
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
	if project == nil {
		t.Error("Expected project. Project is nil")
		return
	}
	if project.Name != "Scrum project" {
		t.Errorf("Unexpected project name %s", project.Name)
	}
}
