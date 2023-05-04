package tui

import (
	"fmt"
	"os"
	"strconv"

	"github.com/anjolaoluwaakindipe/fyne-youtube/app"
	"github.com/anjolaoluwaakindipe/fyne-youtube/appmsg"
	"github.com/anjolaoluwaakindipe/fyne-youtube/tui/state"
	"github.com/anjolaoluwaakindipe/fyne-youtube/videodownload"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// each video Progress
type Progress struct {
	Index            int
	downloadSize     int64
	amountDownloaded int64
	videoName        string
	videoFile        *os.File
	directory        string
}

type PlaylistDownloadModel struct {
	downloadType     videodownload.VideoDownload
	directory        string
	videoId          string
	allVideos        int
	downloadedVideos int
	downloadProgress map[string]Progress
	model            tea.Model
	table            table.Model
	rows             []table.Row
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type StartPlaylistDownload struct{}

func InitPlaylistDownloadModel(videoDownload videodownload.VideoDownload) *PlaylistDownloadModel {
	globalState := state.GlobalStateInstance()
	columns := []table.Column{{Title: "S/N", Width: 5}, {Title: "Name", Width: 30}, {Title: "Progress", Width: 20}, {Title: "Done", Width: 5}}
	rows := make([]table.Row, 0)
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	return &PlaylistDownloadModel{videoId: globalState.GetVideoId(), downloadType: videoDownload, directory: globalState.GetDownloadDirectory(), downloadProgress: make(map[string]Progress), table: t}
}

func (pdm *PlaylistDownloadModel) Init() tea.Cmd {
	progressChan := make(chan appmsg.DownloadProgressMsg)

	go pdm.downloadType.Download(pdm.videoId, pdm.directory, &progressChan)

	go func() {
		for {
			progress := <-progressChan
			app.TuiProgram.Send(progress)
			if progress.IsDone {
				app.TuiProgram.Send(appmsg.DownloadComplete{})
				close(progressChan)
				return
			}
		}
	}()

	return nil
}

func (pdm *PlaylistDownloadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if pdm.table.Focused() {
				pdm.table.Blur()
			} else {
				pdm.table.Focus()
			}
		case tea.KeyUp.String():
			pdm.table.MoveUp(0)
		case tea.KeyDown.String():
			pdm.table.MoveDown(0)
		case "q", "ctrl+c":
			var cmds []tea.Cmd
			cmds = append(cmds, tea.Sequence(pdm.DeleteVideo, tea.Quit))
			return pdm, tea.Batch(cmds...)
		default:
			return pdm, nil
		}

	case appmsg.DownloadProgressMsg:
		if len(pdm.rows) < 1 {
			pdm.rows = make([]table.Row, msg.AllVideos)
		}
		var done string
		if msg.AmountDownloaded == msg.TotalDownloadSize {
			done = "O"
		} else {
			done = "-"
		}
		pdm.rows[msg.Index] = table.Row{strconv.Itoa(msg.Index + 1), msg.VideoName, fmt.Sprintf("%v/%v", msg.AmountDownloaded, msg.TotalDownloadSize), done}
		pdm.table.SetRows(pdm.rows)
		pdm.downloadProgress[msg.DefaultFileName] = Progress{downloadSize: msg.TotalDownloadSize, amountDownloaded: msg.AmountDownloaded, videoName: msg.VideoName, videoFile: msg.VideoFile, directory: state.GlobalStateInstance().GetDownloadDirectory(), Index: msg.Index}
		pdm.allVideos = msg.AllVideos
		pdm.downloadedVideos = msg.VideosDownloaded
		return pdm, nil

	case appmsg.DownloadComplete:
		return InitializeSuccessfulDownloadModel(), nil

	case StartPlaylistDownload:
		pdm.Init()
		return pdm, nil
	}
	pdm.table, cmd = pdm.table.Update(msg)
	return pdm, cmd
}

func (pdm *PlaylistDownloadModel) View() string {
	var s string = ""

	// for _, val := range pdm.downloadProgress {
	// 	s += fmt.Sprintf("%v Downloading %v : %v/%v\n\n", val.Index, val.videoName, val.amountDownloaded, val.amountDownloaded)
	// }
	s += baseStyle.Render(pdm.table.View()) + "\n"
	s += fmt.Sprintf("%v out of %v downloaded", pdm.downloadedVideos, pdm.allVideos)

	return s
}

func (pdm *PlaylistDownloadModel) DeleteVideo() tea.Msg {
	for _, val := range pdm.downloadProgress {
		pdm.downloadType.CancelVideoDownload(val.videoFile, val.directory, val.videoName)
	}
	return nil
}

func (pdm *PlaylistDownloadModel) GetModel() tea.Model {
	return pdm.model
}
