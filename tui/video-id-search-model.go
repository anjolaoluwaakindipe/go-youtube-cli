package tui

import (
	"fmt"
	"strings"

	videodownload "github.com/anjolaoluwaakindipe/fyne-youtube/videodownload"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/anjolaoluwaakindipe/fyne-youtube/tui/state"

	"github.com/charmbracelet/bubbles/textinput"
)

// possible application msgs
type errorMsg error

// model/state
type videoIdSearchModel struct {
	download  videodownload.VideoDownload
	textInput textinput.Model
	err       error
}

type StartVideoIdSearch struct{}

// constructor
func InitVideoIdSearchModel(videoDownload videodownload.VideoDownload) *videoIdSearchModel {
	ti := textinput.New()
	ti.Placeholder = "e.g. YalT4KKnLao"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return &videoIdSearchModel{download: videoDownload, textInput: ti}
}

// Initialization function
func (vm* videoIdSearchModel) Init() tea.Cmd {
	return textinput.Blink
}

// Application UI
func (vm* videoIdSearchModel) View() string {
	s := ""

	s = fmt.Sprintf("\n Pleas input the %s id", vm.download.GetType())

	if vm.err != nil {
		s = "\n An error occurred while typing\n \n Please quit and try again"
	} else {
		s += vm.textInput.View()
	}

	s += "\n \n  Press Ctrl+c to quit. \n"

	return s
}

// event listener
func (vm* videoIdSearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// input keys
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return vm, tea.Quit

		case "enter":
			if strings.ReplaceAll(strings.ReplaceAll(vm.textInput.Value(), " ", ""), "\r\n", "") == "" {
				return vm, nil
			}
			globalState := state.GlobalStateInstance()
			globalState.SetVideoId(strings.ReplaceAll(vm.textInput.Value(), "\n", ""))
			nextModel := InitializeDownloadLocationModel(vm.download)
			return nextModel, func() tea.Msg {
				return StartDownloadLocationModel{}
			}
		}
	case StartVideoIdSearch:
		vm.Init()
	// error handling
	case errorMsg:
		vm.err = msg
		return vm, nil
	}

	// create a new cmd for textinput
	var cm tea.Cmd

	// update textinput and cmd
	vm.textInput, cm = vm.textInput.Update(msg)

	return vm, cm
}
