package cli_test

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCLIInferPlanJSONToCSV(t *testing.T) {
	cmd := exec.Command("go", "run", "./cli", "--from", "json", "--to", "csv", "--infer-plan")
	cmd.Dir = filepath.Join("..", "..")
	cmd.Stdin = bytes.NewBufferString(`{"name":"Ada","age":30}`)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cli error: %v\n%s", err, string(output))
	}

	expected := "age,name\n30,Ada\n"
	if string(output) != expected {
		t.Fatalf("unexpected output\nexpected: %q\nactual: %q", expected, string(output))
	}
}

func TestCLIInspectShape(t *testing.T) {
	cmd := exec.Command("go", "run", "./cli", "--from", "json", "--inspect")
	cmd.Dir = filepath.Join("..", "..")
	cmd.Stdin = bytes.NewBufferString(`[
  {"user":{"id":1},"tags":["a","b"]},
  {"user":{"id":2}}
]`)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cli error: %v\n%s", err, string(output))
	}

	expected := `{"shape":{"fields":[{"path":"tags","type":"string","nullable":true,"repeated":true},{"path":"user","type":"object","nullable":false,"repeated":false},{"path":"user.id","type":"number","nullable":false,"repeated":false}]}}`
	if string(output) != expected {
		t.Fatalf("unexpected output\nexpected: %q\nactual: %q", expected, string(output))
	}
}
