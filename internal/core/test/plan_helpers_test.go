package core_test

import (
	"strings"
	"testing"

	"reshape/internal/core"
)

func TestValidateLossyDecisionsRequiresJoinArray(t *testing.T) {
	plan := core.ConversionPlan{
		JoinArrays: []core.JoinArrayRule{{Path: "tags", Delimiter: ","}},
	}

	err := core.ValidateLossyDecisions(plan)
	if err == nil {
		t.Fatalf("expected error for missing lossy operation")
	}
	if !strings.Contains(err.Error(), "lossy_decisions") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateLossyDecisionsRequiresDropAndCoerce(t *testing.T) {
	plan := core.ConversionPlan{
		DropFields: []string{"secret"},
		TypeCoercions: []core.TypeCoercionRule{
			{Path: "total", TargetType: core.LogicalTypeNumber},
		},
		LossyDecisions: []core.LossyDecision{
			{FieldPath: "secret", Strategy: core.StrategyDropField, Reason: core.LossReasonUserRequest},
		},
	}

	err := core.ValidateLossyDecisions(plan)
	if err == nil {
		t.Fatalf("expected error for missing coercion lossy op")
	}
	if !strings.Contains(err.Error(), "type_coercions") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateLossyDecisionsAcceptsAll(t *testing.T) {
	plan := core.ConversionPlan{
		JoinArrays: []core.JoinArrayRule{{Path: "tags", Delimiter: ","}},
		DropFields: []string{"secret"},
		TypeCoercions: []core.TypeCoercionRule{
			{Path: "total", TargetType: core.LogicalTypeNumber},
		},
		LossyDecisions: []core.LossyDecision{
			{FieldPath: "tags", Strategy: core.StrategyJoinArray, Reason: core.LossReasonFormatLimit},
			{FieldPath: "secret", Strategy: core.StrategyDropField, Reason: core.LossReasonUserRequest},
			{FieldPath: "total", Strategy: core.StrategyCoerceType, Reason: core.LossReasonUserRequest},
		},
	}

	if err := core.ValidateLossyDecisions(plan); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
