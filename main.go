package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	videodownload "github.com/anjolaoluwaakindipe/fyne-youtube/videoDownload"
	"github.com/kkdai/youtube/v2"
)

func main(){

	youtubeClient := youtube.Client{}

	fmt.Println("Please Input a video Id from youtube")
	fmt.Print("-> ")
	reader := bufio.NewReader(os.Stdin)
	
	videoId, stringReadErr := reader.ReadString('\n');

	if stringReadErr != nil{
		log.Fatal(stringReadErr.Error())
		return
	}

	videoId = strings.Replace(videoId, "\r\n", "", 1)

	videoDownloader := videodownload.InitViedoDownload(videoId, &youtubeClient)
	videoDownloader.Download()

	
	fmt.Println("")
	fmt.Println("Download finished")

	
	

	
}