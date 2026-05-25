# Golumen

[![Go Reference](https://pkg.go.dev/badge/github.com/erikw/golumen.svg)](https://pkg.go.dev/github.com/erikw/golumen)
[![Go Version](https://img.shields.io/github/go-mod/go-version/erikw/golumen)](./go.mod)
[![Release](https://img.shields.io/github/v/release/erikw/golumen)](https://github.com/erikw/golumen/releases)
[![License](https://img.shields.io/github/license/erikw/golumen)](./LICENSE)

> Shining a light on your file system.

**Golumen** is a Go CLI for recursively walking the file system and printing the paths it finds. The project is still early-stage, but already provides a simple searchable foundation with structured logging and versioned releases.

## Why Golumen?
The name is a portmanteau of the language **Go** and the Latin word **Lumen**, meaning light or an opening for light. Searching through a deeply nested file system can often feel like navigating a dark labyrinth; Golumen aims to illuminate the path to your data.


## Features
* **Recursive traversal:** Walks directories with Go's `filepath.WalkDir`.
* **Useful defaults:** Skips `.git` directories while traversing.
* **CLI logging:** Supports debug logging with `--debug`.
* **Version reporting:** Prints the build version with `--version`.

## Installation
```bash
go install github.com/erikw/golumen@latest
```

## Usage
```bash
# Walk the current directory
golumen

# Enable debug logging
golumen --debug

# Show the version
golumen --version
```
