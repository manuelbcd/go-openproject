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
	if attachment.FileName != "Leerzeile.png" {
		t.Errorf("Unexpected attachment filename %s", attachment.FileName)
	}
}

func TestAttachmentService_Download(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/api/v3/attachments/5/content"

	raw, err := ioutil.ReadFile("./mocks/get/download-attachment-file.jpg")
	if err != nil {
		t.Error(err.Error())
	}
	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEdpoint)
		fmt.Fprint(w, string(raw))
	})

	file, err := testClient.Attachment.Download("5")
	if err != nil {
		t.Errorf("Error given: %s", err)
	}
	if file == nil {
		t.Error("Expected file (download). File is nil")
		return
	}
	if len(*file) != 111186 {
		t.Errorf("Unexpected downloaded filesize %d", len(*file))
	}
}

func TestAttachmentService_Download_BadStatus(t *testing.T) {
	setup()
	defer teardown()
	testAPIEdpoint := "/api/v3/attachments/5/content"

	testMux.HandleFunc(testAPIEdpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, testAPIEdpoint)
		w.WriteHeader(http.StatusForbidden)
	})

	_, err := testClient.Attachment.Download("5")

	if err == nil {
		t.Errorf("Error %d expected but null received", http.StatusForbidden)
	}
}
