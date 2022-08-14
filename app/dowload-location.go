package app

import (
	"github.com/anjolaoluwaakindipe/fyne-youtube/app/state"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// state / model
type downloadLocationModel struct {
	textInput textinput.Model
}

// constructor
func InitializeDownloadLocationModel() downloadLocationModel {

	ti := textinput.New()
	ti.Placeholder = "e.g. C:\\Users\\<User>\\Video"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return downloadLocationModel{textInput: ti}
}

// init command func
func (sm downloadLocationModel) Init() tea.Cmd {
	return nil
}

// UI layer
func (sm downloadLocationModel) View() string {

	s := "\n Please input a Directory for your Download  \n \n"

	s += sm.textInput.View()

	s += "\n \n  Press Ctrl+c to quit. \n"

	return s
}

// event listener
func (sm downloadLocationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return sm, tea.Quit
		case "enter":
			globalState := state.GlobalStateInstance()
			globalState.SetDownloadDirectory(sm.textInput.Value())
			return InitializeSingleVideoDownloadModel(), tickCmd()
		}
	}

	var cm tea.Cmd

	sm.textInput, cm = sm.textInput.Update(msg)

	return sm, cm
}
