package openproject

// OPGenericDescription is an structure widely used in several OpenProject API objects
type OPGenericDescription struct {
	Format string `json:"format,omitempty" structs:"format,omitempty"`
	Raw    string `json:"raw,omitempty" structs:"raw,omitempty"`
	HTML   string `json:"html,omitempty" structs:"html,omitempty"`
}

// OPGenericLink is a structure widely used in several OpenProject API objects
type OPGenericLink struct {
	Href   string `json:"href" structs:"href"`
	Title  string `json:"title,omitempty" structs:"title,omitempty"`
	Method string `json:"method,omitempty" structs:"method,omitempty"`
	Type   string `json:"type,omitempty" structs:"type,omitempty"`
}

// PaginationParam is a structure widely used in several OpenProject API objects
type PaginationParam struct {
	Total    int `json:"total" structs:"total"`
	Count    int `json:"count" structs:"count"`
	PageSize int `json:"pageSize" structs:"pageSize"`
	Offset   int `json:"offset" structs:"offset"`
}

// Schemas is the object representing OpenProject schemas.
type Schemas struct {
	Type     string `json:"_type,omitempty" structs:"_type,omitempty"`
	Total    int    `json:"total,omitempty" structs:"total,omitempty"`
	Count    int    `json:"count,omitempty" structs:"count,omitempty"`
	Embedded struct {
		Elements []SchemaEmbeddedElement `json:"elements,omitempty" structs:"elements,omitempty"`
	} `json:"_embedded,omitempty" structs:"_embedded,omitempty"`
	Links struct {
		Self *OPGenericLink `json:"self,omitempty" structs:"self,omitempty"`
	}
}

type SchemaEmbeddedElement struct {
	AttributeGroups []struct {
		Name      string   `json:"name,omitempty" structs:"name,omitempty"`
		Type      string   `json:"_type,omitempty" structs:"_type,omitempty"`
		Attribute []string `json:"attributes,omitempty" structs:"attributes,omitempty"`
	} `json:"_attributeGroups,omitempty" structs:"_attributeGroups,omitempty"`
	Type                 string                `json:"_type,omitempty" structs:"_type,omitempty"`
	Dependencies         []interface{}         `json:"_dependencies,omitempty" structs:"_dependencies,omitempty"`
	LockVersion          *SchemaEmbeddedOption `json:"lockVersion,omitempty" structs:"lockVersion,omitempty"`
	ID                   *SchemaEmbeddedOption `json:"id,omitempty" structs:"id,omitempty"`
	Subject              *SchemaEmbeddedOption `json:"subject,omitempty" structs:"subject,omitempty"`
	Description          *SchemaEmbeddedOption `json:"description,omitempty" structs:"description,omitempty"`
	ScheduleManually     *SchemaEmbeddedOption `json:"scheduleManually,omitempty" structs:"scheduleManually,omitempty"`
	StartDate            *SchemaEmbeddedOption `json:"startDate,omitempty" structs:"startDate,omitempty"`
	DueDate              *SchemaEmbeddedOption `json:"dueDate,omitempty" structs:"dueDate,omitempty"`
	DerivedStartDate     *SchemaEmbeddedOption `json:"derivedStartDate,omitempty" structs:"derivedStartDate,omitempty"`
	DerivedDueDate       *SchemaEmbeddedOption `json:"derivedDueDate,omitempty" structs:"derivedDueDate,omitempty"`
	EstimatedTime        *SchemaEmbeddedOption `json:"estimatedTime,omitempty" structs:"estimatedTime,omitempty"`
	DerivedEstimatedTime *SchemaEmbeddedOption `json:"derivedEstimatedTime,omitempty" structs:"derivedEstimatedTime,omitempty"`
	SpentTime            *SchemaEmbeddedOption `json:"spentTime,omitempty" structs:"spentTime,omitempty"`
	PercentageDone       *SchemaEmbeddedOption `json:"percentageDone,omitempty" structs:"percentageDone,omitempty"`
	CreatedAt            *SchemaEmbeddedOption `json:"createdAt,omitempty" structs:"createdAt,omitempty"`
	UpdatedAt            *SchemaEmbeddedOption `json:"updatedAt,omitempty" structs:"updatedAt,omitempty"`
	Author               *SchemaEmbeddedOption `json:"author,omitempty" structs:"author,omitempty"`
	Project              *SchemaEmbeddedOption `json:"project,omitempty" structs:"project,omitempty"`
	Parent               *SchemaEmbeddedOption `json:"parent,omitempty" structs:"parent,omitempty"`
	Assignee             *SchemaEmbeddedOption `json:"assignee,omitempty" structs:"assignee,omitempty"`
	Responsible          *SchemaEmbeddedOption `json:"responsible,omitempty" structs:"responsible,omitempty"`
	Type0                *SchemaEmbeddedOption `json:"type,omitempty" structs:"type,omitempty"`
	Status               *SchemaEmbeddedOption `json:"status,omitempty" structs:"status,omitempty"`
	Category             *SchemaEmbeddedOption `json:"category,omitempty" structs:"category,omitempty"`
	Version              *SchemaEmbeddedOption `json:"version,omitempty" structs:"version,omitempty"`
	Priority             *SchemaEmbeddedOption `json:"priority,omitempty" structs:"priority,omitempty"`
	OverallCosts         *SchemaEmbeddedOption `json:"overallCosts,omitempty" structs:"overallCosts,omitempty"`
	LaborCosts           *SchemaEmbeddedOption `json:"laborCosts,omitempty" structs:"laborCosts,omitempty"`
	MaterialCosts        *SchemaEmbeddedOption `json:"materialCosts,omitempty" structs:"materialCosts,omitempty"`
	CostsByType          *SchemaEmbeddedOption `json:"costsByType,omitempty" structs:"costsByType,omitempty"`
	Position             *SchemaEmbeddedOption `json:"position,omitempty" structs:"position,omitempty"`
	StoryPoints          *SchemaEmbeddedOption `json:"storyPoints,omitempty" structs:"storyPoints,omitempty"`
	RemainingTime        *SchemaEmbeddedOption `json:"remainingTime,omitempty" structs:"remainingTime,omitempty"`
	CustomField29        *SchemaEmbeddedOption `json:"customField29,omitempty" structs:"customField29,omitempty"`
	Links                struct {
		Self *OPGenericLink `json:"self,omitempty" structs:"self,omitempty"`
	} `json:"_links,omitempty" structs:"_links,omitempty"`
}

type SchemaEmbeddedOption struct {
	Type       string `json:"type,omitempty" structs:"type,omitempty"`
	Name       string `json:"name,omitempty" structs:"name,omitempty"`
	Required   bool   `json:"required,omitempty" structs:"required,omitempty"`
	HasDefault bool   `json:"hasDefault,omitempty" structs:"hasDefault,omitempty"`
	Writable   bool   `json:"writable,omitempty" structs:"writable,omitempty"`
	Options    struct {
	} `json:"options,omitempty" structs:"options,omitempty"`
	AttributeGroup string `json:"attributeGroup,omitempty" structs:"attributeGroup,omitempty"`
	Location       string `json:"location,omitempty" structs:"location,omitempty"`
	Links          struct {
	} `json:"_links,omitempty" structs:"_links,omitempty"`
}
