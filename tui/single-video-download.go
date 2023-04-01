package tui

import (
	"fmt"
	"os"
	"time"

	"github.com/anjolaoluwaakindipe/fyne-youtube/appmsg"
	"github.com/anjolaoluwaakindipe/fyne-youtube/tui/state"
	videodownload "github.com/anjolaoluwaakindipe/fyne-youtube/videodownload"
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
type (
	tickMsg      time.Time
	incrementMsg float64
)

// state / model
type singleVideoDownloadModel struct {
	downloadType       videodownload.DownloadType
	videoId            string
	progress           progress.Model
	amountDownloaded   int64
	totalDownloadSize  int64
	directory          string
	videoName          string
	videoAuthor        string
	DownloadedFileName string
	VideoFile          *os.File
}

// constructor
func InitializeSingleVideoDownloadModel() *singleVideoDownloadModel {
	globalState := state.GlobalStateInstance()
	return &singleVideoDownloadModel{videoId: globalState.GetVideoId(), downloadType: globalState.GetDownloadType(), progress: progress.New(progress.WithDefaultGradient()), directory: globalState.GetDownloadDirectory()}
}

// func tickCmd() tea.Cmd {
// 	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
// 		return tickMsg(t)
// 	})
// }

func pauseProgress() tea.Cmd {
	return tea.Tick(time.Millisecond*750, func(t time.Time) tea.Msg {
		return nil
	})
}

// init command func
func (singleVideoDownloadModel) Init() tea.Cmd {
	return nil
}

// UI layer
func (sm singleVideoDownloadModel) View() string {
	s := ""

	s = fmt.Sprintf("\n You are downloading %v by %v  into  %v \n\n", sm.videoName, sm.videoAuthor, sm.directory)

	s += (sm.progress.View() + "\n\n")

	if sm.totalDownloadSize > 0 {
		s += fmt.Sprintf("\n\n Total amount downloaded: %vmb/%vmb \n", sm.amountDownloaded/1000000, sm.totalDownloadSize/1000000)
	}

	s += "\n \n  Press Q or Ctrl+c to quit. \n"

	return s
}

// event listener
func (sm singleVideoDownloadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			var cmds []tea.Cmd
			if sm.amountDownloaded != sm.totalDownloadSize && sm.amountDownloaded > 0 {
				cmds = append(cmds, tea.Sequentially(sm.DeleteVideo, tea.Quit))
				return sm, tea.Batch(cmds...)
			}
			return sm, tea.Quit
		default:
			return sm, nil
		}

	case tea.WindowSizeMsg:
		sm.progress.Width = msg.Width
		if sm.progress.Width > maxWidth {
			sm.progress.Width = maxWidth
		}
		return sm, nil

	// case tickMsg:
	// 	if sm.progress.Percent() == 1.0 {
	// 		return sm, tea.Quit
	// 	}

	// 	// Note that you can also use progress.Model.SetPercent to set the
	// 	// percentage value explicitly, too.
	// 	cmd := sm.progress.IncrPercent(0.25)
	// 	return sm, tea.Batch(tickCmd(), cmd)

	case appmsg.DownloadProgressMsg:
		var cmds []tea.Cmd
		if msg.Progress >= 1.0 {
			cmds = append(cmds, tea.Sequentially(pauseProgress()))
		} else {
			cmd := sm.progress.SetPercent(float64(msg.Progress))
			sm.amountDownloaded = msg.AmountDownloaded
			sm.totalDownloadSize = msg.TotalDownloadSize
			sm.videoAuthor = msg.VideoAuthor
			sm.videoName = msg.VideoName
			sm.DownloadedFileName = msg.DefaultFileName
			sm.VideoFile = msg.VideoFile

			cmds = append(cmds, cmd)
		}

		return sm, tea.Batch(cmds...)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:

		progressModel, cmd := sm.progress.Update(msg)
		sm.progress = progressModel.(progress.Model)
		return sm, cmd

	case appmsg.DownloadComplete:
		return InitializeSuccessfulDownloadModel(), nil

	default:
		return sm, nil
	}
}

// commans
func (sm *singleVideoDownloadModel) DeleteVideo() tea.Msg {
	sm.VideoFile.Close()
	os.Remove(sm.directory + string(os.PathSeparator) + sm.DownloadedFileName)
	return nil
}

// func test() tea.Msg {
// 	for i := 1; i <= 4; i++ {
// 		app.TuiProgram.Send(incrementMsg(0.25))
// 		time.Sleep(time.Second * 5)
// 	}
// 	return nil
// }
