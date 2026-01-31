package core_test

import (
	"strings"
	"testing"

	"reshape/internal/core"
)

func TestValueAtPathPrefersLiteralKey(t *testing.T) {
	record := core.Record{
		"a.b": "literal",
		"a":   map[string]any{"b": "nested"},
	}

	value, exists, err := core.ValueAtPath(record, "a.b")
	if err != nil {
		t.Fatalf("value at path: %v", err)
	}
	if !exists {
		t.Fatalf("expected value at path")
	}
	if value != "literal" {
		t.Fatalf("expected literal value, got %v", value)
	}
}

func TestTransformFlattensNestedObject(t *testing.T) {
	input := core.CanonicalData{Records: []core.Record{
		{"user": map[string]any{"name": "Ada"}},
	}}
	plan := core.ConversionPlan{FlattenFields: []string{"user"}}

	output, _, err := core.TransformData(input, plan)
	if err != nil {
		t.Fatalf("transform: %v", err)
	}
	if _, ok := output.Records[0]["user"]; ok {
		t.Fatalf("expected user key to be removed")
	}
	if output.Records[0]["user.name"] != "Ada" {
		t.Fatalf("expected flattened value, got %v", output.Records[0]["user.name"])
	}
}

func TestTransformRejectsFlattenNonObject(t *testing.T) {
	input := core.CanonicalData{Records: []core.Record{
		{"user": "Ada"},
	}}
	plan := core.ConversionPlan{FlattenFields: []string{"user"}}

	_, _, err := core.TransformData(input, plan)
	if err == nil {
		t.Fatalf("expected error for non-object flatten")
	}
	if !strings.Contains(err.Error(), "not an object") {
		t.Fatalf("unexpected error: %v", err)
	}
}
