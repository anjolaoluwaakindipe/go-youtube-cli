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

type PlaylistDownload struct {
	videoClient *youtube.Client
	mu          sync.Mutex
	videos      []*youtube.PlaylistEntry
	doneCount   int
	utils       utils.Utils
}

func InitPlaylistVideoDownload() VideoDownload {
	return &PlaylistDownload{videoClient: &youtube.Client{}, utils: utils.NewUtilsImpl()}
}

func (pd *PlaylistDownload) GetType() string {
	return PlayList.String()
}

func (pd *PlaylistDownload) Download(videoId string, directoryPath string, progressChan *chan appmsg.DownloadProgressMsg) {
	playlist, err := pd.videoClient.GetPlaylist(videoId)
	if err != nil {
		log.Fatalf("Failed to fetch playlist with id %v: %v", videoId, err)
	}

	playlistDirectory := directoryPath + string(os.PathSeparator) + playlist.Title

	err = os.Mkdir(playlistDirectory, os.ModePerm)

	if err != nil {
		log.Fatalf("Failed to create directory %v : %v", playlistDirectory, err)
	}

	pd.mu.Lock()
	pd.videos = playlist.Videos
	pd.mu.Unlock()

	var wg sync.WaitGroup

	pd.DownloadAllVideos(wg, playlistDirectory, progressChan)
}

func (pd *PlaylistDownload) DownloadAllVideos(wg sync.WaitGroup, playlistDirectory string, progress *chan appmsg.DownloadProgressMsg) {
	for i, video := range pd.videos {
		videoInfo, err := pd.videoClient.GetVideo(video.ID)
		if err != nil {
			log.Fatalf("Failed to fetch video metadata: %v", err)
			continue
		}

		filename := pd.utils.ConvertVideoNameToProperFileName(video.Title)

		wg.Add(1)

		go func(index int, videoInfo *youtube.Video, filename string) {
			defer wg.Done()

			videoFormat := videoInfo.Formats.WithAudioChannels()
			stream, videoSize, getStreamError := pd.videoClient.GetStream(videoInfo, &videoFormat[0])

			if getStreamError != nil {
				log.Fatalf("File Copy Error: %v", getStreamError)
			}

			file, fileCreationErr := os.Create(playlistDirectory + string(os.PathSeparator) + filename)

			if fileCreationErr != nil {
				log.Fatalf("File Creation Error: %v", fileCreationErr)
			}

			pd.showDownloadProgress(index, file, videoSize, videoInfo, filename, progress)

			_, filecopyErr := io.Copy(file, stream)

			if filecopyErr != nil {
				log.Fatalf("File Copy Error: %v", filecopyErr)
			}
		}(i, videoInfo, filename)
	}
}

// shows the download progress of the yooutbe video
func (pd *PlaylistDownload) showDownloadProgress(index int, file *os.File, expectedSize int64, video *youtube.Video, downloadedFileName string, progressChan *chan appmsg.DownloadProgressMsg) {
	// run the concurrent function
	go func() {
		for {
			if file == nil {
				break
			}
			// make a mutex lock so to prevent simultaneous access
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
			pd.mu.Lock()
			*progressChan <- appmsg.DownloadProgressMsg{Index: index, Progress: progress, TotalDownloadSize: expectedSize, AmountDownloaded: amountDownloaded, VideoName: video.Title, VideoAuthor: video.Author, DefaultFileName: downloadedFileName, VideoFile: file, IsDone: false, VideosDownloaded: pd.doneCount, AllVideos: len(pd.videos)}
			pd.mu.Unlock()
			// check if they are equal thus breaking the loop
			if amountDownloaded == expectedSize {
				pd.mu.Lock()
				pd.doneCount += 1
				pd.mu.Unlock()
				if pd.doneCount == len(pd.videos) {
					pd.mu.Lock()
					*progressChan <- appmsg.DownloadProgressMsg{Index: index, Progress: progress, TotalDownloadSize: expectedSize, AmountDownloaded: amountDownloaded, VideoName: video.Title, VideoAuthor: video.Author, DefaultFileName: downloadedFileName, VideoFile: file, IsDone: true, VideosDownloaded: pd.doneCount, AllVideos: len(pd.videos)}
					pd.mu.Unlock()
				}
				file.Close()
				break
			}

			// checks every 500 millisecond
			time.Sleep(time.Millisecond * 500)

			// unlocks the mutex
		}
	}()
}
func (pd *PlaylistDownload) CancelVideoDownload(file *os.File, directory string, filename string) {}
