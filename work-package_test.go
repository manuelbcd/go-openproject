package openproject

import (
	"encoding/json"
	"fmt"
	"github.com/trivago/tgo/tcontainer"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"
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

func TestWorkPackageFields_MarshalJSON_OmitsEmptyFields(t *testing.T) {
	thisDate, _ := time.Parse("2006-01-02T15:04:05-0700", "2020-12-31")
	thisDateTime := Time(thisDate)

	i := &WorkPackage{
		Description: &WPDescription{
			Format: "html",
			Raw:    "Content example",
			HTML:   "<html>Content example</html>",
		},
		CreatedAt: &thisDateTime,
		Type:      "Task",
	}

	rawdata, err := json.Marshal(i)
	if err != nil {
		t.Errorf("Expected nil err, received %s", err)
	}

	// convert json to map and see if unset keys are there
	issuef := tcontainer.NewMarshalMap()
	err = json.Unmarshal(rawdata, &issuef)
	if err != nil {
		t.Errorf("Expected nil err, received %s", err)
	}

	_, err = issuef.Int("description/html")
	if err == nil {
		t.Error("Expected non nil error, received nil")
	}

	// verify that the field that should be there, is.
	name, err := issuef.String("description/raw")
	if err != nil {
		t.Errorf("Expected nil err, received %s", err)
	}

	if name != "Content example" {
		t.Errorf("Expected Story, received %s", name)
	}

}

func TestWorkPackageCustomFields_MarshalJSON_Success(t *testing.T) {
	thisDate, _ := time.Parse("2006-01-02T15:04:05-0700", "2020-12-31")
	thisDateTime := Time(thisDate)

	i := &WorkPackage{
		Description: &WPDescription{
			Format: "html",
			Raw:    "Content example",
			HTML:   "<html>Content example</html>",
		},
		CreatedAt: &thisDateTime,
		Custom: tcontainer.MarshalMap{
			"customfield_A": "testA",
		},
	}

	bytes, err := json.Marshal(i)
	if err != nil {
		t.Errorf("Expected nil err, received %s", err)
	}

	received := new(WorkPackage)
	// the order of json might be different. so unmarshal it again and compare objects
	err = json.Unmarshal(bytes, received)
	if err != nil {
		t.Errorf("Expected nil err, received %s", err)
	}

	if !reflect.DeepEqual(i, received) {
		t.Errorf("Received object different from expected. Expected %+v, received %+v", i, received)
	}
}
