# crispy-broccoli

A small single-binary Go CLI program that demonstrates a simple terminal UI using Bubble Tea and Bubbles.

What the program does

- On start it shows a small menu asking "What kind of Bubble Tea would you like to order?" with three choices:
  - Taro
  - Coffee
  - Lychee

- You navigate with the arrow keys (or `j`/`k`) and press Enter to select.

- If you choose `Taro`, the program opens a second input UI that asks:
  "What’s your favorite Pokémon?"

  You can type a short answer (text input) then press Enter or Esc to quit that input UI and return. The program then exits and prints your original menu choice.

Usage (module-aware, recommended)

1. Ensure dependencies are present (downloads modules):

	go mod tidy

2. Run the program:

	go run Main/main.go

3. Build a binary:

	go build -o bin/crispy-broccoli Main/main.go

Example session

1. Start the program:

	go run Main/main.go

2. Use the arrow keys to select "Taro" and press Enter.
3. Type a Pokémon name (e.g. Pikachu) and press Enter.
4. After the UI exits the program prints the final menu selection, for example:

	---
	You chose Taro!

Notes for contributors

- Source: `Main/main.go` (menu & UI flow) and `Main/taro.go` (example text input — excluded from normal builds by the `example` build tag).
- Module path: `github.com/nenxuto/crispy-broccoli` (see `go.mod`).
- To run examples that are excluded by build tags use `-tags=example` with `go run` or `go build`.

If you'd like I can add a small GitHub Actions workflow to run `go build` on push/PRs or expand the README with a CONTRIBUTING section.
