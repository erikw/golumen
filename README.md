# Golumen

> Shining a light on your file system.

**Golumen** is a minimalist, high-performance CLI tool for finding files and directories. It is designed to be a modern, Go-based alternative to `find`, focusing on speed through concurrency and clarity through intuitive output.

## Why Golumen?
The name is a portmanteau of the language **Go** and the Latin word **Lumen**, meaning light or an opening for light. Searching through a deeply nested file system can often feel like navigating a dark labyrinth; Golumen aims to illuminate the path to your data.


## Features
* **Concurrent Execution:** Utilizes Go routines to traverse directory trees in parallel.
* **Intuitive Syntax:** Simple command-line arguments that favor human readability.
* **Smart Defaults:** Automatically respects `.gitignore` and skips hidden directories.
* **Metric-Driven:** Optimized using `filepath.WalkDir` for minimal system overhead.

## Usage
```bash
# Basic search
golumen "filename"

# Search within a specific directory
golumen --path ./src "main.go"
```
