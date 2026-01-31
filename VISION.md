# reshape – Vision

## One-sentence goal

reshape is a deterministic data reshaping engine that converts data between formats by translating them through an explicit, inspectable intermediate representation, making all lossy decisions visible, testable, and overridable.

## What reshape is

reshape is:

- A data normalization and projection engine
- Format-agnostic at its core
- Explicit about assumptions and data loss
- Deterministic and test-driven
- Designed for inspection, planning, and automation

reshape converts data by:

1. Decoding input formats into an internal canonical representation
2. Applying an explicit transformation plan
3. Encoding the result into an output format

## What reshape is not

reshape is NOT:

- A “best effort” converter that hides decisions
- A format-to-format shortcut tool
- A schema-less magic transformer
- A streaming ETL system
- An implicit inference engine

Any ambiguity or loss must be surfaced explicitly.

## Design values

- Explicitness over convenience
- Determinism over cleverness
- Traceability over automation
- Internal correctness over CLI polish

If a transformation cannot be explained, it should not exist.
