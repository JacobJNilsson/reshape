# Contributing (AI Agents)

This file describes how AI agents should reason when modifying reshape.

## How to think before implementing

Before writing code, ask:

1. Is this a decode, plan, transform, or encode concern?
2. Does this introduce an implicit or hidden decision?
3. Can this decision be represented explicitly in a plan?
4. Will this be overridable later?
5. Does this violate determinism?

If unsure, prefer explicit types over flags.

## Implementation priorities

When multiple approaches are possible, prefer:
- Explicit structs over booleans
- Named strategies over implicit behavior
- Deterministic ordering over convenience
- Plan-driven behavior over inference shortcuts

## Testing expectations

All meaningful behavior must be covered by tests that:
- Show full input data
- Assert full output data
- Assert warnings and lossy decisions
- Do not rely on format-specific shortcuts

Tests should fail if:
- Lossy behavior does not emit warnings
- Transform logic bypasses the plan
- Output ordering becomes unstable

## What not to optimize for

Do NOT optimize for:
- Minimal lines of code
- Clever abstractions
- Maximum format support early

Correct structure is more important than feature count.
