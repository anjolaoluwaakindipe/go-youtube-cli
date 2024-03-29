package appmsg

import "os"

type DownloadProgressMsg struct {
	Progress          float64
	TotalDownloadSize int64
	AmountDownloaded  int64
	VideoName         string
	VideoAuthor       string
	DefaultFileName   string
	VideoFile         *os.File
	IsDone            bool
	VideosDownloaded  int
	AllVideos         int
	Index             int
}

type DownloadErrorMsg error

type DownloadAnotherVideoMsg struct{}

type QuitMsg struct{}

type DownloadComplete struct{}

type StartDownload struct{}
