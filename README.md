# Reshape

`reshape` is a deterministic data reshaping engine that converts data between formats by translating them through an explicit, inspectable intermediate representation, making all lossy decisions visible, testable, and overridable.

Supported formats:

- `json`
- `csv`

## Basic usage

Convert JSON to CSV with an inferred plan (stdin → stdout):

```bash
cat input.json | go run ./cli --from json --to csv --infer-plan
```

Convert CSV to JSON using a conversion plan file:

```bash
go run ./cli --from csv --to json --plan plan.json data.csv
```

Notes:

- `--from` and `--to` are required and must be one of the supported formats.
- Provide at most one input path; omit it to read from stdin.
- Warnings about lossy operations are printed to stderr.

## Conversion plan overview

Reshape uses a JSON conversion plan to make transformation decisions explicit. The plan can be fully authored or inferred for CSV targets.

Common plan fields:

- `flatten_fields`: paths to flatten nested objects into dotted keys.
- `explode_arrays`: paths to expand array items into multiple records.
- `join_arrays`: array join rules with a delimiter.
- `type_coercions`: coerce field types (string/number/boolean).
- `default_values`: set defaults when fields are missing.
- `drop_fields`: remove fields entirely.
- `lossy_operations`: explicit acknowledgements required for lossy actions.

Example (trimmed) plan:

```json
{
  "flatten_fields": ["customer.address"],
  "join_arrays": [{"path": "tags", "delimiter": ","}],
  "type_coercions": [{"path": "total", "target_type": "number"}],
  "lossy_operations": [
    {"path": "tags", "operation": "join_array", "reason": "CSV does not support arrays"},
    {"path": "total", "operation": "coerce_type", "reason": "normalize numbers"}
  ]
}
```

## How Reshape works

1) Parse input into a canonical model: `schema` + `records`.
2) Infer a schema from records to capture logical field types.
3) Load or infer a conversion plan, then normalize it for deterministic ordering.
4) Apply transformations in order: flatten → explode arrays → join arrays → coerce types → defaults → drop fields.
5) Rebuild the output schema and render to the target format.

Lossy transformations (joining arrays, type coercions, dropping fields) require explicit `lossy_operations` entries; otherwise the CLI returns an error. Warnings are emitted when lossy steps run.
