package app

import (
	videodownload "github.com/anjolaoluwaakindipe/fyne-youtube/videoDownload"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/anjolaoluwaakindipe/fyne-youtube/app/state"

	"github.com/charmbracelet/bubbles/textinput"
)

// possible application msgs
type errorMsg error

// model/state
type videoIdSearchModel struct {
	videoType videodownload.DownloadType
	textInput textinput.Model
	err error
}

// constructor
func InitVideoIdSearchModel() videoIdSearchModel {
	ti := textinput.New()
	ti.Placeholder = "e.g. YalT4KKnLao"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	globalState := state.GlobalStateInstance()
	return videoIdSearchModel{videoType: globalState.GetDownloadType(), textInput: ti}
}

// Initialization function
func (vm videoIdSearchModel) Init() tea.Cmd {
	return textinput.Blink
}

// Application UI 
func (vm videoIdSearchModel) View() string {
	s:= ""
	switch {
	case vm.videoType == videodownload.SingleVideo:
		s = "\n Please input the video id \n \n"
	case vm.videoType == videodownload.PlayList:
		s= "\n Please input the playlist id\n \n"
	}

	if vm.err != nil{
		s = "\n An error occurred while typing\n \n Please quit and try again"
	}else{
		s += vm.textInput.View()
	}

	s += "\n \n  Press Ctrl+c to quit. \n"

	return s
}

// event listener
func (vm videoIdSearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// input keys 
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return vm, tea.Quit

		case "enter":
			globalState := state.GlobalStateInstance()
			globalState.SetVideoId(vm.textInput.Value())
			return InitializeDownloadLocationModel(), nil
		}

	// error handling
	case errorMsg:
		vm.err = msg
		return vm, nil
	}

	// create a new cmd for textinput
	var cm tea.Cmd

	// update textinput and cmd
	vm.textInput , cm = vm.textInput.Update(msg)

	return vm, cm
}
