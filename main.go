package main

import (
	"fmt"
	"os"

	"github.com/anjolaoluwaakindipe/fyne-youtube/app"
	"github.com/anjolaoluwaakindipe/fyne-youtube/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {


	// fmt.Println("Please Input a video Id from youtube")
	// fmt.Print("-> ")
	// reader := bufio.NewReader(os.Stdin)

	// videoId, stringReadErr := reader.ReadString('\n');

	// if stringReadErr != nil{
	// 	log.Fatal(stringReadErr.Error())
	// 	return
	// }

	// videoId = strings.Replace(videoId, "\r\n", "", 1)

	// videodownload.InitViedoDownload(videoId, &youtubeClient)

	app.TuiProgram = tea.NewProgram(tui.InitialStartingUIModel())
	if err:=app.TuiProgram.Start(); err!=nil {
		fmt.Printf("Alas, there's been an error: %v", err)
	    os.Exit(1)
	}
	// videoDownloader.Download()
	// fmt.Println("Starting download")
	// videodownload.VideoDownloadInstance.SingleVideoDownload("ljuKnv9D5JU")()

	// fmt.Println("")
	// fmt.Println("Download finished")

}
