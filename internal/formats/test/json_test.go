package formats_test

import (
	"encoding/json"
	"strings"
	"testing"

	"reshape/internal/core"
	"reshape/internal/formats"
)

func TestParseJSONObject(t *testing.T) {
	input := []byte(`{"name":"Ada","age":30}`)

	data, err := formats.ParseJSON(input)
	if err != nil {
		t.Fatalf("parse json: %v", err)
	}

	if len(data.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(data.Records))
	}
	if data.Records[0]["name"] != "Ada" {
		t.Fatalf("expected name Ada, got %v", data.Records[0]["name"])
	}
}

func TestParseJSONRejectsNonObjectArray(t *testing.T) {
	_, err := formats.ParseJSON([]byte(`["a","b"]`))
	if err == nil {
		t.Fatalf("expected error for non-object array")
	}
	if !strings.Contains(err.Error(), "non-object") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderJSONShapes(t *testing.T) {
	empty, err := formats.RenderJSON(core.CanonicalData{})
	if err != nil {
		t.Fatalf("render empty: %v", err)
	}
	if string(empty) != "[]" {
		t.Fatalf("expected empty array JSON, got %s", string(empty))
	}

	data, err := formats.ParseJSON([]byte(`{"name":"Ada"}`))
	if err != nil {
		t.Fatalf("parse json: %v", err)
	}
	one, err := formats.RenderJSON(data)
	if err != nil {
		t.Fatalf("render json: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(one, &decoded); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}
	if decoded["name"] != "Ada" {
		t.Fatalf("unexpected output: %v", decoded)
	}
}
