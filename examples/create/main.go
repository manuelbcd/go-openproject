package main


import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"

	openproj "go-openproject"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	r := bufio.NewReader(os.Stdin)

	fmt.Print("OpenProject URL: ")
	openProjURL, _ := r.ReadString('\n')

	fmt.Print("OpenProject Username: ")
	username, _ := r.ReadString('\n')

	fmt.Print("OpenProject Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)

	tp := openproj.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	client, err := openproj.NewClient(tp.Client(), strings.TrimSpace(openProjURL))
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	i := openproj.WorkPackage{
		Subject: "This is my test work package",
		Description: &openproj.WPDescription {
			Format: "",
			Raw: "",
			Html: "Just a demo workpackage",
		},
	}

	wpResponse, _, err := client.WorkPackage.Create(&i, "demo-project")
	if err != nil {
		panic(err)
	}

	//fmt.Printf("\nSubject: %s \n StartDate: %s\n", wpResponse.Subject, wpResponse.Description.Html)

	fmt.Printf(prettyPrint(wpResponse))
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

