package main

import (
	"debug/pe"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
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
		os.Exit(0)
	}
}

var choices = []string{"Pick a File"}

type model struct {
	cursor int
	choice string
}

type filePickerModel struct {
	filepicker   filepicker.Model
	selectedFile string
	quitting     bool
	err          error
}

type Section struct {
	Name             string
	VirtualAddress   uint32
	SizeOfRawData    uint32
	PointerToRawData uint32
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
			if m.choice == "Pick a File" {
				Coffee()
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
	s.WriteString("\nWelcome to Crispy Broccoli - Go String Parser\n\n")

	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\nPress esc to quit\n")
	s.WriteString("\nEnter to Continue\n")
	return s.String()
}

func coffeeModel() filePickerModel {
	fp := filepicker.New()
	return filePickerModel{
		filepicker:   fp,
		selectedFile: "",
		quitting:     false,
		err:          nil,
	}
}

func Coffee() {
	p := tea.NewProgram(coffeeModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m filePickerModel) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m filePickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.selectedFile = path

		// Read the file.
		fileData, err := fileReader(m.selectedFile)
		if err != nil {
			// Show read error in the view for a short time.
			m.err = err
			m.selectedFile = ""
			return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
		}
		const patternLength = 17
		var matches []int
		i := 0
		for i <= len(fileData)-patternLength {
			if fileData[i] == 0x48 &&
				fileData[i+1] == 0x8D &&
				(fileData[i+2] == 0x0D || fileData[i+2] == 0x05) &&
				(fileData[i+7] == 0xBF || fileData[i+7] == 0xBB) &&
				fileData[i+12] == 0xE8 {
				matches = append(matches, i)
				i += patternLength
			} else {
				i++
			}
		}
		if len(matches) == 0 {
			fmt.Print("No Broccoli detected.")
		} else {
			f, err := os.Open(m.selectedFile)
			if err != nil {
				panic(err)
			}
			peFile, err := pe.NewFile(f)
			if err != nil {
				panic(err)
			}
			var imageBase uint64
			if optHdr, ok := peFile.OptionalHeader.(*pe.OptionalHeader64); ok {
				imageBase = optHdr.ImageBase
			} else {
				panic("PE is not 64-bit or missing optional header")
			}
			for _, off := range matches {
				rva, err := fileOffsetToRVA(uint32(off), peFile.Sections)
				if err != nil {
					fmt.Printf("❌ Offset 0x%X: %v\n", off, err)
					continue
				}
				virtualAddress := imageBase + uint64(rva)
				disp, err := readDisplacement(fileData, off)
				if err != nil {
					fmt.Println("❌ Error reading displacement:", err)
				}
				ripAfter := virtualAddress + 7
				stringVA := ripAfter + uint64(disp)
				stringOffset, err := vaToFileOffset(stringVA, imageBase, peFile.Sections)
				if err != nil {
					fmt.Println("❌ Error converting VA to file offset:", err)
				}
				if off+12 > len(fileData) {
					fmt.Printf("❌ Offset 0x%X: string length out of bounds\n", off)
					continue
				}
				length := binary.LittleEndian.Uint32(fileData[off+8 : off+12])
				if int(stringOffset)+int(length) > len(fileData) {
					fmt.Printf("❌ Offset 0x%X: string offset out of bounds (offset=0x%X, len=%d)\n", off, stringOffset, length)
					continue
				}
				stringBytes := fileData[stringOffset : stringOffset+length]
				if len(stringBytes) == 0 {
					fmt.Printf("❌ Offset 0x%X: extracted empty string\n", off)
					continue
				}
				fmt.Printf("String: %s\n", string(stringBytes))
			}
		}
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m filePickerModel) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.selectedFile == "" {
		s.WriteString("Pick a file:")
	} else {
		s.WriteString("Selected file: " + m.filepicker.Styles.Selected.Render(m.selectedFile))
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")
	return s.String()
}

func fileReader(path string) ([]byte, error) {
	filePath := path
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return nil, err
	}
	return data, nil
}

func fileOffsetToRVA(fileOffset uint32, sections []*pe.Section) (uint32, error) {
	for _, s := range sections {
		start := s.Offset
		end := s.Offset + s.Size
		if fileOffset >= start && fileOffset < end {
			return s.VirtualAddress + (fileOffset - start), nil
		}
	}
	return 0, errors.New("file offset not in any section")
}

func readDisplacement(data []byte, offset int) (int32, error) {
	if offset+7 > len(data) {
		return 0, fmt.Errorf("offset %d out of bounds", offset)
	}
	dispBytes := data[offset+3 : offset+7]
	disp := int32(binary.LittleEndian.Uint32(dispBytes))
	return disp, nil
}

func vaToFileOffset(va uint64, imageBase uint64, sections []*pe.Section) (uint32, error) {
	rva := uint32(va - imageBase)
	for _, s := range sections {
		vaStart := s.VirtualAddress
		vaEnd := s.VirtualAddress + s.VirtualSize
		if rva >= vaStart && rva < vaEnd {
			offset := (rva - vaStart) + s.Offset
			return offset, nil
		}
	}
	return 0, fmt.Errorf("VA 0x%X not found in any section", va)
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
	fmt.Println("		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢻⣿⡆⢸⣿⣿⠀⣼⡿⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀")
	fmt.Println("		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⣿⣿⣿⣿⣿⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀")
	fmt.Println("		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠛⠋⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀")
	fmt.Println("  ____      _                   ____                          _ _ ")
	fmt.Println(" / ___|_ __(_)___ _ __  _   _  | __ ) _ __ ___   ___ ___ ___ | (_)")
	fmt.Println("| |   | '__| / __| '_ \\| | | | |  _ \\| '__/ _ \\ / __/ __/ _ \\| | |")
	fmt.Println("| |___| |  | \\__ \\ |_) | |_| | | |_) | | | (_) | (_| (_| (_) | | |")
	fmt.Println(" \\____|_|  |_|___/ .__/ \\__, | |____/|_|  \\___/ \\___\\___\\___/|_|_|")
	fmt.Println("                 |_|    |___/                                     ")
}
