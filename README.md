# Golumen

[![Go Reference](https://pkg.go.dev/badge/github.com/erikw/golumen.svg)](https://pkg.go.dev/github.com/erikw/golumen)
[![Go Version](https://img.shields.io/github/go-mod/go-version/erikw/golumen)](./go.mod)
[![CI](https://github.com/erikw/golumen/actions/workflows/ci.yml/badge.svg)](https://github.com/erikw/golumen/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/erikw/golumen)](https://github.com/erikw/golumen/releases)
[![License](https://img.shields.io/github/license/erikw/golumen)](./LICENSE)

> Shining a light on your file system.

**Golumen** is a Go CLI for recursively walking the file system and printing the paths it finds. The project is still early-stage, but already provides a simple searchable foundation with structured logging and versioned releases.

## Why Golumen?
The name is a portmanteau of the language **Go** and the Latin word **Lumen**, meaning light or an opening for light. Searching through a deeply nested file system can often feel like navigating a dark labyrinth; Golumen aims to illuminate the path to your data.


## Features
* **Recursive traversal:** Walks directory trees recursively and prints matching paths.
* **Optional symlink following:** Can follow symlinked directories with `-f`/`--follow` while avoiding loops.
* **Useful defaults:** Skips `.git` directories while traversing.
* **CLI logging:** Supports debug logging with `--debug`.
* **Version reporting:** Prints the build version with `--version`.

## Installation
```bash
go install github.com/erikw/golumen@latest
```

## Usage
```bash
# Walk the current directory for Go files
golumen '*.go'

# Walk a specific path for Go files
golumen '*.go' ./cmd

# Follow symlinked directories during traversal
golumen --follow '*.go' ./cmd

# Enable debug logging
golumen --debug '*.go'

# Show the version
golumen --version
```

## Testing
Run the full test suite from the repository root with:

```bash
make test
```

This wraps the standard Go command:

```bash
go test ./...
```

To run a specific package's tests directly:

```bash
go test ./cmd
```
