# OpenProject Go Client Library

[Go](https://golang.org/) client library for [OpenProject](https://www.openproject.org)

## API doc
https://docs.openproject.org/api

## Usage examples

### Single work-package request
Basic work-package retrieval (Single work-package with ID 36353 from community.openproject.org)
Please check [examples](https://github.com/manuelbcd/go-openproject/tree/master/examples) folder for different use-cases.

```go
import (
	"fmt"
	openproj "github.com/manuelbcd/go-openproject"
)

func main() {
	client, _ := openproj.NewClient(nil, "https://community.openproject.org/")
	wpResponse, _, err := client.WorkPackage.Get("36353", nil)
	if err != nil {
		panic(err)
	}

	// Output specific fields from response
	fmt.Printf("\n\nSubject: %s \nDescription: %s\n\n", wpResponse.Subject, wpResponse.Description.Raw)
}
```

### Create a work package
Create a single work package

```go
package main

import (
	"fmt"
	"strings"

	openproj "github.com/manuelbcd/go-openproject"
)

func main() {
	client, err := openproj.NewClient(nil, "https://youropenproject.url")
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	i := openproj.WorkPackage{
		Subject: "This is my test work package",
		Description: &openproj.WPDescription{
			Format: "textile",
			Raw:    "This is just a demo workpackage description",
		},
	}

	wpResponse, _, err := client.WorkPackage.Create(&i, "demo-project")
	if err != nil {
		panic(err)
	}

	// Output specific fields from response
	fmt.Printf("\n\nSubject: %s \nDescription: %s\n\n", wpResponse.Subject, wpResponse.Description.Raw)
}
```
## Supported objects
| Endpoint | GET single | GET many | POST single | POST many | DELETE single | DELETE many |
| ------------- | ------------- | ------------- | ------------- | ------------- | ------------- | ------------- |
| Attachments (Info) | OK | OK | *implementing* | - | *pending* | - |
| Attachments (Download) | OK | - | - | - | - | - |
| Categories | OK | OK | - | - | - | - |
| Documents | *implementing* | - | - | - | - | - |
| Projects  | OK | OK | OK | *pending* | *pending* | *pending* | *pending* |
| Queries | OK | OK | - | - | - | - |
| Schemas | *pending* |
| Statuses | OK | OK | *pending* | *pending* | *pending* | *pending* |
| Users | OK | OK | OK | | OK | *pending* |
| Wiki Pages | OK | *pending* | *pending* | *pending* | *pending* | *pending* |
| WorkPackages | OK | OK | OK | | OK | |

## Thanks
Thanks [Wieland](https://github.com/wielinde), [Oliver](https://github.com/oliverguenther) and [OpenProject](https://github.com/opf/openproject) team for your support.

Thank you very much [Andy Grunwald](https://github.com/andygrunwald) for the idea and your base code.

Inspired in [Go Jira library](https://github.com/andygrunwald/go-jira) 

