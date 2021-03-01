package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	openproj "github.com/manuelbcd/go-openproject"
)

const openProjURL = "https://community.openproject.org/"

func main() {

	client, err := openproj.NewClient(nil, openProjURL)
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	opt := &openproj.FilterOptions{
		Fields: []openproj.OptionsFields{
			{
				Field:    "status",
				Operator: openproj.EQUAL,
				Value:    "21",
			},
		},
	}

	wpResponse, resp, err := client.WorkPackage.Get("", opt)
	if err != nil {
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Printf(string(body))
		panic(err)
	}

	// Output specific fields from response
	//fmt.Printf("\n\nSubject: %s \nDescription: %s\n\n", wpResponse.Subject, wpResponse.Description.Raw)

	// Raw output of the whole object (only for debug)
	fmt.Printf(prettyPrint(wpResponse))
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
