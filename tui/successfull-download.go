package tui

import (
	"fmt"

	"github.com/anjolaoluwaakindipe/fyne-youtube/appmsg"
	tea "github.com/charmbracelet/bubbletea"
)

type successfulDownloadModel struct {
	choices  []RetryOptions
	cursor   int
	selected int
}

type RetryOptions struct {
	text    string
	message tea.Msg
}

func InitializeSuccessfulDownloadModel() successfulDownloadModel {
	options := []RetryOptions{{text: "Download another video?", message: appmsg.DownloadAnotherVideoMsg{}}}

	return successfulDownloadModel{choices: options, selected: 0}
}

func (sdm successfulDownloadModel) Init() tea.Cmd {
	return nil
}

func (sdm successfulDownloadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return sdm, tea.Quit
		case "down", "j":
			if sdm.cursor < len(sdm.choices)-1 {
				sdm.cursor++
			}

		case "up", "k":
			if sdm.cursor > 0 {
				sdm.cursor--
			}

		case " ":
			sdm.selected = sdm.cursor

		case "enter":
			return sdm, func() tea.Msg { return sdm.choices[sdm.selected].message }
		}

	case appmsg.DownloadAnotherVideoMsg:
		return InitialStartingUIModel(), nil
	case appmsg.QuitMsg:
		return sdm, tea.Quit
	}

	return sdm, nil
}

func (sdm successfulDownloadModel) View() string {
	s := ""

	s += "Video downloaded successfully!!! \n \n "

	for i, choice := range sdm.choices {
		cursor := " "
		if sdm.cursor == i {
			cursor = ">"
		}

		// if choices
		checked := " "
		if i == sdm.selected {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s \n", cursor, checked, choice.text)
	}

	s += "\n Press \"Q or Ctrl+c\" to quit."

	return s
}
