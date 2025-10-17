# crispy-broccoli

A playful single-binary Go CLI program demonstrating a tiny Bubble Tea TUI.

Overview

- On start the program prints a small broccoli ASCII-art banner.
- It then displays a menu with three choices:
  - Enter Text
  - Goof off
  - Act a Fool

Controls

- Navigate with the arrow keys or `j`/`k`.
- Press Enter to select an option.
- Press Esc or Ctrl+C to quit.

Behavior

- Enter Text: opens a text-input UI (Bubbles textinput) prompting you to "Enter some important text". After pressing Enter the program prints the entered text and briefly pauses.
- Goof off: prints "You chose to goof off!" and sleeps for 2 seconds, then exits.
- Act a Fool: prints "You chose to act a fool!" and sleeps for 2 seconds, then exits.

Usage

Module-aware (recommended):

1. Download dependencies:

	go mod tidy

2. Run the program:

	go run Main/main.go

3. Build a binary:

	go build -o bin/crispy-broccoli Main/main.go

Example session

1. Start the program:

	go run Main/main.go

2. Choose "Enter Text", type some text (e.g. "Hello world") and press Enter.
3. The program will print the text you entered and then exit.

Notes for contributors

- Source: `Main/main.go` (menu, text-input flow) and `Main/taro.go` (a secondary example file; excluded from normal builds with the `example` build tag).
- Module path: `github.com/nenxuto/crispy-broccoli` (see `go.mod`).
- To run files excluded by build tags use `-tags=example` with `go run` or `go build`.

If you want, I can add a small GitHub Actions workflow to automatically run `go build` on pushes/PRs, or expand this README with a CONTRIBUTING guide.
