package openproject

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestQueryService_GetByID_Success(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/api/v3/queries/1"
	raw, err := os.ReadFile("./mocks/get/get-query.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEdpoint)

		_, err := fmt.Fprint(w, string(raw))
		if err != nil {
			t.Error(err)
		}
	})

	if query, _, err := testClient.Query.Get("1"); err != nil {
		t.Errorf("Error given: %s", err)
	} else if query == nil {
		t.Error("Expected user. QueryResult is nil")
	}
}

func TestQueryService_GetList_Success(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/api/v3/queries"
	raw, err := os.ReadFile("./mocks/get/get-queries-no-filters.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEdpoint)

		_, err := fmt.Fprint(w, string(raw))
		if err != nil {
			t.Error(err)
		}
	})

	queries, _, err := testClient.Query.GetList(0, 10)
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
	raw, err := os.ReadFile("./mocks/post/post-query.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc("/api/v3/queries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testRequestURL(t, r, "/api/v3/queries")

		w.WriteHeader(http.StatusCreated)
		_, err2 := fmt.Fprint(w, string(raw))
		if err2 != nil {
			t.Error(err2)
		}
	})

	i := &QueryResult{
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
		t.Error("Expected query object. QueryResult object is nil")
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
		_, err := fmt.Fprint(w, `{}`)
		if err != nil {
			t.Error(err)
		}
	})

	resp, err := testClient.Query.Delete("554")
	if resp.StatusCode != 204 {
		t.Error("QueryResult not deleted.")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}
