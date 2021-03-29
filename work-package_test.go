package openproject

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestWorkPackageService_Get_Success(t *testing.T) {
	setup()
	defer teardown()
	raw, err := ioutil.ReadFile("./mocks/get/get-workpackage.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc("/api/v3/work_packages/36350", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, "/api/v3/work_packages/36350")

		fmt.Fprint(w, string(raw))
	})

	workpkg, _, err := testClient.WorkPackage.Get("36350")
	if workpkg == nil {
		t.Error("Expected work-package. Work-package is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}

func TestWorkPackageService_Get_SearchListSuccess(t *testing.T) {
	setup()
	defer teardown()
	raw, err := ioutil.ReadFile("./mocks/get/get-workpackages-filtered.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc("/api/v3/work_packages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, "/api/v3/work_packages?filters=%5B%7B%22status%22%3A%7B%22operator%22%3A%22%3D%22%2C%22values%22%3A%5B%2221%22%5D%7D%7D%5D")

		fmt.Fprint(w, string(raw))
	})

	opt := &FilterOptions{
		Fields: []OptionsFields{
			{
				Field:    "status",
				Operator: Equal,
				Value:    "21",
			},
		},
	}

	_, resp, err := testClient.WorkPackage.GetList(opt)

	if resp == nil {
		t.Errorf("Null response: %+v", resp)
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
	if resp.Total != 2 {
		t.Errorf("Total should populate with 2, %v given", resp.Total)
	}
}

func TestWorkPackageService_Create(t *testing.T) {
	setup()
	defer teardown()
	raw, err := ioutil.ReadFile("./mocks/post/post-workpackage.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc("/api/v3/projects/demo-project/work_packages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testRequestURL(t, r, "/api/v3/projects/demo-project/work_packages")

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, string(raw))
	})

	i := &WorkPackage{
		Subject: "Just another test work-package",
		Description: &WPDescription{
			Format: "textile",
			Raw:    "This is just a demo work-package description",
		},
	}
	wp, _, err := testClient.WorkPackage.Create(i, "demo-project")
	if wp == nil {
		t.Error("Expected work-package. Work-package is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}

func TestWorkPackageService_Delete(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/api/v3/work_packages/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testRequestURL(t, r, "/api/v3/work_packages/123")

		w.WriteHeader(http.StatusNoContent)
		fmt.Fprint(w, `{}`)
	})

	resp, err := testClient.WorkPackage.Delete("123")
	if resp.StatusCode != 204 {
		t.Error("Work-package not deleted.")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}
