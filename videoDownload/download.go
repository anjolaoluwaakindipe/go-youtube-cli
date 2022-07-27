package videodownload

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gosuri/uilive"
	"github.com/kkdai/youtube/v2"
)

// VideoDownload properties
type VideoDownload struct {
	videoUrl      string
	videoClient   *youtube.Client
	consoleWriter *uilive.Writer
	mu            sync.Mutex
}

// Video Download Constructor
func InitViedoDownload(videoUrl string, videoClient *youtube.Client) *VideoDownload {
	return &VideoDownload{videoUrl: videoUrl, videoClient: videoClient, consoleWriter: uilive.New()}
}

/* Video Download Constructors */

// Gets the video stream from youtube
func (vd *VideoDownload) getStream(downloadType DownloadType) (video *youtube.Video, stream io.ReadCloser, videoSize int64, getStreamError error) {

	var videoFetchingErr error

	// get video information and metadata for a specific download type e.g A single video or a playlist
	if downloadType == SingleVideo {
		video, videoFetchingErr = vd.videoClient.GetVideo(vd.videoUrl)
	} else {
		return nil, nil, 0, errors.New("did not give a valid video type")
	}

	if videoFetchingErr != nil {
		log.Fatal(videoFetchingErr.Error())
		return nil, nil, 0, videoFetchingErr
	}
	// audio channel with format list
	format := video.Formats.WithAudioChannels()

	// stream, total youtube video downloaded and stream error
	stream, videoSize, getStreamError = vd.videoClient.GetStream(video, &format[0])
	return
}

// shows the download progress of the yooutbe video
func (vd *VideoDownload) showDownloadProgress(file *os.File, expectedSize int64) {

	go func() {

		for {
			vd.mu.Lock()

			fileInfo, _ := file.Stat()
			fmt.Fprintf(vd.consoleWriter, "Status: %v mb / %v mb \n", fileInfo.Size()/1000000, expectedSize/1000000)
			if fileInfo.Size() == expectedSize {
				break
			}
			time.Sleep(time.Millisecond * 300)
			vd.mu.Unlock()
		}
		fmt.Fprintf(vd.consoleWriter, "Finished downloading, Total video size: %v mb\n", expectedSize)
		vd.consoleWriter.Stop()
	}()
}

// download the video stream into a file
func (vd *VideoDownload) Download() {

	vd.consoleWriter.Start()
	fmt.Print("Getting youtube video stream \n")
	video, stream, videoSize, streamErr := vd.getStream(SingleVideo)

	fmt.Printf("Starting download for %v by %v... \n", video.Title, video.Author)

	if streamErr != nil {
		log.Fatal(streamErr.Error())
		return
	}

	dowloadedFileName := strings.ReplaceAll(video.Title, " ", "_") + ".mp4"

	file, fileCreationErr := os.Create(dowloadedFileName)

	if fileCreationErr != nil {
		log.Fatal(fileCreationErr.Error())
		return
	}

	vd.showDownloadProgress(file, videoSize)
	// channels
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, syscall.SIGTERM)

	go func() {
		signal := <-signalChan
		fmt.Println(signal)
		fileInfo, _ := file.Stat()
		fmt.Println(fileInfo.Size())
		fmt.Println(videoSize)
		if fileInfo.Size() < videoSize {
			file.Close()
			fileRemoverErr := os.Remove(dowloadedFileName)
			if fileRemoverErr != nil {
				log.Fatal("File Removal Error: " + fileRemoverErr.Error())
				return
			}
			os.Exit(0)

		}

	}()

	_, fileCopyErr := io.Copy(file, stream)

	if fileCopyErr != nil {

		log.Print("File Copy Error: " + fileCopyErr.Error())
		return
	}

	vd.consoleWriter.Stop()

	file.Close()
}
