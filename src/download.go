package src

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func RunProgram() {
	// read file and get contents
	filePath, skipFile, key, iv, _ := InputPromt()

	// create all required folders
	CreateRequiredFolders(*filePath)

	// extracts the links from main file
	//	change the file name to index of link
	//      create download input file
	//      create ffmpeg playlist
	// if there is only one link append the first link to the other links
	links, playlist := ReadFile(*filePath)

	// copy text files to required folder
	linksFile := WriteLinksToFile(*filePath, &links)
	// WritePlaylistToFile(*filePath, &playlist)

	// skip download if there is skip flag
	if !*skipFile {
		// download the files with new names
		DownloadFiles(linksFile)
	}

	// copy to backup folder
	BackupDownloadFiles(*filePath)

	// manually check if key is defined
	if *key != "" {
		DecryptFiles(*filePath, *key, *iv)
	}

	// ~~remove file type identifier bits if key is not defined~~
	// loop through the playlist and check if files are ffmpeg processable
	// remove the unncessary lines and delete the unncessary files
	CheckIfFilesCorrupted(*filePath, &playlist)

	WritePlaylistToFile(*filePath, &playlist)

	// merge using ffmpeg

	MergePlaylist(*filePath)

	fmt.Println("==================== Download completed ======================")
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func DecryptFiles(filePath string, key string, iv string) {
}

func GetBasename(str string) string {
	u, err := url.Parse(str)
	Check(err)

	u.RawQuery = ""
	u.Fragment = ""

	basename := filepath.Base(u.String())

	ext := filepath.Ext(basename)

	basename = strings.TrimSuffix(basename, ext)
	return basename
}

func GetBaseUrl(in string) string {
	u, err := url.Parse(in)
	Check(err)

	dirPath := path.Dir(u.EscapedPath())

	return fmt.Sprintf(
		"%s://%s%s",
		u.Scheme,
		u.Host,
		dirPath,
	)
}

func ReadFile(filePath string) ([]string, []string) {
	fmt.Printf("========================================== Working in %v Reading\n", filePath)

	fileName := GetBasename(filePath)

	fmt.Println(filePath, fileName)

	file, err := os.Open(filePath)
	Check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	links := make([]string, 0)
	playlist := make([]string, 0)

	var firstLink string
	var firstLinkBase string

	counter := 1

	for scanner.Scan() {
		line := scanner.Text()

		if line != "" && string(line[0]) == "<" {
			continue
		}
		if line != "" && string(line[0]) == "#" {
			continue
		}

		if match, _ := regexp.MatchString(`^(https?://).+`, line); match {

			fmt.Println(line, firstLink)

			if firstLink == "" {
				firstLink = line
				firstLinkBase = GetBaseUrl(line)
			}

			links = append(links, line)
			playlist = append(playlist, fmt.Sprintf("file %v_%v", fileName, counter))
			counter++
		} else {
			link := fmt.Sprintf("%v/%v", firstLinkBase, line)

			links = append(links, link)
			playlist = append(playlist, fmt.Sprintf("file %v_%v", fileName, counter))
			counter++
		}

	}

	if firstLink == "" {
		panic("No links found")
	}

	return links, playlist
}

func CreateRequiredFolders(filePath string) {
	workspace := GetBasename(filePath)

	fmt.Println("=================== Creating required folders")
	err := os.MkdirAll(workspace, os.ModePerm)
	Check(err)

	err = os.MkdirAll(filepath.Join(workspace, "parts"), os.ModePerm)
	Check(err)

	// err = os.MkdirAll(filepath.Join(workspace, "backup"), os.ModePerm)
	// Check(err)
}

func InputPromt() (*string, *bool, *string, *string, *bool) {
	filePath := flag.String("i", "", "Input File")
	skip := flag.Bool("s", false, "Skip Files")
	keyFile := flag.String("k", "", "KeyFile")
	urlEncode := flag.Bool("u", false, "Url Encode")

	var ivFile *string

	if *keyFile != "" {
		ivFile = flag.String("iv", "", "IV string")

		if *ivFile == "" {
			panic("Iv string not defined")
		}
	}

	flag.Parse()

	return filePath, skip, keyFile, ivFile, urlEncode
}

func WriteLinksToFile(filePath string, links *[]string) string {
	fileName := GetBasename(filePath)
	fmt.Println("================================Preparing links", fileName)

	linkFileName := fmt.Sprintf("%v/%v.links.txt", fileName, fileName)

	newFile, err := os.Create(linkFileName)
	Check(err)
	defer newFile.Close()

	writer := bufio.NewWriter(newFile)

	for index, url := range *links {

		fragmentName := fmt.Sprintf("%v_%v", fileName, index+1)

		_, err := fmt.Fprintf(writer, "%v \n\t dir=%v \n\t out=parts/%v \n", url, fileName, fragmentName)

		Check(err)
	}
	writer.Flush()

	return linkFileName
}

func WritePlaylistToFile(filePath string, playlist *[]string) string {
	fileName := GetBasename(filePath)
	fmt.Println("================================Preparing playlist", fileName)

	linkFileName := fmt.Sprintf("%v/%v.playlist.txt", fileName, fileName)

	newFile, err := os.Create(linkFileName)
	Check(err)
	defer newFile.Close()

	writer := bufio.NewWriter(newFile)

	for _, v := range *playlist {

		line := v + "\n"

		_, err := writer.WriteString(line)

		Check(err)
	}
	writer.Flush()

	CopyFile(linkFileName, fmt.Sprintf("%v/parts/%v.playlist.txt", fileName, fileName))

	return linkFileName
}

func DownloadFiles(inputFilename string) {
	args := []string{
		"-U", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:134.0) Gecko/20100101 Firefox/134.0",
		"--header=Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"--header=Accept-Language: en-US,en;q=0.5",
		"--header=Accept-Encoding: gzip, deflate, br",
		"--header=Connection: keep-alive",
		"--header=Upgrade-Insecure-Requests: 1",
		"--header=Sec-Fetch-Dest: document",
		"--header=Sec-Fetch-Mode: navigate",
		"--header=Sec-Fetch-Site: none",
		"--header=Sec-Fetch-User: ?1",
		"--file-allocation=prealloc",
		"-j", "10", "-s", "10", "-x", "10", "-c",
		"-i", inputFilename,
	}

	cmd := exec.Command("aria2c", args...)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()

	fmt.Println("======================Starting download===================")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

// CopyFile copies a single file
func CopyFile(srcFile, dstFile string) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

// CopyDir copies a whole directory recursively
func CopyDir(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return CopyFile(path, dstPath)
	})
}

