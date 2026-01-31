package core_test

import (
	"reflect"
	"testing"

	"reshape/internal/core"
)

func TestTransformNestedJSONToCSVWithPlan(t *testing.T) {
	inputData := core.CanonicalData{Records: []core.Record{
		{
			"user":    map[string]any{"email": "a@example.com", "name": "Ada"},
			"tags":    []any{"alpha", "beta"},
			"metrics": map[string]any{"scores": []any{1.0, 2.0}, "active": true},
		},
	}}

	plan := core.ConversionPlan{
		FlattenFields: []string{"metrics", "user"},
		JoinArrays: []core.JoinArrayRule{
			{Path: "metrics.scores", Delimiter: ";"},
			{Path: "tags", Delimiter: ";"},
		},
		LossyOperations: []core.LossyOperation{
			{Path: "metrics.scores", Operation: core.LossyOperationJoinArray, Reason: "CSV requires scalars"},
			{Path: "tags", Operation: core.LossyOperationJoinArray, Reason: "CSV requires scalars"},
		},
	}

	transformed, _, err := core.TransformData(inputData, plan)
	if err != nil {
		t.Fatalf("transform: %v", err)
	}

	if len(transformed.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(transformed.Records))
	}
	record := transformed.Records[0]
	if record["metrics.active"] != true {
		t.Fatalf("expected metrics.active true, got %v", record["metrics.active"])
	}
	if record["metrics.scores"] != "1;2" {
		t.Fatalf("expected metrics.scores joined, got %v", record["metrics.scores"])
	}
	if record["tags"] != "alpha;beta" {
		t.Fatalf("expected tags joined, got %v", record["tags"])
	}
	if record["user.email"] != "a@example.com" || record["user.name"] != "Ada" {
		t.Fatalf("unexpected user fields: %v", record)
	}
}

func TestTransformDeterministic(t *testing.T) {
	inputData := core.CanonicalData{Records: []core.Record{
		{"name": "Ada", "age": 30.0},
		{"name": "Linus", "age": 55.0},
	}}

	plan := core.ConversionPlan{
		TypeCoercions: []core.TypeCoercionRule{
			{Path: "age", TargetType: core.LogicalTypeString},
		},
		LossyOperations: []core.LossyOperation{
			{Path: "age", Operation: core.LossyOperationCoerceType, Reason: "string output"},
		},
	}

	first, _, err := core.TransformData(inputData, plan)
	if err != nil {
		t.Fatalf("transform first: %v", err)
	}
	second, _, err := core.TransformData(inputData, plan)
	if err != nil {
		t.Fatalf("transform second: %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("transform not deterministic\nfirst: %+v\nsecond: %+v", first, second)
	}
}
