package openproject

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func TestAutoPageTurn(t *testing.T) {
	setup()
	defer teardown()

	endpoint := "/api/v3/users"
	testMux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, endpoint)
		offsetStr := r.URL.Query().Get("offset")
		if offsetStr == "" {
			t.Error("Expected offset query parameter")
		}
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			t.Errorf("Expected offset to be an integer, %s given", offsetStr)
		}
		pageSizeStr := r.URL.Query().Get("pageSize")
		if pageSizeStr == "" {
			t.Error("Expected pageSize query parameter")
		}
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			t.Errorf("Expected pageSize to be an integer, %s given", pageSizeStr)
		}
		if pageSize != 10 {
			t.Errorf("Expected pageSize to be 10, %d given", pageSize)
		}
		mockJsonPath := fmt.Sprintf("./mocks/get/get-users-pagination-%d.json", offset-1)
		t.Logf("mock data from %s", mockJsonPath)
		raw, _ := os.ReadFile(mockJsonPath)
		fmt.Fprint(w, string(raw))
	})

	users, err := AutoPageTurn(nil, 10, testClient.User.GetList)
	if err != nil {
		t.Errorf("Error given: %s", err)
		return
	}
	if users == nil {
		t.Error("Expected user list but received nil")
		return
	}
	// check total in data and the length of elements
	totalInJson := users.Total
	lenOfElement := len(users.Embedded.Elements)
	if totalInJson != lenOfElement {
		t.Errorf("Expected total in json %d is not equal the length of elements %d.", totalInJson, lenOfElement)
	}
}
