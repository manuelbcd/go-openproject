package openproject

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestCategoryService_Get(t *testing.T) {
	setup()
	defer teardown()
	testAPIEndpoint := "/api/v3/categories/1"

	raw, err := ioutil.ReadFile("./mocks/get/get-category.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEndpoint)
		fmt.Fprint(w, string(raw))
	})

	category, _, err := testClient.Category.Get("1")
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
	if category == nil {
		t.Error("Expected category. Category is nil")
		return
	}
	if category.Name != "Category 1 (to be changed in Project settings)" {
		t.Errorf("Unexpected category name %s", category.Name)
	}
}

func TestCategoryService_GetList(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/api/v3/projects/demo-project/categories/"

	raw, err := ioutil.ReadFile("./mocks/get/get-categories-from-project-no-filters.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEdpoint)
		fmt.Fprint(w, string(raw))
	})

	categories, _, err := testClient.Category.GetList("demo-project", 0, 10)
	if categories == nil {
		t.Error("Expected category list from project, but received nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}
