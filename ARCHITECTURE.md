# Architecture

## High-level pipeline

All data transformations follow this pipeline:

Input Format
→ Decode
→ Internal Truth
→ Transformation Plan
→ Transform
→ Encode
→ Output Format

No step may be skipped.

## Core layers

- decode: format-specific parsing (JSON, CSV, etc.)
- truth: canonical internal data model
- plan: explicit description of transformations and decisions
- transform: applies a plan to internal truth
- encode: format-specific output rendering
- cli: thin interface only, no logic

## Invariants (must never be violated)

- All transformations go through the internal truth model
- Format-specific logic exists only in decode and encode layers
- All lossy operations:
  - are explicit
  - emit warnings
  - are representable in the transformation plan
- Identical input + plan must produce identical output
- Transform logic must not depend on CLI flags
- Inference produces plans, not transformed data

## Forbidden patterns

- jsonToCsv, csvToJson, or any format-to-format functions
- Implicit defaults inside transform logic
- Silent data loss
- Business logic inside the CLI layer
- Map iteration order affecting output

## Expected evolution

The architecture must support:
- Schema inspection
- Plan inference
- Plan overrides
- Future API usage without changes to core logic
