package main

import (
	"fmt"
	"iter"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	calculator "github.com/n9iels/dice-calculator/internal"
)

type question struct {
	text  string
	input textinput.Model
}

type model struct {
	calculator calculator.Calculator
	questions  []question
	focusIndex int
	output     iter.Seq[calculator.CalculatorOutput]
}

var (
	boldTextStyle = lipgloss.NewStyle().Bold(true)
	focusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("248"))
)

func initialModel() model {
	m := model{
		questions: []question{
			{text: "Amount of dice sides"},
			{text: "Amount of dice"},
			{text: "Minimum roll for a success"},
			{text: "Minimum roll for exploding (0 to disable)"},
			{text: "Amount of rolls to do for the calculation"},
		},
	}

	var t textinput.Model
	for i := range m.questions {
		t = textinput.New()
		t.PromptStyle = blurredStyle
		t.TextStyle = blurredStyle

		if i == 0 {
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		}

		m.questions[i].input = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			return initialModel(), nil
		case "enter", "tab", "shift+tab", "up", "down":
			s := msg.String()

			if m.output != nil {
				return m, nil
			}

			if s == "enter" && m.focusIndex == len(m.questions)-1 {
				diceSides, _ := strconv.Atoi(m.questions[0].input.Value())
				amountOfDice, _ := strconv.Atoi(m.questions[1].input.Value())
				minimumRollForSuccess, _ := strconv.Atoi(m.questions[2].input.Value())
				miniumRollToExplode, _ := strconv.Atoi(m.questions[3].input.Value())
				amountOfRolls, _ := strconv.Atoi(m.questions[4].input.Value())

				m.calculator = calculator.Calculator{
					DiceSides:             diceSides,
					AmountOfDice:          amountOfDice,
					MinimumRollForSuccess: minimumRollForSuccess,
					MinimumRollToExplode:  miniumRollToExplode,
					AmountOfRolls:         amountOfRolls,
				}

				m.questions[m.focusIndex].input.Blur()
				m.output = m.calculator.Calculate()

				return m, nil
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.questions) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.questions)
			}

			cmds := make([]tea.Cmd, len(m.questions))
			for i := range m.questions {
				if i == m.focusIndex {
					m.questions[i].input.PromptStyle = focusedStyle
					m.questions[i].input.TextStyle = focusedStyle
					cmds[i] = m.questions[i].input.Focus()
					continue
				}

				m.questions[i].input.Blur()
				m.questions[i].input.PromptStyle = blurredStyle
				m.questions[i].input.TextStyle = blurredStyle
			}

			return m, tea.Batch(cmds...)
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "backspace":
			cmd := m.updateInputs(msg)
			return m, tea.Batch(cmd)
		}
	}

	cmd := m.updateInputs(msg)
	return m, tea.Batch(cmd)

}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.questions))

	for i := range m.questions {
		m.questions[i].input, cmds[i] = m.questions[i].input.Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	var s strings.Builder
	var asciiArt = `
______ _            _____       _            _       _             
|  _  (_)          /  __ \     | |          | |     | |            
| | | |_  ___ ___  | /  \/ __ _| | ___ _   _| | __ _| |_ ___  _ __ 
| | | | |/ __/ _ \ | |    / __ | |/ __| | | | |/ __ | __/ _ \| '__|
| |/ /| | (_|  __/ | \__/\ (_| | | (__| |_| | | (_| | || (_) | |   
|___/ |_|\___\___|  \____/\__,_|_|\___|\__,_|_|\__,_|\__\___/|_|   `

	s.WriteString(asciiArt)
	s.WriteRune('\n')
	s.WriteRune('\n')

	for _, question := range m.questions {
		s.WriteString(question.text)
		s.WriteRune('\n')
		s.WriteString(question.input.View())
		s.WriteRune('\n')
		s.WriteRune('\n')
	}

	if m.output != nil {
		for o := range m.output {
			t := fmt.Sprintf("Roll with %d successes was made %d out of %d times with a probability of %.2f",
				o.AmountOfSuccess,
				o.AmountOfRolls,
				m.calculator.AmountOfRolls,
				float64(o.AmountOfRolls)/float64(m.calculator.AmountOfRolls),
			)

			s.WriteString(boldTextStyle.Render(t))
			s.WriteRune('\n')
			s.WriteRune('\n')
		}
	}

	s.WriteString(blurredStyle.Render("<enter> next • <shift+enter> previous • <r> restart • <q> quit"))

	return s.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
