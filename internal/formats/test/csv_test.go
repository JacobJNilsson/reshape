package formats_test

import (
	"strings"
	"testing"

	"reshape/internal/core"
	"reshape/internal/formats"
)

func TestParseCSVValid(t *testing.T) {
	input := []byte("name,age\nAda,30\nLinus,\n")

	data, err := formats.ParseCSV(input)
	if err != nil {
		t.Fatalf("parse csv: %v", err)
	}
	if len(data.Records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(data.Records))
	}
	if data.Records[1]["age"] != nil {
		t.Fatalf("expected nil age for empty column, got %v", data.Records[1]["age"])
	}
}

func TestParseCSVRejectsEmpty(t *testing.T) {
	_, err := formats.ParseCSV([]byte(""))
	if err == nil {
		t.Fatalf("expected error for empty csv")
	}
}

func TestParseCSVRejectsMismatchedRow(t *testing.T) {
	_, err := formats.ParseCSV([]byte("a,b\n1\n"))
	if err == nil {
		t.Fatalf("expected error for mismatched columns")
	}
	if !strings.Contains(err.Error(), "column count") && !strings.Contains(err.Error(), "wrong number of fields") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderCSVUsesSchemaPaths(t *testing.T) {
	data := core.CanonicalData{
		Schema: core.CanonicalSchema{Fields: []core.FieldDefinition{
			{Path: "user.name", Type: core.LogicalTypeString},
		}},
		Records: []core.Record{
			{"user": map[string]any{"name": "Ada"}},
		},
	}

	output, err := formats.RenderCSV(data)
	if err != nil {
		t.Fatalf("render csv: %v", err)
	}

	if string(output) != "user.name\nAda\n" {
		t.Fatalf("unexpected csv output: %s", string(output))
	}
}

func TestRenderCSVRejectsObjectValues(t *testing.T) {
	data := core.CanonicalData{
		Schema: core.CanonicalSchema{Fields: []core.FieldDefinition{
			{Path: "user", Type: core.LogicalTypeObject},
		}},
		Records: []core.Record{
			{"user": map[string]any{"name": "Ada"}},
		},
	}

	_, err := formats.RenderCSV(data)
	if err == nil {
		t.Fatalf("expected error for object value")
	}
	if !strings.Contains(err.Error(), "requires scalar") {
		t.Fatalf("unexpected error: %v", err)
	}
}
