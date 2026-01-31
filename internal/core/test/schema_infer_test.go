package core_test

import (
	"testing"

	"reshape/internal/core"
)

func TestBuildShapeFromRecords(t *testing.T) {
	records := []core.Record{
		{"name": "Ada", "age": 30.0, "tags": []any{"a", "b"}},
		{"name": nil, "tags": []any{}},
	}

	shape := core.BuildShapeFromRecords(records)

	field := fieldByPath(shape, "name")
	if field.Path == "" {
		t.Fatalf("expected name field")
	}
	if !field.Nullable {
		t.Fatalf("expected name to be nullable")
	}
	if field.Type != core.LogicalTypeString {
		t.Fatalf("expected name type string, got %s", field.Type)
	}

	ageField := fieldByPath(shape, "age")
	if ageField.Path == "" {
		t.Fatalf("expected age field")
	}
	if !ageField.Nullable {
		t.Fatalf("expected age to be nullable")
	}
	if ageField.Type != core.LogicalTypeNumber {
		t.Fatalf("expected age type number, got %s", ageField.Type)
	}

	tagsField := fieldByPath(shape, "tags")
	if tagsField.Path == "" {
		t.Fatalf("expected tags field")
	}
	if !tagsField.Repeated {
		t.Fatalf("expected tags to be repeated")
	}
	if tagsField.Type != core.LogicalTypeArray {
		t.Fatalf("expected tags type array, got %s", tagsField.Type)
	}
}

func fieldByPath(shape core.DataShape, path string) core.FieldDefinition {
	for _, field := range shape.Fields {
		if field.Path == path {
			return field
		}
	}
	return core.FieldDefinition{}
}
