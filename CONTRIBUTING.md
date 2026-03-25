# Contributing to Algoryn Fabric

Thank you for your interest in contributing to Algoryn Fabric.
This repository defines the shared contracts for the Algoryn ecosystem —
changes here affect every tool in the stack.

## Before you contribute

Fabric defines contracts, not implementations. Every change to a type,
field, or constant is a potential breaking change for Pulse, Relay,
Beacon, and any third-party integration.

Please open an issue before submitting a PR for any of the following:
- Adding or removing fields from existing types
- Renaming types or constants
- Changing field types

Small additions (new constants, new payload types for new event types)
can go directly to a PR.

## Development
```bash
git clone https://github.com/algoryn-io/fabric
cd fabric
go mod tidy
go test ./...
```

## Compatibility policy

Fabric follows semantic versioning strictly:

- **Patch** (v0.x.Y) — documentation, comments, non-breaking additions
- **Minor** (v0.X.0) — new types, new constants, new optional fields
- **Major** (vX.0.0) — breaking changes to existing contracts

While in v0.x.x, minor versions may introduce breaking changes
with a deprecation notice in the PR description and CHANGELOG.

## Submitting a PR

1. Fork the repo
2. Create a branch: `feat/your-change` or `fix/your-fix`
3. Make your changes
4. Run `go vet ./...` and `go test ./...`
5. Update `CHANGELOG.md`
6. Submit the PR with a clear description of what changes and why

## Questions

Open an issue with the `question` label.