package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// download m3u8/html files
func main() {
	fmt.Println("Start downloading m3u8")

	// input file
	filePath, skip := inputPrompt()

	// read file
	downloadLinks, fileLinks := readFile(filePath)

	// get links
	linksFilename := prepareLinks(downloadLinks, *filePath)

	// write to playlist
	files := preparePlaylist(fileLinks, *filePath)

	fmt.Printf(files)

	// add -skip  flag
	if *skip != false {
		// run aria for the input file
		downloadUsingAria(linksFilename)
	}

	// check for png files
	//checkIfPng()

	// run ffmpeg to merge files using the input
	// mergePlaylist()

	fmt.Printf("\nDownloaded completed for %v", getBasename(*filePath))

}

func inputPrompt() (*string, *bool) {

	file := flag.String("i", "", "Skip files")
	skip := flag.Bool("s", false, "Skip files")

	flag.Parse()

	return file, skip
}

func getBasename(str string) string {
	return filepath.Base(str)
}

func readFile(filePath *string) (*[]string, *[]string) {
	fmt.Printf("Working on %v \n", *filePath)

	file, err := os.Open(*filePath)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	links := make([]string, 0)
	videoInfo := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()

		if match, _ := regexp.MatchString("(https?:\\/\\/).+", line); match {
			links = append(links, line)
			videoInfo = append(videoInfo, fmt.Sprintf("file %v", getBasename(line)))
		} else {

			videoInfo = append(videoInfo, line)
		}
	}

	return &links, &videoInfo
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func prepareLinks(file *[]string, filename string) string {

	withoutExtension := fileNameWithoutExtension(filename)

	err := os.MkdirAll(withoutExtension, os.ModePerm)
	check(err)

	linkFile := fmt.Sprintf("%v/%v.links.txt", withoutExtension, withoutExtension)

	newFile, err := os.Create(linkFile)
	check(err)

	defer newFile.Close()
	newFile.Sync()

	// check if has only one  link

	writer := bufio.NewWriter(newFile)

	for _, v := range *file {
		_, err := writer.WriteString(fmt.Sprintf("%v \n\t out=\"%v/%v\" \n", v, withoutExtension, getBasename(v)))
		check(err)

	}
	writer.Flush()

	// use link to to create other links

	return linkFile

}

func preparePlaylist(file *[]string, filename string) string {

	fmt.Println(*file)

	withoutExtension := fileNameWithoutExtension(filename)

	err := os.MkdirAll(withoutExtension, os.ModePerm)
	check(err)

	linkFile := fmt.Sprintf("%v/%v.playlist.txt", withoutExtension, withoutExtension)

	newFile, err := os.Create(linkFile)
	check(err)
	defer newFile.Close()
	newFile.Sync()

	// check if has only one  link

	writer := bufio.NewWriter(newFile)

	for _, v := range *file {

		line := v + "\n"

		_, err := writer.WriteString(line)
		check(err)

	}
	writer.Flush()

	// use link to to create other links

	return linkFile
}

func checkIfPng(files *[]string) bool {
	// check if header is png
	//

	return false
}
func decrypt(files *[]string) {

}

func downloadUsingAria(inputFilename string) {

}

func joinPlaylist() {

}

func check(err error) {

	if err != nil {

		log.Fatal(err)
		panic(err)

	}

}

func writeAriaLinks(links *[]string) string {

	return ""
}

func writePlaylist(file *[]string) string {

	return ""
}

func mergePlaylist(fileName string) string {

	return ""
}
