package main

import (
	"encoding/json"
	"fmt"
	openproj "github.com/manuelbcd/go-openproject"
	"io"
)

const openProjURL = "https://community.openproject.org/"

func main() {

	client, err := openproj.NewClient(nil, openProjURL)
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	wpResponse, resp, err := client.WorkPackage.Get("36353")
	if err != nil {
		body, err := io.ReadAll(resp.Body)
		fmt.Print(string(body))
		panic(err)
	}

	// Output specific fields from response
	fmt.Printf("\n\nSubject: %s \nDescription: %s\n\n", wpResponse.Subject, wpResponse.Description.Raw)

	// Raw output of the whole object (debug only)
	fmt.Print(prettyPrint(wpResponse))
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
