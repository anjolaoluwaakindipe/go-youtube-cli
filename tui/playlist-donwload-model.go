package tui

import (
	"fmt"
	"os"

	"github.com/anjolaoluwaakindipe/fyne-youtube/app"
	"github.com/anjolaoluwaakindipe/fyne-youtube/appmsg"
	"github.com/anjolaoluwaakindipe/fyne-youtube/tui/state"
	"github.com/anjolaoluwaakindipe/fyne-youtube/videodownload"
	tea "github.com/charmbracelet/bubbletea"
)

// each video Progress
type Progress struct {
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
}

type StartPlaylistDownload struct {
}

func InitPlaylistDownloadModel(videoDownload videodownload.VideoDownload) *PlaylistDownloadModel {

	globalState := state.GlobalStateInstance()
	return &PlaylistDownloadModel{videoId: globalState.GetVideoId(), downloadType: videoDownload, directory: globalState.GetDownloadDirectory(), downloadProgress: make(map[string]Progress)}
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c":
			var cmds []tea.Cmd
			cmds = append(cmds, tea.Sequence(pdm.DeleteVideo, tea.Quit))
			return pdm, tea.Batch(cmds...)
		default:
			return pdm, nil
		}

	case appmsg.DownloadProgressMsg:
		pdm.downloadProgress[msg.DefaultFileName] = Progress{downloadSize: msg.TotalDownloadSize, amountDownloaded: msg.AmountDownloaded, videoName: msg.VideoName, videoFile: msg.VideoFile, directory: state.GlobalStateInstance().GetDownloadDirectory()}
		pdm.allVideos = msg.AllVideos
		pdm.downloadedVideos = msg.VideosDownloaded

	case appmsg.DownloadComplete:
		return InitializeSuccessfulDownloadModel(), nil

	case StartPlaylistDownload:
		pdm.Init()
		return pdm, nil
	}

	return pdm, nil
}

func (pdm *PlaylistDownloadModel) View() string {
	var s string = ""

	// for _, val := range pdm.downloadProgress {
	// 	s += fmt.Sprintf("Downloading %v : %v/%v\n\n", val.videoName, val.amountDownloaded, val.amountDownloaded)
	// }
	s += fmt.Sprintf("%v out of %v downloaded", pdm.downloadedVideos, pdm.allVideos)

	return s
}

func (pdm *PlaylistDownloadModel) DeleteVideo() tea.Msg {
	for _, val := range pdm.downloadProgress {
		pdm.downloadType.CancelVideoDownload(val.videoFile, val.directory, val.videoName)
	}
	return nil
}
