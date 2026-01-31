package core

// LossyOperationType describes a lossy conversion.
type LossyOperationType string

const (
	LossyOperationJoinArray  LossyOperationType = "join_array"
	LossyOperationDropField  LossyOperationType = "drop_field"
	LossyOperationCoerceType LossyOperationType = "coerce_type"
)

// LossyOperation records explicit approval for a lossy action.
type LossyOperation struct {
	Path      string             `json:"path"`
	Operation LossyOperationType `json:"operation"`
	Reason    string             `json:"reason,omitempty"`
}

// JoinArrayRule defines how to join arrays.
type JoinArrayRule struct {
	Path      string `json:"path"`
	Delimiter string `json:"delimiter"`
}

// TypeCoercionRule defines type coercion for a field.
type TypeCoercionRule struct {
	Path       string      `json:"path"`
	TargetType LogicalType `json:"target_type"`
}

// DefaultValueRule defines a default value for a field.
type DefaultValueRule struct {
	Path  string `json:"path"`
	Value any    `json:"value"`
}

// ConversionPlan defines explicit transformation decisions.
type ConversionPlan struct {
	FlattenFields   []string           `json:"flatten_fields,omitempty"`
	JoinArrays      []JoinArrayRule    `json:"join_arrays,omitempty"`
	ExplodeArrays   []string           `json:"explode_arrays,omitempty"`
	TypeCoercions   []TypeCoercionRule `json:"type_coercions,omitempty"`
	DefaultValues   []DefaultValueRule `json:"default_values,omitempty"`
	DropFields      []string           `json:"drop_fields,omitempty"`
	LossyOperations []LossyOperation   `json:"lossy_operations,omitempty"`
}
