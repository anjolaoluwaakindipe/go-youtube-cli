package app

import (
	"fmt"
	"time"

	"github.com/anjolaoluwaakindipe/fyne-youtube/app/state"
	videodownload "github.com/anjolaoluwaakindipe/fyne-youtube/videoDownload"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

const (
	padding  = 2
	maxWidth = 80
)

// messages
type tickMsg time.Time

// state / model
type singleVideoDownloadModel struct {
	downloadType videodownload.DownloadType
	videoId      string
	progress     progress.Model
	directory    string
	numberChan chan float64
}

// constructor
func InitializeSingleVideoDownloadModel() *singleVideoDownloadModel {
	globalState := state.GlobalStateInstance()
	return &singleVideoDownloadModel{videoId: globalState.GetVideoId(), downloadType: globalState.GetDownloadType(), progress: progress.New(progress.WithDefaultGradient()), directory: globalState.GetDownloadDirectory()}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// init command func
func (singleVideoDownloadModel) Init() tea.Cmd {
	return nil
}

// UI layer
func (sm singleVideoDownloadModel) View() string {
	s := ""

	s = fmt.Sprintf("\n your video  id is %v and directory is %v \n\n", sm.videoId, sm.directory)

	s += (sm.progress.View() + "\n\n")

	s += "\n \n  Press Q or Ctrl+c to quit. \n"

	return s
}

// event listener
func (sm singleVideoDownloadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return sm, tea.Quit

	case tea.WindowSizeMsg:
		sm.progress.Width = msg.Width - padding*2 - 4
		if sm.progress.Width > maxWidth {
			sm.progress.Width = maxWidth
		}
		return sm, nil

	case tickMsg:
		if sm.progress.Percent() == 1.0 {
			return sm, tea.Quit
		}
		
		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
		cmd := sm.progress.IncrPercent(0.25)
		return sm, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := sm.progress.Update(msg)
		sm.progress = progressModel.(progress.Model)
		return sm, cmd

	default:
		return sm, nil
	}

}

func(sm singleVideoDownloadModel) test(numberChan chan float64){
	go func(){
		for i:=1; i <=4 ; i++{
			numberChan <- 0.25
			time.Sleep(time.Second * 5);
		}
	}()
}
