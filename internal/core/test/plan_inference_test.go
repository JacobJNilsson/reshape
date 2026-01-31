package core_test

import (
	"testing"

	"reshape/internal/core"
)

func TestInferConversionPlanForCSV(t *testing.T) {
	data := core.CanonicalData{Values: core.DataValues{Records: []core.Record{
		{
			"user":  map[string]any{"name": "Ada"},
			"tags":  []any{"a", "b"},
			"items": []any{map[string]any{"sku": "1"}},
			"flags": []any{true, false},
		},
	}}}

	plan := core.InferConversionPlan(data, "csv")

	if !contains(plan.FlattenFields, "user") {
		t.Fatalf("expected flatten for user")
	}
	if !contains(plan.ExplodeArrays, "items") {
		t.Fatalf("expected explode for items")
	}
	if !hasJoinArray(plan.JoinArrays, "tags") {
		t.Fatalf("expected join for tags")
	}
	if !hasJoinArray(plan.JoinArrays, "flags") {
		t.Fatalf("expected join for flags")
	}
	if !hasLossyOp(plan.LossyOperations, "tags", core.LossyOperationJoinArray) {
		t.Fatalf("expected lossy op for tags")
	}
	if !hasLossyOp(plan.LossyOperations, "flags", core.LossyOperationJoinArray) {
		t.Fatalf("expected lossy op for flags")
	}
}

func contains(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}

func hasJoinArray(rules []core.JoinArrayRule, path string) bool {
	for _, rule := range rules {
		if rule.Path == path {
			return true
		}
	}
	return false
}

func hasLossyOp(ops []core.LossyOperation, path string, op core.LossyOperationType) bool {
	for _, entry := range ops {
		if entry.Path == path && entry.Operation == op {
			return true
		}
	}
	return false
}