func BackupDownloadFiles(filePath string) {
	fileName := GetBasename(filePath)

	partsDir := fmt.Sprintf("%v/parts", fileName)
	backupDir := fmt.Sprintf("%v/backup", fileName)

	err := CopyDir(partsDir, backupDir)

	Check(err)
}

func RemoveFirstNBytesInPlace(filename string, n int) error {
	tempFile := filename + ".tmp"

	// Open original file
	in, err := os.Open(filename)
	Check(err)
	defer in.Close()

	// Create temporary file
	out, err := os.Create(tempFile)
	Check(err)
	defer out.Close()

	// Skip first N bytes
	_, err = in.Seek(int64(n), io.SeekStart)
	Check(err)

	// Copy rest to temp file
	_, err = io.Copy(out, in)
	Check(err)

	in.Close()
	out.Close()

	// Replace original file with temp file
	err = os.Rename(tempFile, filename)
	return err
}

func CheckIfFilesCorrupted(filePath string, playlist *[]string) {
	fileName := GetBasename(filePath)

	var filteredPlaylist []string

	for _, value := range *playlist {

		parts := strings.TrimPrefix(value, "file ")
		toCheck := fmt.Sprintf("%v/parts/%v", fileName, parts)

		args := []string{
			"-v", "error",
			"-i", toCheck,
			"-f", "null", "-",
		}

		cmd := exec.Command("ffmpeg", args...)
		cmd.Stdout, _ = os.Open(os.DevNull)
		cmd.Stderr, _ = os.Open(os.DevNull)

		fmt.Println("======================Checking ===================")

		if err := cmd.Run(); err != nil {
			fmt.Printf("%v is corrupted. Attempting to delete...\n", toCheck)

			if _, statErr := os.Stat(toCheck); statErr != nil {
				if os.IsNotExist(statErr) {
					fmt.Printf("File %v does not exist, nothing to remove.\n", toCheck)
				} else {
					fmt.Printf("Error checking file %v: %v\n", toCheck, statErr)
				}
			} else if removeErr := os.Remove(toCheck); removeErr != nil {
				fmt.Printf("Failed to remove %v: %v\n", toCheck, removeErr)
			} else {
				fmt.Printf("Deleted corrupted file: %v\n", toCheck)
			}

			continue
		}
		filteredPlaylist = append(filteredPlaylist, value)
	}

	*playlist = filteredPlaylist
}

