package main

import (
	"encoding/json"
	"fmt"
	openproj "github.com/manuelbcd/go-openproject"
	"io"
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
				Operator: openproj.Equal,
				Value:    "21",
			},
		},
	}

	wpResponse, resp, err := client.WorkPackage.GetList(opt, 0, 10)
	if err != nil {
		body, err := io.ReadAll(resp.Body)
		fmt.Print(string(body))
		panic(err)
	}

	// Raw output of the whole object (only for debug)
	// fmt.Printf(prettyPrint(wpResponse))

	fmt.Printf("\nWorkpackages: %d \n\n", resp.Total)

	for _, wp := range wpResponse.Embedded.Elements {
		fmt.Printf("\n\nId: %d ", wp.ID)
		fmt.Printf("\nStatus: %s ", wp.Links.Status.Title)
		fmt.Printf("Subject: %.*s ", 15, wp.Subject)
		fmt.Printf("\nDescription: %.*s\n", 25, wp.Description.Raw)
	}

	// Or you can use auto page turn
	// Use careful when dealing large amounts of data because it will set all objects in memory.
	allUser, err := openproj.AutoPageTurn(opt, 20, client.User.GetList)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n%s", prettyPrint(allUser))
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
