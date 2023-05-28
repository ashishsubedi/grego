package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#00FF00"))

func main() {
	err := tea.NewProgram(&Model{words: make(map[int]string), meanings: make(map[int]string)}, tea.WithAltScreen()).Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getTerminalDimensions() (int, int, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	dimensions := strings.Split(strings.TrimSpace(string(out)), " ")
	width, err := strconv.Atoi(dimensions[1])
	if err != nil {
		return 0, 0, err
	}

	height, err := strconv.Atoi(dimensions[0])
	if err != nil {
		return 0, 0, err
	}

	return width, height, nil
}

type Model struct {
	words       map[int]string
	meanings    map[int]string
	index       int
	totalWords  int
	showMeaning bool
	viewport    viewport.Model
	ready       bool
}

func (m *Model) Init() tea.Cmd {
	file, err := os.ReadFile("data/vocab-list.txt")
	check(err)

	lines := strings.Split(string(file), "\n")
	m.totalWords = len(lines)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })
	for i, line := range lines {
		words := strings.Split(line, "-")

		if len(words) != 2 {
			continue
		}

		m.words[i] = strings.Trim(words[0], " \n")
		m.meanings[i] = strings.Trim(words[1], " \n")
	}
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "down", "j":
			m.index++
			if m.index >= m.totalWords {
				m.index = m.totalWords - 1
			}
		case "up", "k":
			m.index--
			if m.index < 0 {

				m.index = 0
			}
		case "s":
			// show
			m.showMeaning = !m.showMeaning
		case "r":
			// Update index to random value
			m.index = rand.Intn(m.totalWords)
		}
	case tea.WindowSizeMsg:
		adjustViewport(m, msg)

	}

	m.viewport, _ = m.viewport.Update(msg)

	return m, nil
}

func adjustViewport(m *Model, msg tea.WindowSizeMsg) {
	footerHeight := lipgloss.Height(m.footerView())
	headerHeight := lipgloss.Height(m.topSection())
	verticalMarginHeight := footerHeight + headerHeight

	if !m.ready {

		m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
		m.viewport.YPosition = headerHeight
		m.viewport.HighPerformanceRendering = false

		m.ready = true

		m.viewport.YPosition = headerHeight + 1
	} else {
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - verticalMarginHeight
	}
}

func (m *Model) footerView() string {
	instructions := "\n↑,k: prev\t↓,j: next\ts: Toggle show/hide meaning\tr: Random Word"
	stats := fmt.Sprintf("\nTotal Words: %d\tindex: %d", m.totalWords, m.index)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		instructions,
		stats,
	)
}
func (m *Model) topSection() string {
	wordSection := fmt.Sprintf("Word: %s", style.Render(m.words[m.index]))
	meaningSection := "\n"
	if m.showMeaning {
		meaningSection = fmt.Sprintf("\nMeaning: %s", m.meanings[m.index])
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		wordSection,
		meaningSection,
	)
}
func (m *Model) View() string {

	// bottomSection := fmt.Sprintf("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.topSection(),
		m.viewport.View(),
		m.footerView(),
	)
}
