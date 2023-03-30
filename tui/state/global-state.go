package state

import (
	"sync"

	videodownload "github.com/anjolaoluwaakindipe/fyne-youtube/videodownload"
)

type globalState struct {
	videoId           string
	downloadType      videodownload.DownloadType
	downloadDirectory string
}

// gloabal state singleton
var instance globalState

// var
var once sync.Once

func GlobalStateInstance() *globalState{
	once.Do(func() {
		instance = globalState{}
	})

	return &instance
}

// getters 
func (gs *globalState) GetVideoId () string{
	return gs.videoId
}
func (gs *globalState) GetDownloadType () videodownload.DownloadType{
	return gs.downloadType
}
func (gs *globalState) GetDownloadDirectory () string{
	return gs.downloadDirectory
}

// setters
func (gs *globalState) SetVideoId (newVideoId string) {
	gs.videoId = newVideoId
}
func (gs *globalState) SetDownloadType (newDownloadType videodownload.DownloadType) {
	gs.downloadType = newDownloadType
}
func (gs *globalState) SetDownloadDirectory (newDonwloadDirectory string) {
	gs.downloadDirectory = newDonwloadDirectory
}