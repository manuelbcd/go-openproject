package openproject

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestQueryService_GetByID_Success(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/api/v3/queries/1"
	raw, err := ioutil.ReadFile("./mocks/get/get-query.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEdpoint)

		fmt.Fprint(w, string(raw))
	})

	if query, _, err := testClient.Query.Get("1"); err != nil {
		t.Errorf("Error given: %s", err)
	} else if query == nil {
		t.Error("Expected user. Query is nil")
	}
}

func TestQueryService_GetList_Success(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/api/v3/queries"
	raw, err := ioutil.ReadFile("./mocks/get/get-queries-no-filters.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEdpoint)

		fmt.Fprint(w, string(raw))
	})

	queries, _, err := testClient.Query.GetList()
	if queries == nil {
		t.Error("Expected query list but received nil")
	}
	if queries.Total != 25 {
		t.Error(fmt.Sprintf("Expected 25 queries in response but received %d", 25))
	}
	if queries.Embedded.Elements[0].Name != "Never" {
		errString := "Expected query name \"Never\" in pos 1 of received list"
		errString += fmt.Sprintf("\n (got \"%s\"", queries.Embedded.Elements[1].Name)
		t.Error(errString)
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}

func TestQueryService_Create(t *testing.T) {
	setup()
	defer teardown()
	raw, err := ioutil.ReadFile("./mocks/post/post-query.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc("/api/v3/queries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testRequestURL(t, r, "/api/v3/queries")

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, string(raw))
	})

	i := &Query{
		Name:    "My filter-query",
		Starred: true,
		Project: "/api/v3/projects/1",
		Columns: []OPGenericLink{
			{Href: "/api/v3/queries/columns/id"},
			{Href: "/api/v3/queries/columns/subject"},
			{Href: "/api/v3/queries/columns/status"},
		},
		SortBy: []OPGenericLink{
			{Href: "/api/v3/queries/sort_bys/id-asc"},
		},
		TimelineVisible: true,
		Hidden:          false,
		Public:          false,
		Sums:            false,
	}
	wp, _, err := testClient.Query.Create(i)
	if wp == nil {
		t.Error("Expected query object. Query object is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}

func TestQueryService_Delete(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/api/v3/queries/554", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testRequestURL(t, r, "/api/v3/queries/554")

		w.WriteHeader(http.StatusNoContent)
		fmt.Fprint(w, `{}`)
	})

	resp, err := testClient.Query.Delete("554")
	if resp.StatusCode != 204 {
		t.Error("Query not deleted.")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}
