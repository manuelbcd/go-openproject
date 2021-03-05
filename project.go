package openproject

/**
ProjectService handles projects for the OpenProject instance / API.
*/
type ProjectService struct {
	client *Client
}

// ProjectList represent a list of Projects
type ProjectList []struct {
	Expand          string          `json:"expand" structs:"expand"`
	Self            string          `json:"self" structs:"self"`
	ID              string          `json:"id" structs:"id"`
	Key             string          `json:"key" structs:"key"`
	Name            string          `json:"name" structs:"name"`
	AvatarUrls      AvatarUrls      `json:"avatarUrls" structs:"avatarUrls"`
	ProjectTypeKey  string          `json:"projectTypeKey" structs:"projectTypeKey"`
	ProjectCategory ProjectCategory `json:"projectCategory,omitempty" structs:"projectsCategory,omitempty"`
	IssueTypes      []IssueType     `json:"issueTypes,omitempty" structs:"issueTypes,omitempty"`
}

/** Project structure representing OpenProject project */
type Project struct {
	Expand          string             `json:"expand,omitempty" structs:"expand,omitempty"`
	Self            string             `json:"self,omitempty" structs:"self,omitempty"`
	ID              string             `json:"id,omitempty" structs:"id,omitempty"`
	Key             string             `json:"key,omitempty" structs:"key,omitempty"`
	Description     string             `json:"description,omitempty" structs:"description,omitempty"`
	Lead            User               `json:"lead,omitempty" structs:"lead,omitempty"`
	Components      []ProjectComponent `json:"components,omitempty" structs:"components,omitempty"`
	IssueTypes      []IssueType        `json:"issueTypes,omitempty" structs:"issueTypes,omitempty"`
	URL             string             `json:"url,omitempty" structs:"url,omitempty"`
	Email           string             `json:"email,omitempty" structs:"email,omitempty"`
	AssigneeType    string             `json:"assigneeType,omitempty" structs:"assigneeType,omitempty"`
	Versions        []Version          `json:"versions,omitempty" structs:"versions,omitempty"`
	Name            string             `json:"name,omitempty" structs:"name,omitempty"`
	Roles           map[string]string  `json:"roles,omitempty" structs:"roles,omitempty"`
	AvatarUrls      AvatarUrls         `json:"avatarUrls,omitempty" structs:"avatarUrls,omitempty"`
	ProjectCategory ProjectCategory    `json:"projectCategory,omitempty" structs:"projectCategory,omitempty"`
}