package videodownload

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kkdai/youtube/v2"
)

type VideoDownload struct {
	videoUrl    string
	videoClient *youtube.Client
}

func InitViedoDownload(videoUrl string, videoClient *youtube.Client) *VideoDownload {
	return &VideoDownload{videoUrl: videoUrl, videoClient: videoClient}
}

func (vd *VideoDownload) Download() {
		video, videoFetchingErr := vd.videoClient.GetVideo(vd.videoUrl)
	
	if videoFetchingErr != nil {
		log.Fatal(videoFetchingErr.Error())
		return
	}

	format := video.Formats.WithAudioChannels();

	stream, bytesStreamed, streamErr := vd.videoClient.GetStream(video, &format[0])

	if streamErr != nil{
		log.Fatal(streamErr.Error())
		return
	}

	dowloadedFileName := video.Title + fmt.Sprintf(" (%v)", video.Author) + ".mp4";

	file, fileCreationErr := os.Create(dowloadedFileName);

	defer file.Close();

	if fileCreationErr != nil{
		log.Fatal(fileCreationErr.Error())
		return 
	}

	
	go func(){
		for(true){
			
			fileInfo, _ := file.Stat()
			fmt.Printf("\r Status: %v mb / %v mb",fileInfo.Size()/1000000, bytesStreamed/1000000)
		}
	}()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan , os.Interrupt)
	signal.Notify(signalChan , syscall.SIGTERM)
	go func(){
		signal := <- signalChan
		file.Close()
		fmt.Println(signal)
		fileInfo, _ := file.Stat()
		if fileInfo.Size() < bytesStreamed{
			fileRemoverErr := os.Remove(dowloadedFileName);
			if fileRemoverErr != nil{
				log.Fatal(fileRemoverErr.Error())
				return
			}
		}
	}()

	_, fileCopyErr := io.Copy(file, stream)

	if fileCopyErr != nil{
		
		log.Fatal(fileCopyErr.Error())
		return 
	}

}