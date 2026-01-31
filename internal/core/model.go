package core

// LogicalType describes the canonical field type.
type LogicalType string

const (
	LogicalTypeString  LogicalType = "string"
	LogicalTypeNumber  LogicalType = "number"
	LogicalTypeBoolean LogicalType = "boolean"
	LogicalTypeObject  LogicalType = "object"
	LogicalTypeArray   LogicalType = "array"
)

// FieldConstraints describe optional validation constraints.
type FieldConstraints struct {
	Minimum *float64 `json:"minimum,omitempty"`
	Maximum *float64 `json:"maximum,omitempty"`
	Enum    []string `json:"enum,omitempty"`
	Pattern string   `json:"pattern,omitempty"`
}

// FieldDefinition describes a canonical field.
type FieldDefinition struct {
	Path        string           `json:"path"`
	Type        LogicalType      `json:"type"`
	Nullable    bool             `json:"nullable"`
	Repeated    bool             `json:"repeated"`
	Constraints FieldConstraints `json:"constraints,omitempty"`
}

// CanonicalSchema describes the schema for canonical records.
type CanonicalSchema struct {
	Fields []FieldDefinition `json:"fields"`
}

// Record represents a canonical record.
type Record map[string]any

// CanonicalData is the canonical representation of structured data.
type CanonicalData struct {
	Schema  CanonicalSchema `json:"schema"`
	Records []Record        `json:"records"`
}

// Warning captures lossy or noteworthy operations.
type Warning struct {
	Path    string
	Message string
}
