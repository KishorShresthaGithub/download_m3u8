package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

// download m3u8/html files
func main() {
	fmt.Println("Start downloading m3u8")

	// input file
	filePath, skip := inputPrompt()

	// read file
	file := readFile(filePath)

	// get links
	links := prepareLinks(file)
	// get timestamps and filenames
	linkFilename := writeAriaLinks(links)

	// write to aria file
	files := preparePlaylist(file)
	// write to ffmpeg file
	playlistFilename := writePlaylist(files)

	// add -skip  flag
	if *skip != false {

		// run aria for the input file
		downloadUsingAria(linkFilename)
	}

	// check for png files
	checkIfPng(files)

	// run ffmpeg to merge files using the input
	mergePlaylist(playlistFilename)

	fmt.Printf("Downloaded completed for %v", file)

}

func inputPrompt() (*string, *bool) {

	file := flag.String("i", "", "Skip files")
	skip := flag.Bool("s", false, "Skip files")

	flag.Parse()

	return file, skip
}

func readFile(filePath *string) *[]string {
	fmt.Printf("Working on %v \n", *filePath)

	file, err := os.Open(*filePath)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	links := make([]string, 0)
	videoInfo := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()

		if match, _ := regexp.MatchString("(https?:\\/\\/).+", line); match == true {
			links = append(links, line)
		}

		videoInfo = append(videoInfo, line)

	}

	fmt.Println(links)

	arr := []string{"asd"}

	return &arr
}
func prepareLinks(file *[]string) *[]string {

	// check if has only one  link
	// use link to to create other links

	return nil

}

func preparePlaylist(file *[]string) *[]string {
	//store

	return nil
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
