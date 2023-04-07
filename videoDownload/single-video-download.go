package videodownload

import (
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/anjolaoluwaakindipe/fyne-youtube/app"
	"github.com/anjolaoluwaakindipe/fyne-youtube/appmsg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kkdai/youtube/v2"
)

// SingleVideoDownload properties
type SingleVideoDownload struct {
	videoClient *youtube.Client
	mu          sync.Mutex
	program     *tea.Program
}

// Video Download Constructor
func InitSingleVideoDownload() VideDownload {
	return &SingleVideoDownload{videoClient: &youtube.Client{}, program: app.TuiProgram}
}

/* Video Download Constructors */

// Gets the video stream from youtube
func (vd *SingleVideoDownload) getStream(downloadType DownloadType, videoUrl string) (video *youtube.Video, stream io.ReadCloser, videoSize int64, getStreamError error) {
	var videoFetchingErr error

	// get video information and metadata for a specific download type e.g A single video or a playlist
	if downloadType == SingleVideo {
		video, videoFetchingErr = vd.videoClient.GetVideo(videoUrl)
	} else {
		return nil, nil, 0, errors.New("did not give a valid video type")
	}

	if videoFetchingErr != nil {
		log.Fatalln("Video Fetching error: " + videoFetchingErr.Error())
		return nil, nil, 0, videoFetchingErr
	}
	// audio channel with format list
	format := video.Formats.WithAudioChannels()

	// stream, total youtube video downloaded and stream error
	stream, videoSize, getStreamError = vd.videoClient.GetStream(video, &format[0])

	return
}

// shows the download progress of the yooutbe video
func (vd *SingleVideoDownload) showDownloadProgress(file *os.File, expectedSize int64, video *youtube.Video, downloadedFileName string) {
	// run the concurrent function
	go func() {
		for {

			if file == nil {
				break
			}
			// make a mutex lock so to prevent simultaneous access
			vd.mu.Lock()
			// get file info
			if file == nil {
				break
			}
			fileInfo, _ := file.Stat()
			// get the amount downloaded from the size of the file created and the expected size from the stream
			if file == nil {
				break
			}
			amountDownloaded := fileInfo.Size()
			progress := float64(amountDownloaded) / float64(expectedSize)

			// check if they are equal thus breaking the loop
			if amountDownloaded == expectedSize {
				break
			}

			// makes sure that progress messages are sent to the ui
			if vd.program != nil {
				vd.program.Send(appmsg.DownloadProgressMsg{Progress: progress, TotalDownloadSize: expectedSize, AmountDownloaded: amountDownloaded, VideoName: video.Title, VideoAuthor: video.Author, DefaultFileName: downloadedFileName, VideoFile: file})
			}

			// checks every 500 millisecond
			time.Sleep(time.Millisecond * 500)

			// unlocks the mutex
			vd.mu.Unlock()
		}
		file.Close()
		vd.program.Send(appmsg.DownloadComplete{})
	}()
}

// download the video stream into a file
func (vd *SingleVideoDownload) Download(videoId string, directoryPath string) {
	// call stream method to get the video stream
	video, stream, videoSize, streamErr := vd.getStream(SingleVideo, videoId)

	// check for stream error
	if streamErr != nil {
		log.Fatalln("File Streaming Error: " + streamErr.Error() + " 1 ")
	}

	// format the filename so that it is a valid name that can be used with any os
	downloadedFileName := strings.ReplaceAll(video.Title, " ", "_") + ".mp4"

	// create the file
	file, fileCreationErr := os.Create(directoryPath + "/" + downloadedFileName)

	// cheeck for errors while creating the file
	if fileCreationErr != nil {
		log.Fatalln("File Creation Error: " + fileCreationErr.Error() + " 2 ")
	}

	// begin display the downloand progress of the video
	vd.showDownloadProgress(file, videoSize, video, downloadedFileName)

	_, fileCopyErr := io.Copy(file, stream)

	// check for any errors that may occur while downloading
	if fileCopyErr != nil {
		log.Fatal("File Copy Error: " + fileCopyErr.Error())
	}

}

func (vd *SingleVideDownload) CancelVideoDownload(file *os.File, directory string, filename string) {
	file.Close()
	os.Remove(directory + string(os.PathSeparator) + filename)
}

func (vd *SingleVideDownload) GetType() string {
	return "Single Video"
}
