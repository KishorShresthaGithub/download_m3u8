package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func isHTMLFile(filename string) bool {
	return filepath.Ext(filename) == ".html"
}

func isM3U8File(filename string) bool {
	return filepath.Ext(filename) == ".m3u8"
}

func run() {

	// Define the directory path (Change this as needed)
	dirPath := "./"

	// Open the directory
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	arrHtml := make([]string, 0)
	arrM3u8 := make([]string, 0)

	// Loop through and print file names
	for _, file := range files {
		if file.IsDir() { // Skip directories
			continue
		}
		fileName := file.Name()

		if isHTMLFile(fileName) {
			fmt.Println(file)
			arrHtml = append(arrHtml, fileName)
		}
		if isM3U8File(fileName) {
			arrM3u8 = append(arrM3u8, fileName)
		}
	}

	htmlChunks := chunkArray(arrHtml)
	m3u8Chunks := chunkArray(arrM3u8)

	runCopy("html", htmlChunks)
	runCopy("mu", m3u8Chunks)

}

func runCopy(typ string, chunks [][]string) {

	for index, chunk := range chunks {

		dirname := fmt.Sprintf("../%v%v", typ, index)
		err := os.MkdirAll(dirname, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating destination folder:", err)
			return
		}

		// Copy files
		for _, file := range chunk {
			srcPath := filepath.Join(".", file)
			dstPath := filepath.Join(dirname, file)

			script1 := filepath.Join("../", "download_m3u8.sh")
			script2 := filepath.Join("../", "play.sh")

			dstscript1 := filepath.Join(dirname, "download_m3u8.sh")
			dstscript2 := filepath.Join(dirname, "play.sh")

			err := copyFile(srcPath, dstPath)
			if err != nil {
				fmt.Println("Failed to copy:", file, err)
			} else {
				fmt.Println("Copied:", file)
			}
			err = copyFile(script1, dstscript1)

			if err != nil {
				fmt.Println("Failed to copy:", file, err)
			} else {
				fmt.Println("Copied:", file)
			}
			err = copyFile(script2, dstscript2)

			if err != nil {
				fmt.Println("Failed to copy:", file, err)
			} else {
				fmt.Println("Copied:", file)
			}
		}

	}
}

func chunkArray(arr []string) [][]string {
	chunks := make([][]string, 0)

	for i := 0; i < len(arr); i += 5 {
		temp := make([]string, 0)

		if i+5 < len(arr) {
			temp = arr[i : i+5]
		} else {
			temp = arr[i:]
		}

		chunks = append(chunks, temp)
	}
	return chunks
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}
