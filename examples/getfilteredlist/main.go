package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	openproj "github.com/manuelbcd/go-openproject"
)

const openProjURL = "https://community.openproject.org"

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

	wpResponse, resp, err := client.WorkPackage.GetList(opt)
	if err != nil {
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Printf(string(body))
		panic(err)
	}

	// Raw output of the whole object (only for debug)
	// fmt.Printf(prettyPrint(wpResponse))

	fmt.Printf("\nWorkpackages: %d \n\n", resp.Total)

	for _, wp := range wpResponse {
		fmt.Printf("\n\nId: %d ", wp.Id)
		fmt.Printf("\nStatus: %s ", wp.Links.Status.Title)
		fmt.Printf("Subject: %.*s ", 15, wp.Subject)
		fmt.Printf("\nDescription: %.*s\n", 25, wp.Description.Raw)
	}
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
