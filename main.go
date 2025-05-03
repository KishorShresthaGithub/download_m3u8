package main

import (
	module "f/src"
	"fmt"
)

// download m3u8/html files
func main() {
	fmt.Println("Program to download m3u8")

	// input file
	filePath, skip := module.InputPrompt()

	// read file
	downloadLinks, fileLinks := module.ReadFile(filePath)

	workspace := module.FileNameWithoutExtension(*filePath)
	module.CreateRequiredFolders(workspace)

	// get links
	linksFilename := module.PrepareLinks(downloadLinks, workspace)

	// write to playlist
	module.PreparePlaylist(fileLinks, workspace)

	// add -skip  flag
	if *skip != true {
		// run aria for the input file
		module.DownloadUsingAria(linksFilename)
	}

	// check for png files
	module.ProcessFilesIfPng(workspace)

	module.CopyNeighborIfDoesntExist(fileLinks,workspace)

	// run ffmpeg to merge files using the input
	module.MergePlaylist(workspace)

	fmt.Printf("\nDownloaded completed for %v \n", module.GetBasename(*filePath))

}
