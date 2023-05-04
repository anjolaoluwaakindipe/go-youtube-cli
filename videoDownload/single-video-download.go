package videodownload

import (
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/anjolaoluwaakindipe/fyne-youtube/appmsg"
	"github.com/anjolaoluwaakindipe/fyne-youtube/utils"
	"github.com/kkdai/youtube/v2"
)

// SingleVideoDownload properties
type SingleVideoDownload struct {
	videoClient *youtube.Client
	mu          sync.Mutex
	utils       utils.Utils
}

// Video Download Constructor
func InitSingleVideoDownload() VideoDownload {
	return &SingleVideoDownload{videoClient: &youtube.Client{}, utils: utils.NewUtilsImpl()}
}

/* Video Download Constructors */

// Gets the video stream from youtube
func (vd *SingleVideoDownload) getStream(videoUrl string) (video *youtube.Video, stream io.ReadCloser, videoSize int64, getStreamError error) {
	var videoFetchingErr error

	// get video information and metadata for a specific download type e.g A single video or a playlist
	video, videoFetchingErr = vd.videoClient.GetVideo(videoUrl)

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
func (vd *SingleVideoDownload) showDownloadProgress(file *os.File, expectedSize int64, video *youtube.Video, downloadedFileName string, progressChan *chan appmsg.DownloadProgressMsg) {
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

			// makes sure that progress messages are sent to the ui
			*progressChan <- appmsg.DownloadProgressMsg{Progress: progress, TotalDownloadSize: expectedSize, AmountDownloaded: amountDownloaded, VideoName: video.Title, VideoAuthor: video.Author, DefaultFileName: downloadedFileName, VideoFile: file, IsDone: false}
			// check if they are equal thus breaking the loop
			if amountDownloaded == expectedSize {
				*progressChan <- appmsg.DownloadProgressMsg{Progress: progress, TotalDownloadSize: expectedSize, AmountDownloaded: amountDownloaded, VideoName: video.Title, VideoAuthor: video.Author, DefaultFileName: downloadedFileName, VideoFile: file, IsDone: true}
				break
			}

			// checks every 500 millisecond
			time.Sleep(time.Millisecond * 500)

			// unlocks the mutex
			vd.mu.Unlock()
		}
		file.Close()
	}()
}

// download the video stream into a file
func (vd *SingleVideoDownload) Download(videoId string, directoryPath string, progressChan *chan appmsg.DownloadProgressMsg) {
	// call stream method to get the video stream
	video, stream, videoSize, streamErr := vd.getStream(videoId)

	// check for stream error
	if streamErr != nil {
		log.Fatalln("File Streaming Error: " + streamErr.Error() + " 1 ")
	}

	// format the filename so that it is a valid name that can be used with any os
	downloadedFileName := vd.utils.ConvertVideoNameToProperFileName(video.Title)

	// create the file
	file, fileCreationErr := os.Create(directoryPath + string(os.PathSeparator) + downloadedFileName)

	// cheeck for errors while creating the file
	if fileCreationErr != nil {
		log.Fatalln("File Creation Error: " + fileCreationErr.Error() + " 2 ")
	}
	// begin display the downloand progress of the video
	vd.showDownloadProgress(file, videoSize, video, downloadedFileName, progressChan)

	_, fileCopyErr := io.Copy(file, stream)

	// check for any errors that may occur while downloading
	if fileCopyErr != nil {
		log.Fatal("File Copy Error: " + fileCopyErr.Error())
	}
}

func (vd *SingleVideoDownload) CancelVideoDownload(file *os.File, directory string, filename string) {
	file.Close()
	os.Remove(directory + string(os.PathSeparator) + filename)
}

func (vd *SingleVideoDownload) GetType() string {
	return SingleVideo.String()
}
