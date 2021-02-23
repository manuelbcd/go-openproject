# OpenProject Go Client Library

[Go](https://golang.org/) client library for [OpenProject](https://www.openproject.org) inspired (more than inspired) in [Go Jira library](https://github.com/andygrunwald/go-jira) 


## OpenProject official API documentation
https://docs.openproject.org/api

## Usage example

```go
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

	wpResponse, _, err := client.WorkPackage.Get("36353", nil)
	if err != nil {
		panic(err)
	}

	// Output specific fields from response
	fmt.Printf("\n\nSubject: %s \nDescription: %s\n\n", wpResponse.Subject, wpResponse.Description.Raw)
}
```

## Thanks

Thank you very much Andy Grunwald for the inspiration
