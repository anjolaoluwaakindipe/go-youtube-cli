package tui

import (
	"fmt"

	"github.com/anjolaoluwaakindipe/fyne-youtube/tui/state"
	videodownload "github.com/anjolaoluwaakindipe/fyne-youtube/videodownload"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// styles
var headerStyle = lipgloss.NewStyle().Bold(true).
	Italic(true).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("228")).
	BorderBackground(lipgloss.Color("63")).Render

// model to store application state
type startingUIModel struct {
	choices  []DownloadOption
	cursor   int
	selected int
}

// Download option struct
type DownloadOption struct {
	text      string
	videoType videodownload.DownloadType
}

// constructor to initilaize model. NOTE: This could also have been a variable
func InitialStartingUIModel() startingUIModel {
	// create option text for type of download a user wants to execute
	options := []DownloadOption{{text: "Download a single video", videoType: videodownload.SingleVideo}, {text: "Download a playlist", videoType: videodownload.PlayList}}
	return startingUIModel{choices: options, selected: 0}
}

// Init method for the model that must return a bumble tea cmd( initial I/O functionality can be put inside)
func (m startingUIModel) Init() tea.Cmd {
	return nil
}

func (m startingUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case " ":
			m.selected = m.cursor

		case "enter":
			globalState := state.GlobalStateInstance()
			globalState.SetDownloadType(m.choices[m.selected].videoType)
			return InitVideoIdSearchModel(), nil

		}

	}
	return m, nil
}

func (m startingUIModel) View() string {
	s := ""
	s = headerStyle("Welcome to You Download")
	s += "\n\nPlease select the type of download you want to do:\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		// Is this choice selected?
		checked := " " // not selected
		if i == m.selected {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.text)
	}
	s += "\n Press \"Enter\" to continue"
	s += "\n Press \"Space bar\" to select an Option"
	s += "\n Press Q or Ctrl+c to quit. \n"

	return s
}
