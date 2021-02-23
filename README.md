# OpenProject Go Client Library

[Go](https://golang.org/) client library for [OpenProject](https://www.openproject.org) inspired in [Go Jira library](https://github.com/andygrunwald/go-jira) 


## OpenProject official API documentation
https://docs.openproject.org/api

## Usage example

Basic work-package retrieval (Single work-package with ID 36353 from community.openproject.org)

```go
import (
	"fmt"
	openproj "github.com/manuelbcd/go-openproject"
)

const openProjURL = "https://community.openproject.org/"

func main() {
	client, _ := openproj.NewClient(nil, openProjURL)
	wpResponse, _, err := client.WorkPackage.Get("36353", nil)
	if err != nil {
		panic(err)
	}

	// Output specific fields from response
	fmt.Printf("\n\nSubject: %s \nDescription: %s\n\n", wpResponse.Subject, wpResponse.Description.Raw)
}
```

## Thanks

Thank you very much [Andy Grunwald](https://github.com/andygrunwald) for the idea and your base code.
