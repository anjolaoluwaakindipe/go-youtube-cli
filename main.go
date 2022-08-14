package main

import (
	"fmt"
	"os"

	"github.com/anjolaoluwaakindipe/fyne-youtube/app"
	tea "github.com/charmbracelet/bubbletea"
)


func main(){

	// youtubeClient := youtube.Client{}

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

	p := tea.NewProgram(app.InitialStartingUIModel())
	if err:=p.Start(); err!=nil {
		fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
	}
	// videoDownloader.Download()

	
	// fmt.Println("")
	// fmt.Println("Download finished")

	
	

	
}