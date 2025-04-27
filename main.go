package main

import (
	"fmt"
	module "main/src"
)

// download m3u8/html files
func main() {
	fmt.Println("Program to download m3u8")

	// input file
	filePath, skip := module.InputPrompt()

	// read file
	downloadLinks, fileLinks := module.ReadFile(filePath)

	// get links
	linksFilename := module.PrepareLinks(downloadLinks, *filePath)

	// write to playlist
	module.PreparePlaylist(fileLinks, *filePath)

	// add -skip  flag
	if *skip != true {
		// run aria for the input file
		module.DownloadUsingAria(linksFilename)
	}

	// check for png files
	module.ProcessFilesIfPng(module.FileNameWithoutExtension(*filePath))

	// run ffmpeg to merge files using the input
	module.MergePlaylist(module.FileNameWithoutExtension(*filePath))

	fmt.Printf("\nDownloaded completed for %v \n", module.GetBasename(*filePath))

}
