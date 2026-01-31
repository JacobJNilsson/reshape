package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"reshape/internal/core"
	"reshape/internal/formats"
)

func main() {
	fromFlag := flag.String("from", "", "input format: json or csv")
	toFlag := flag.String("to", "", "output format: json or csv")
	planPath := flag.String("plan", "", "path to conversion plan JSON")
	inferPlan := flag.Bool("infer-plan", false, "infer a conversion plan")
	flag.Parse()

	if *fromFlag == "" || *toFlag == "" {
		exitWithError(errors.New("--from and --to are required"))
	}
	if *inferPlan && *planPath != "" {
		exitWithError(errors.New("--plan and --infer-plan cannot be used together"))
	}

	inputPath := ""
	args := flag.Args()
	if len(args) > 1 {
		exitWithError(errors.New("only one input path argument is supported"))
	}
	if len(args) == 1 {
		inputPath = args[0]
	}

	inputBytes, err := readInput(inputPath)
	if err != nil {
		exitWithError(err)
	}

	inputData, err := parseInput(*fromFlag, inputBytes)
	if err != nil {
		exitWithError(err)
	}

	plan := core.ConversionPlan{}
	if *planPath != "" {
		planBytes, err := os.ReadFile(*planPath)
		if err != nil {
			exitWithError(err)
		}
		if err := json.Unmarshal(planBytes, &plan); err != nil {
			exitWithError(err)
		}
	}
	if *inferPlan {
		plan = core.InferConversionPlan(inputData, *toFlag)
	}

	transformed, warnings, err := core.TransformData(inputData, plan)
	if err != nil {
		exitWithError(err)
	}

	outputBytes, err := renderOutput(*toFlag, transformed)
	if err != nil {
		exitWithError(err)
	}

	if len(warnings) > 0 {
		for _, warning := range warnings {
			message := fmt.Sprintf("warning: %s", warning.Message)
			if warning.Path != "" {
				message = fmt.Sprintf("%s (path: %s)", message, warning.Path)
			}
			fmt.Fprintln(os.Stderr, message)
		}
	}

	if _, err := os.Stdout.Write(outputBytes); err != nil {
		exitWithError(err)
	}
}

func readInput(path string) ([]byte, error) {
	if path == "" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}

func parseInput(format string, input []byte) (core.CanonicalData, error) {
	switch format {
	case "json":
		return formats.ParseJSON(input)
	case "csv":
		return formats.ParseCSV(input)
	default:
		return core.CanonicalData{}, errors.New("unsupported --from format")
	}
}

func renderOutput(format string, data core.CanonicalData) ([]byte, error) {
	switch format {
	case "json":
		return formats.RenderJSON(data)
	case "csv":
		return formats.RenderCSV(data)
	default:
		return nil, errors.New("unsupported --to format")
	}
}

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, "error:", err.Error())
	os.Exit(1)
}
