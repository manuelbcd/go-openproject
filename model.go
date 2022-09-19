package openproject

// OPGenericDescription is an structure widely used in several OpenProject API objects
type OPGenericDescription struct {
	Format string `json:"format,omitempty" structs:"format,omitempty"`
	Raw    string `json:"raw,omitempty" structs:"raw,omitempty"`
	HTML   string `json:"html,omitempty" structs:"html,omitempty"`
}

// OPGenericLink is an structure widely used in several OpenProject API objects
type OPGenericLink struct {
	Href   string `json:"href"`
	Title  string `json:"title,omitempty"`
	Method string `json:"method,omitempty"`
	Type   string `json:"type,omitempty"`
}

type PaginationParam struct {
	Total    int `json:"total" structs:"total"`
	Count    int `json:"count" structs:"count"`
	PageSize int `json:"pageSize" structs:"pageSize"`
	Offset   int `json:"offset" structs:"offset"`
}
