package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	openproj "github.com/manuelbcd/go-openproject"
)

const openProjURL = "https://community.openproject.org"

func main() {

	client, err := openproj.NewClient(nil, openProjURL)
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	// Get attachment info
	attachmentResp, resp, err := client.Attachment.Get("15713")
	if err != nil {
		body, err := io.ReadAll(resp.Body)
		fmt.Print(string(body))
		panic(err)
	}

	fmt.Printf("\nAttachment: %s \n\n", attachmentResp.FileName)
	// Raw output of the whole object (debug only)
	fmt.Print(prettyPrint(attachmentResp))

	// Download attachment file
	attachmentFile, err := client.Attachment.Download("15713")
	if err != nil {
		body, err := io.ReadAll(resp.Body)
		fmt.Print(string(body))
		panic(err)
	}

	// write the whole file at once
	file, err := os.Create(attachmentResp.FileName)
	if err != nil {
		return
	}
	defer file.Close()
	file.Write(*attachmentFile)
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
