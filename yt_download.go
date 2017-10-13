package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	. "github.com/kkdai/youtube"
)

func main() {
	currentFile, _ := filepath.Abs(os.Args[0])
	log.Println("download to file=", currentFile)

	// NewYoutube(debug) if debug parameter will set true we can log of messages
	y := NewYoutube(true)
	y.DecodeURL("https://www.youtube.com/watch?v=rFejpH_tAHM")
	y.StartDownload(currentFile)
}
