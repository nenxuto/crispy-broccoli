# crispy-broccoli
# crispy-broccoli

A focused command-line helper that scans a selected file for a specific binary pattern and extracts embedded strings from PE (Portable Executable) files.

At startup the program prints a decorative ASCII "broccoli" header and presents a minimal interactive menu (single item: "Pick a File"). Use the file picker to select a binary, and the program will search its bytes for a particular instruction pattern, compute the referenced string address, and print any extracted strings.

What it does

- Prints an ASCII-art header on start.
- Presents a single-menu UI: "Pick a File".
- When you pick a file the program:
  - Reads the file bytes.
  - Scans for a specific 17-byte instruction pattern (the program looks for opcodes matching a heuristic pattern used to locate embedded strings).
  - For each match it converts the file offset to an RVA/VA using the PE sections and the ImageBase, reads a displacement and length from nearby bytes, and extracts the target string bytes.
  - Prints any found strings to stdout. If none are found it prints "No Broccoli detected.".

Keybindings and navigation

- Arrow keys or `j`/`k` to move (menu has one item).
- Enter to select the highlighted menu entry.
- Esc or Ctrl+C to quit.

Typical output

- If matches are found the program will print lines like:

  String: <extracted-bytes-as-text>

- If no matches are found you'll see:

  No Broccoli detected.

Usage

1. Ensure you have Go installed (the repo uses Go modules).

2. Download dependencies and tidy the module:

  go mod tidy

3. Run the program:

  go run Main/main.go

4. Build a binary:

  go build -o bin/crispy-broccoli Main/main.go

Module

This repository uses Go modules. Module path (see `go.mod`): `github.com/nenxuto/crispy-broccoli`.

Notes for contributors

- Main logic is implemented in `Main/main.go`.
- The program uses the Bubbles filepicker component to let the user choose a file interactively.
- The PE parsing uses `debug/pe` to map file offsets to RVAs/virtual addresses; the code handles 64-bit Optional Header (`OptionalHeader64`).
- The byte-pattern search is a heuristic for locating instructions that reference embedded strings; review and adjust the pattern if you need to target different compilers/optimizations.

Safety & privacy

- The program only reads local files selected by the user and prints extracted strings — it does not transmit data over the network.

Next steps (optional)

- Add a command-line mode to process files non-interactively (e.g., `crispy-broccoli scan <path>`).
- Add unit tests for the parsing helpers (`fileOffsetToRVA`, `vaToFileOffset`, `readDisplacement`).
- Add a small CI workflow to run `go build` on pushes/PRs.

If you want any of the next steps implemented, tell me which and I will add them.
