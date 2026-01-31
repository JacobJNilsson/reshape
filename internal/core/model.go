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

// DataShape describes the paths, types, and repetition for canonical records.
type DataShape struct {
	Fields []FieldDefinition `json:"fields"`
}

// DataValues contains the canonical records.
type DataValues struct {
	Records []Record `json:"records"`
}

// Record represents a canonical record.
type Record map[string]any

// CanonicalData is the canonical representation of structured data.
type CanonicalData struct {
	Shape  DataShape  `json:"shape"`
	Values DataValues `json:"values"`
}

// Warning captures lossy or noteworthy operations.
type Warning struct {
	Path    string
	Message string
}
