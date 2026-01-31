package core_test

import (
	"strings"
	"testing"

	"reshape/internal/core"
	"reshape/internal/formats"
)

func TestJSONToCSVHostile(t *testing.T) {
	input := []byte(`[
  {"user":{"id":1,"name":"Ada"},"meta":{"active":true},"tags":["a,b","c"],"scores":[1,2,3],"notes":"first"},
  {"user":{"id":2,"name":"Linus"},"meta":{"active":false},"tags":["solo"],"scores":[4]}
]`)

	data, err := formats.ParseJSON(input)
	if err != nil {
		t.Fatalf("parse json: %v", err)
	}

	plan := core.InferConversionPlan(data, "csv")
	transformed, warnings, err := core.TransformData(data, plan)
	if err != nil {
		t.Fatalf("transform: %v", err)
	}

	output, err := formats.RenderCSV(transformed)
	if err != nil {
		t.Fatalf("render csv: %v", err)
	}

	expectedCSV := strings.Join([]string{
		"meta.active,notes,scores,tags,user.id,user.name",
		"true,first,\"1,2,3\",\"a,b,c\",1,Ada",
		"false,,4,solo,2,Linus",
		"",
	}, "\n")
	if string(output) != expectedCSV {
		t.Fatalf("unexpected csv output\nexpected: %q\nactual: %q", expectedCSV, string(output))
	}

	expectedWarnings := []core.Warning{
		{Code: core.WarningCodeJoinArray, Path: "scores", Message: core.WarningMessage(core.WarningCodeJoinArray)},
		{Code: core.WarningCodeJoinArray, Path: "tags", Message: core.WarningMessage(core.WarningCodeJoinArray)},
	}
	if len(warnings) != len(expectedWarnings) {
		t.Fatalf("unexpected warning count\nexpected: %d\nactual: %d", len(expectedWarnings), len(warnings))
	}
	for index, warning := range warnings {
		expected := expectedWarnings[index]
		if warning.Code != expected.Code || warning.Path != expected.Path {
			t.Fatalf("unexpected warning at %d\nexpected: %#v\nactual: %#v", index, expected, warning)
		}
	}
}
