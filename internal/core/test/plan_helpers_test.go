package core_test

import (
	"strings"
	"testing"

	"reshape/internal/core"
)

func TestValidateLossyOperationsRequiresJoinArray(t *testing.T) {
	plan := core.ConversionPlan{
		JoinArrays: []core.JoinArrayRule{{Path: "tags", Delimiter: ","}},
	}

	err := core.ValidateLossyOperations(plan)
	if err == nil {
		t.Fatalf("expected error for missing lossy operation")
	}
	if !strings.Contains(err.Error(), "lossy_operations") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateLossyOperationsRequiresDropAndCoerce(t *testing.T) {
	plan := core.ConversionPlan{
		DropFields: []string{"secret"},
		TypeCoercions: []core.TypeCoercionRule{
			{Path: "total", TargetType: core.LogicalTypeNumber},
		},
		LossyOperations: []core.LossyOperation{
			{Path: "secret", Operation: core.LossyOperationDropField, Reason: "remove"},
		},
	}

	err := core.ValidateLossyOperations(plan)
	if err == nil {
		t.Fatalf("expected error for missing coercion lossy op")
	}
	if !strings.Contains(err.Error(), "type_coercions") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateLossyOperationsAcceptsAll(t *testing.T) {
	plan := core.ConversionPlan{
		JoinArrays: []core.JoinArrayRule{{Path: "tags", Delimiter: ","}},
		DropFields: []string{"secret"},
		TypeCoercions: []core.TypeCoercionRule{
			{Path: "total", TargetType: core.LogicalTypeNumber},
		},
		LossyOperations: []core.LossyOperation{
			{Path: "tags", Operation: core.LossyOperationJoinArray, Reason: "csv"},
			{Path: "secret", Operation: core.LossyOperationDropField, Reason: "remove"},
			{Path: "total", Operation: core.LossyOperationCoerceType, Reason: "normalize"},
		},
	}

	if err := core.ValidateLossyOperations(plan); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
