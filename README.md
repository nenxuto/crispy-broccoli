# 🥦Crispy Broccoli🥦

A command line tool built with [Bubbletea](https://github.com/charmbracelet/bubbletea/tree/main) that scans a selected file and extracts embedded strings from Golang PE files.

What it does

- When you pick a file the program:
  - Reads the file bytes.
  - Scans for a specific 17-byte instruction pattern (the program looks for opcodes matching a heuristic pattern used to locate embedded strings).
  - For each match it converts the file offset to an RVA/VA using the PE sections and the ImageBase, reads a displacement and length from nearby bytes, and extracts the target string bytes.
  - Prints any found strings to a file in the working directory.

Keybindings and navigation

- Enter to select the highlighted menu entry.
- Use the arrow keys to select a file, left and right to enter/exit a directory.
- Esc or Ctrl+C to quit.

Usage

1. Ensure you have Go installed.
2. Download dependencies and tidy the module:
  go mod tidy
3. Run the program:
  go run Main/main.go
4. Build a binary:
  go build -o bin/crispy-broccoli Main/main.go

Module
This repository uses Go modules. Module path (see `go.mod`): `github.com/nenxuto/crispy-broccoli`.

Notes for contributors
- The program uses the Bubbles filepicker component to let the user choose a file interactively.
- The PE parsing uses `debug/pe` to map file offsets to RVAs/virtual addresses; the code handles 64-bit Optional Header (`OptionalHeader64`).
- The byte-pattern search is a heuristic for locating instructions that reference embedded strings; review and adjust the pattern if you need to target different compilers/optimizations.
