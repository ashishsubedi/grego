package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

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

type Model struct {
	words       map[int]string
	meanings    map[int]string
	index       int
	totalWords  int
	showMeaning bool
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
	}
	return m, nil
}

func (m *Model) View() string {
	s := ""
	s += "Word: "
	s += style.Render(fmt.Sprintf("%s\n", m.words[m.index]))

	if m.showMeaning {
		s += fmt.Sprintf("\nMeaning: %s\n", m.meanings[m.index])
	}
	s += fmt.Sprintf("\n\n\n\n\n\n\n\n\n\n\n\n\n")
	// Instruction and index count at end
	s += fmt.Sprintf("\n↑,k: prev\t↓,j: next\ts: Toogle show/hide meaning\tr: Random Word")
	s += fmt.Sprintf("\nTotal Words: %d\tindex: %d", m.totalWords, m.index)
	return s
}
