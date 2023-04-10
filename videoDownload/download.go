package videodownload

import (
	"os"

	"github.com/anjolaoluwaakindipe/fyne-youtube/appmsg"
)

// Download type enum
type DownloadType int

const (
	SingleVideo DownloadType = iota
	PlayList
)

func (dt DownloadType) String() string {
	switch dt {
	case SingleVideo:
		return "Video"
	case PlayList:
		return "Playlist"
	}
	return "unknown"
}

type VideoDownload interface {
	GetType() string
	Download(videoId string, directoryPath string, progressChan *chan appmsg.DownloadProgressMsg)
	CancelVideoDownload(file *os.File, directory string, filename string)
}
