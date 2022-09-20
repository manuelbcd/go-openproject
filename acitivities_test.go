package openproject

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestActivitiesServiceGet(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/api/v3/activities/357978"

	raw, err := os.ReadFile("./mocks/get/get-activity.json")
	if err != nil {
		t.Error(err.Error())
		return
	}
	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEdpoint)
		fmt.Fprint(w, string(raw))
	})

	activity, _, err := testClient.Activities.Get("357978")
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
	if activity == nil {
		t.Error("Expected activities. Activities is nil")
		return
	}
	if activity.Id != 357978 {
		t.Errorf("Unexpected activity id %d", activity.Id)
	}
}

func TestActivitiesService_GetFromWPHref(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/api/v3/work_packages/36353/activities"

	raw, err := os.ReadFile("./mocks/get/get-activities.json")
	if err != nil {
		t.Error(err.Error())
		return
	}

	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEdpoint)
		fmt.Fprint(w, string(raw))
	})

	activities, _, err := testClient.Activities.GetFromWPHref("/api/v3/work_packages/36353/activities")
	if err != nil {
		t.Errorf("Error given: %s", err)
		return
	}
	if activities == nil {
		t.Error("Expected activities. Activities is nil")
		return
	}
	if activities.Count != 6 {
		t.Errorf("Unexpected activities count %d", activities.Count)
	}
}
