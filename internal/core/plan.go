package core

// LossReason describes why a lossy decision is needed.
type LossReason string

const (
	LossReasonFormatLimit LossReason = "format_limit"
	LossReasonUserRequest LossReason = "user_request"
)

// Strategy describes how a lossy decision is applied.
type Strategy string

const (
	StrategyJoinArray Strategy = "join_array"
	StrategyDropField Strategy = "drop_field"
	StrategyCoerceType Strategy = "coerce_type"
)

// LossyDecision records explicit approval for a lossy action.
type LossyDecision struct {
	FieldPath string     `json:"field_path"`
	Reason    LossReason `json:"reason"`
	Strategy  Strategy   `json:"strategy"`
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
	LossyDecisions  []LossyDecision    `json:"lossy_decisions,omitempty"`
}
