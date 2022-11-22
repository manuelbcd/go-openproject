package openproject

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestUserService_Get_SearchListNoFiltersSuccess(t *testing.T) {
	setup()
	defer teardown()
	raw, err := ioutil.ReadFile("./mocks/get/get-users-no-filters.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc("/api/v3/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, "/api/v3/users")

		fmt.Fprint(w, string(raw))
	})

	userList, resp, err := testClient.User.GetList(nil, 0, 10)

	if resp == nil {
		t.Errorf("Null response: %+v", resp)
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
	if userList.Total != 2 {
		t.Errorf("Total should populate with 2, %v given", resp.Total)
	}
	if userList.Embedded.Elements[0].FirstName != "John" {
		t.Errorf("Expected John as firstName of element 0 but received %s", userList.Embedded.Elements[0].FirstName)
	}
}

func TestUserService_GetByID_Success(t *testing.T) {
	setup()
	defer teardown()
	raw, err := ioutil.ReadFile("./mocks/get/get-user.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc("/api/v3/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, "/api/v3/users?id=1")

		fmt.Fprint(w, string(raw))
	})

	if user, _, err := testClient.User.Get("1"); err != nil {
		t.Errorf("Error given: %s", err)
	} else if user == nil {
		t.Error("Expected user. User is nil")
	}
}

func TestUserService_Create(t *testing.T) {
	setup()
	defer teardown()
	raw, err := ioutil.ReadFile("./mocks/post/post-user.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc("/api/v3/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testRequestURL(t, r, "/api/v3/users")

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, string(raw))
	})

	p := &User{
		Type:      "User",
		Login:     "john.smith@acme.com",
		Admin:     false,
		FirstName: "John",
		LastName:  "Smith",
		Email:     "john.smith@acme.com",
		Status:    "active",
		Language:  "en",
		Password:  "AB12345pass",
	}
	wp, _, err := testClient.User.Create(p)
	if wp == nil {
		t.Error("Expected user but user is nil")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}

func TestUserService_Delete(t *testing.T) {
	setup()
	defer teardown()
	testMux.HandleFunc("/api/v3/users/4", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testRequestURL(t, r, "/api/v3/users/4")

		w.WriteHeader(http.StatusNoContent)
		fmt.Fprint(w, `{}`)
	})

	resp, err := testClient.User.Delete("4")
	if resp.StatusCode != 204 {
		t.Error("User not deleted.")
	}
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
}
