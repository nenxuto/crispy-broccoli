package main

import (
	"bytes"
	"debug/macho"
	"debug/pe"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

/* =========================
   Binary Abstraction
========================= */

type Binary interface {
	ImageBase() uint64
	Sections() []BinarySection
	Close() error
}

type BinarySection struct {
	Name   string
	Addr   uint64
	Size   uint64
	Offset uint64
}

/* =========================
   PE Implementation
========================= */

type PEBinary struct {
	f  *pe.File
	os *os.File
}

func openPE(path string) (*PEBinary, error) {
	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	pf, err := pe.NewFile(fh)
	if err != nil {
		fh.Close()
		return nil, err
	}
	return &PEBinary{f: pf, os: fh}, nil
}

func (p *PEBinary) ImageBase() uint64 {
	if opt, ok := p.f.OptionalHeader.(*pe.OptionalHeader64); ok {
		return opt.ImageBase
	}
	return 0
}

func (p *PEBinary) Sections() []BinarySection {
	var out []BinarySection
	for _, s := range p.f.Sections {
		out = append(out, BinarySection{
			Name:   s.Name,
			Addr:   uint64(s.VirtualAddress),
			Size:   uint64(s.VirtualSize),
			Offset: uint64(s.Offset),
		})
	}
	return out
}

func (p *PEBinary) Close() error {
	p.f.Close()
	return p.os.Close()
}

/* =========================
   Mach-O Implementation
========================= */

type MachOBinary struct {
	f *macho.File
}

func openMachO(path string) (*MachOBinary, error) {
	f, err := macho.Open(path)
	if err != nil {
		return nil, err
	}
	return &MachOBinary{f: f}, nil
}

func (m *MachOBinary) ImageBase() uint64 {
	return 0
}

func (m *MachOBinary) Sections() []BinarySection {
	var out []BinarySection
	for _, s := range m.f.Sections {
		out = append(out, BinarySection{
			Name:   s.Name,
			Addr:   s.Addr,
			Size:   s.Size,
			Offset: uint64(s.Offset),
		})
	}
	return out
}

func (m *MachOBinary) Close() error {
	return m.f.Close()
}

/* =========================
   Binary Detection
========================= */

func detectBinary(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var magic [4]byte
	if _, err := io.ReadFull(f, magic[:]); err != nil {
		return "", err
	}

	// PE
	if magic[0] == 'M' && magic[1] == 'Z' {
		return "pe", nil
	}

	be := binary.BigEndian.Uint32(magic[:])
	le := binary.LittleEndian.Uint32(magic[:])

	switch {
	// Mach-O 32
	case be == 0xfeedface || le == 0xfeedface:
		return "macho", nil
	// Mach-O 64
	case be == 0xfeedfacf || le == 0xfeedfacf:
		return "macho", nil
	// Fat Mach-O
	case be == 0xcafebabe || le == 0xcafebabe:
		return "macho", nil
	}

	return "", errors.New("unknown binary format")
}

/* =========================
   Cross-Platform Output Path
========================= */

func outputFilePath(filename string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "windows", "darwin":
		return filepath.Join(home, "Desktop", filename), nil

	case "linux":
		desktop := filepath.Join(home, "Desktop")
		if info, err := os.Stat(desktop); err == nil && info.IsDir() {
			return filepath.Join(desktop, filename), nil
		}
		return filename, nil

	default:
		return filename, nil
	}
}

/* =========================
   UI
========================= */

type model struct{}

type filePickerModel struct {
	filepicker filepicker.Model
}

func main() {
	printBroccoli()
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("Crispy Broccoli")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		switch k.String() {
		case "enter":
			p := tea.NewProgram(filePickerModel{filepicker: filepicker.New()})
			p.Run()
			return m, tea.Quit
		case "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return "\nWelcome to Crispy Broccoli\n\n(•) Pick a File\n\nPress Enter\n"
}

func (m filePickerModel) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m filePickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		processFile(path)
		time.Sleep(3 * time.Second)
		return m, tea.Quit
	}
	return m, cmd
}

func (m filePickerModel) View() string {
	return "\nPick a file:\n\n" + m.filepicker.View()
}

/* =========================
   Core Logic
========================= */

func processFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !bytes.Contains(data, []byte("Go build ID:")) {
		fmt.Println("Not a Go binary")
		return
	}

	kind, err := detectBinary(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	var bin Binary
	if kind == "pe" {
		bin, err = openPE(path)
	} else {
		bin, err = openMachO(path)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	defer bin.Close()

	outPath, err := outputFilePath("Strings_Output.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	out, err := os.Create(outPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer out.Close()

	imageBase := bin.ImageBase()
	sections := bin.Sections()

	for i := 0; i+17 < len(data); i++ {
		if data[i] == 0x48 && data[i+12] == 0xE8 {
			va, err := fileOffsetToVA(uint64(i), imageBase, sections)
			if err != nil {
				continue
			}

			disp := int32(binary.LittleEndian.Uint32(data[i+3 : i+7]))
			strVA := va + 7 + uint64(disp)

			off, err := vaToFileOffset(strVA, imageBase, sections)
			if err != nil {
				continue
			}

			length := binary.LittleEndian.Uint32(data[i+8 : i+12])
			if int(off+uint64(length)) > len(data) {
				continue
			}

			s := data[off : off+uint64(length)]
			if isPrintable(s) {
				fmt.Fprintf(out, "VA 0x%X: %s\n", va, strings.ReplaceAll(string(s), "\n", " "))
			}
		}
	}

	fmt.Printf("Strings extracted successfully to:\n%s\n", outPath)
}

/* =========================
   Helpers
========================= */

func fileOffsetToVA(off, base uint64, secs []BinarySection) (uint64, error) {
	for _, s := range secs {
		if off >= s.Offset && off < s.Offset+s.Size {
			return base + s.Addr + (off - s.Offset), nil
		}
	}
	return 0, errors.New("offset not mapped")
}

func vaToFileOffset(va, base uint64, secs []BinarySection) (uint64, error) {
	rva := va - base
	for _, s := range secs {
		if rva >= s.Addr && rva < s.Addr+s.Size {
			return s.Offset + (rva - s.Addr), nil
		}
	}
	return 0, errors.New("va not mapped")
}

func isPrintable(b []byte) bool {
	for _, c := range b {
		if c < 32 || c > 126 {
			return false
		}
	}
	return true
}

func printBroccoli() {
	fmt.Println("  ____      _                   ____                          _ _ ")
	fmt.Println(" / ___|_ __(_)___ _ __  _   _  | __ ) _ __ ___   ___ ___ ___ | (_)")
	fmt.Println("| |   | '__| / __| '_ \\| | | | |  _ \\| '__/ _ \\ / __/ __/ _ \\| | |")
	fmt.Println("| |___| |  | \\__ \\ |_) | |_| | | |_) | | | (_) | (_| (_| (_) | | |")
	fmt.Println(" \\____|_|  |_|___/ .__/ \\__, | |____/|_|  \\___/ \\___\\___\\___/|_|_|")
	fmt.Println("                 |_|    |___/                                     ")
}
