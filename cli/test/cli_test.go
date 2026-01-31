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
