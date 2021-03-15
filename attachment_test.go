package openproject

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAttachmentService_Get(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/api/v3/attachments/5"

	raw, err := ioutil.ReadFile("./mocks/get/get-attachment-from-workpackage.json")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEdpoint)
		fmt.Fprint(w, string(raw))
	})

	attachment, _, err := testClient.Attachment.Get("5")
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
	if attachment == nil {
		t.Error("Expected attachment. Attachment is nil")
		return
	}
	if attachment.FileName != "Rollen-Ticket-Sichtbarkeit.jpg" {
		t.Errorf("Unexpected attachment filename %s", attachment.FileName)
	}
}
