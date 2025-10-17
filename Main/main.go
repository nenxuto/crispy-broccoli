package main

// A simple example that shows how to retrieve a value from a Bubble Tea
// program after the Bubble Tea has exited.

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	printBroccoli()
	p := tea.NewProgram(model{})

	m, err := p.Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}
	if m, ok := m.(model); ok && m.choice != "" {
		tea.Quit()
	}
}

var choices = []string{"Enter Text", "Goof off", "Act a Fool"}

type model struct {
	cursor int
	choice string
}

type textInputModel struct {
	textInput textinput.Model
	err       error
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("Crispy Broccoli")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			m.choice = choices[m.cursor]
			if m.choice == "Enter Text" {
				Taro()
			}
			if m.choice == "Goof off" {
				fmt.Println("\nYou chose to goof off!")
				time.Sleep(2.0 * time.Second)
			}
			if m.choice == "Act a Fool" {
				fmt.Println("\nYou chose to act a fool!")
				time.Sleep(2.0 * time.Second)
			}
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	s.WriteString("\nWhat would you like to do?\n\n")

	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\npress esc to quit\n")
	s.WriteString("\nArrow keys to navigate. Enter to select.\n")
	return s.String()
}

func Taro() {
	p := tea.NewProgram(taroModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func taroModel() textInputModel {
	ti := textinput.New()
	ti.Placeholder = "Enter some text"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return textInputModel{
		textInput: ti,
		err:       nil,
	}
}

type errMsg error

func (m textInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m textInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			fmt.Println("\nYou entered:", m.textInput.Value())
			time.Sleep(2.0 * time.Second)
			return m, tea.Quit
		}

	case errMsg:
		m.err = error(msg)
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m textInputModel) View() string {
	return fmt.Sprintf(
		"\nEnter some important text\n\n%s\n\n%s",
		m.textInput.View(),
		"press esc to quit",
	) + "\n"
}

func printBroccoli() {
	fmt.Println("		 ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣤⣤⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀")
	fmt.Println("		⠀⠀⠀⠀⢀⣤⣤⣤⣤⣄⣀⣾⣿⣿⣿⣿⣿⣿⣷⣀⣠⣤⣤⣤⣤⡀⠀⠀⠀⠀")
	fmt.Println("		⠀⠀⠀⠀⣿⣿⣿⣿⣿⣿⣿⣿⡿⠿⠿⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀")
	fmt.Println("		⠀⠀⠀⠀⠹⣿⣿⣿⣿⣿⣿⣿⣿⣶⣶⣄⠘⠋⣉⣀⣈⠙⢿⣿⣿⣁⡀⠀⠀⠀")
	fmt.Println("		⠀⠀⣠⣶⣿⣿⡿⣿⣿⡿⠿⠿⣿⣿⣿⣿⣷⣾⣿⣿⣿⣷⣀⣿⣿⣿⣿⣦⡀⠀")
	fmt.Println("		⠀⣼⣿⣿⣟⣁⣤⣤⣈⣴⣶⣦⡈⢻⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡷⠀")
	fmt.Println("		⠀⠈⠿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠿⠁⠀")
	fmt.Println("		⠀⠀⠀⠈⠉⠛⢻⣿⣿⣿⣿⣿⣿⠋⠀⠙⣿⣿⣿⣿⣿⡟⠛⠛⠛⠋⠁⠀⠀⠀")
	fmt.Println("		⠀⠀⠀⠀⠀⠀⠈⠛⠿⠿⠿⠋⠁⣴⣿⡆⠈⠙⠛⠛⠋⡀⠀⠀⠀⠀⠀⠀⠀⠀")
	fmt.Println("		⠀⠀⠀⠀⠀⠀⠀⠀⠰⣶⣶⡀⠀⢸⣿⡇⠀⠀⢰⣾⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀")
	fmt.Println("		⠀⠀⠀⠀⠀⠀⠀⠀⠀⢹⣿⣧⠀⢸⣿⡇⠀⢀⣿⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀")
	fmt.Println("		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢻⣿⡆⢸⣿⣿⠀⣼⡿⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀")
	fmt.Println("		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢿⣿⣿⣿⣿⣿⣿⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀")
	fmt.Println("		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⣿⣿⣿⣿⣿⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀")
	fmt.Println("		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠛⠋⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀")
	fmt.Println("  ____      _                   ____                          _ _ ")
	fmt.Println(" / ___|_ __(_)___ _ __  _   _  | __ ) _ __ ___   ___ ___ ___ | (_)")
	fmt.Println("| |   | '__| / __| '_ \\| | | | |  _ \\| '__/ _ \\ / __/ __/ _ \\| | |")
	fmt.Println("| |___| |  | \\__ \\ |_) | |_| | | |_) | | | (_) | (_| (_| (_) | | |")
	fmt.Println(" \\____|_|  |_|___/ .__/ \\__, | |____/|_|  \\___/ \\___\\___\\___/|_|_|")
	fmt.Println("                 |_|    |___/                                     ")
}
