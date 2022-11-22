package openproject

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestProjectService_GetList(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/api/v3/projects"

	raw, err := os.ReadFile("./mocks/get/get-projects-no-filters.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEdpoint)
		fmt.Fprint(w, string(raw))
	})

	projects, _, err := testClient.Project.GetList(0, 10)
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

	raw, err := os.ReadFile("./mocks/get/get-project.json")
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

func TestProjectService_Create(t *testing.T) {
	setup()
	defer teardown()
	raw, err := os.ReadFile("./mocks/post/post-project.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc("/api/v3/projects", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testRequestURL(t, r, "/api/v3/projects")

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, string(raw))
	})

	p := &Project{
		Type:   "Project",
		Name:   "Scrum project 1",
		Active: true,
		Public: true,
		Description: &ProjDescription{
			Format: "textile",
			Raw:    "This is a short summary of the goals of another demo Scrum project.",
			HTML:   "<p class=\"op-uc-p\">This is a short summary of the goals of another demo Scrum project.</p>",
		},
	}
	wp, _, err := testClient.Project.Create(p)
	if wp == nil {
		t.Error("Expected project but project is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}