func MergePlaylist(filePath string) {
	filename := GetBasename(filePath)

	playlist := filepath.Join(filename, "parts", fmt.Sprintf("%v.playlist.txt", filename))

	args := []string{
		"-f", "concat",
		"-safe", "0",
		"-hwaccel", "cuda",
		"-i", playlist,
		"-fflags", "+genpts",
 		"-preset", "p1", 
		"-c:v", "h264_nvenc",
		"-c:a", "copy",
		fmt.Sprintf("%v.mp4", filename),
	}

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()

	fmt.Println("======================Merging files===================")

	err := cmd.Run()
	Check(err)
}

// func ProcessFilesIfPng(workspace string) {
// 	fmt.Println("======================= Process if pngg")
// 	// Your main logic
// 	parts := filepath.Join(workspace, "parts")
// 	backup := filepath.Join(workspace, "backup")
//
// 	err := os.MkdirAll(backup, os.ModePerm)
// 	Check(err)
//
// 	entries, err := os.ReadDir(parts)
// 	Check(err)
//
// 	// Filter to files only (not directories)
// 	var files []os.DirEntry
// 	for _, entry := range entries {
// 		if !entry.IsDir() {
// 			files = append(files, entry)
// 		}
// 	}
// 	if len(files) == 0 {
// 		fmt.Println("No files found in directory.")
// 		return
// 	}
//
// 	if check := CheckifPng(filepath.Join(parts, files[0].Name())); !check {
// 		return
// 	}
// 	fmt.Println("....................................................")
//
// 	for _, srcEntry := range files {
// 		srcPath := filepath.Join(parts, srcEntry.Name())
//
// 		if strings.Contains(srcEntry.Name(), "playlist.txt") {
// 			continue
// 		}
//
// 		// 1. Copy original to backup
// 		backupPath := filepath.Join(backup, srcEntry.Name())
// 		Check(copyFile(srcPath, backupPath))
//
// 		// 2. Trim the first 8 bytes of the original file
// 		in, err := os.Open(srcPath)
// 		Check(err)
// 		defer in.Close()
//
// 		_, err = in.Seek(8, io.SeekStart)
// 		Check(err)
//
// 		tmpPath := srcPath + ".tmp"
// 		out, err := os.Create(tmpPath)
// 		Check(err)
// 		defer out.Close()
//
// 		_, err = io.Copy(out, in)
// 		Check(err)
//
// 		// Close files before renaming
// 		in.Close()
// 		out.Close()
//
// 		// Replace original with trimmed version
// 		Check(os.Rename(tmpPath, srcPath))
// 	}
// }
