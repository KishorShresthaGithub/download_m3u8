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
	"slices"
	"strings"
)

func InputPrompt() (*string, *bool) {

	file := flag.String("i", "", "Skip files")
	skip := flag.Bool("s", false, "Skip files")

	flag.Parse()

	return file, skip
}

func GetBasename(str string) string {
	u, err := url.Parse(str)
	Check(err)

	u.RawQuery = ""
	u.Fragment = ""

	return filepath.Base(u.String())
}

func CheckIfLineIsLink(text string) bool {
	match, _ := regexp.MatchString("^(https?:\\/\\/).+", text)
	return match
}

func CheckIfLineStartsWithHash(text string) bool {
	match, _ := regexp.MatchString("^#.+$", text)
	return match
}
func GetBaseUrl(in string) string {

	u, err := url.Parse(in)
	Check(err)

	dirPath := path.Dir(u.Path)

	return fmt.Sprintf(
		"%s://%s%s",
		u.Scheme,
		u.Host,
		dirPath,
	)

}

func ReadFile(filePath *string) (*[]string, *[]string) {
	fmt.Printf("=================`Working on %v +++++ Reading \n", *filePath)

	file, err := os.Open(*filePath)
	Check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	links := make([]string, 0)
	videoInfo := make([]string, 0)

	firstLinkBaseUrl := ""
	firstLinkBasename := ""
	fileNameOnlyMode := false

	for scanner.Scan() {
		line := scanner.Text()
		basename := GetBasename(line)

		if match, _ := regexp.MatchString("^<.+", line); match {
			continue
		}

		if CheckIfLineStartsWithHash(line) {
			videoInfo = append(videoInfo, line)
		} else if CheckIfLineIsLink(line) {
			links = append(links, line)
			videoInfo = append(videoInfo, fmt.Sprintf("file %v", basename))

			if firstLinkBaseUrl == "" {
				firstLinkBaseUrl = GetBaseUrl(line)
				firstLinkBasename = GetBasename(line)
			}

		} else {
			fileNameOnlyMode = true
			videoInfo = append(videoInfo, fmt.Sprintf("file %v", basename))
			links = append(links, fmt.Sprintf("%v/%v", firstLinkBaseUrl, line))

		}

	}

	if fileNameOnlyMode {
		name := fmt.Sprintf("file %v", firstLinkBasename)

		index := slices.Index(videoInfo, name)

		fmt.Println("++++++++++++++++++++++++++++++++++++++", index)

		if index > -1 {
			videoInfo = slices.Delete(videoInfo, index, index+1)
		}

	}

	return &links, &videoInfo
}

func FileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func CreateRequiredFolders(workspace string) {

	fmt.Println("=================== Creating required folders")
	err := os.MkdirAll(workspace, os.ModePerm)
	Check(err)

	err = os.MkdirAll(filepath.Join(workspace, "parts"), os.ModePerm)
	Check(err)

}

func PrepareLinks(file *[]string, workspace string) string {
	fmt.Println("====================== Preparing Links")

	linkFile := fmt.Sprintf("%v/%v.links.txt", workspace, workspace)

	newFile, err := os.Create(linkFile)
	Check(err)

	defer newFile.Close()
	newFile.Sync()

	writer := bufio.NewWriter(newFile)

	for _, url := range *file {
		_, err := writer.WriteString(fmt.Sprintf("%v \n\t dir=%v \n\t out=parts/%v \n", url, workspace, GetBasename(url)))

		Check(err)

	}
	writer.Flush()

	return linkFile

}

func PreparePlaylist(file *[]string, workspace string) string {

	fmt.Println("====================== Preparing Playlists")

	linkFile := fmt.Sprintf("%v/parts/%v.playlist.txt", workspace, workspace)

	newFile, err := os.Create(linkFile)
	Check(err)
	defer newFile.Close()
	newFile.Sync()

	writer := bufio.NewWriter(newFile)

	for _, v := range *file {

		line := v + "\n"

		_, err := writer.WriteString(line)
		Check(err)

	}
	writer.Flush()
	return linkFile
}

func CheckifPng(filename string) bool {

	f, err := os.Open(filename)
	Check(err)
	defer f.Close()

	sig := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	buf := make([]byte, 8)
	if _, err := f.Read(buf); err != nil {
		return false
	}
	for i := range sig {
		if buf[i] != sig[i] {
			return false
		}
	}
	return true
}

// Helper to copy a file
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	Check(err)
	defer in.Close()

	out, err := os.Create(dst)
	Check(err)
	defer out.Close()

	_, err = io.Copy(out, in)
	Check(err)
	return nil
}

func ProcessFilesIfPng(workspace string) {
	fmt.Println("======================= Process if pngg")
	// Your main logic
	parts := filepath.Join(workspace, "parts")
	backup := filepath.Join(workspace, "backup")

	err := os.MkdirAll(backup, os.ModePerm)
	Check(err)

	entries, err := os.ReadDir(parts)
	Check(err)

	// Filter to files only (not directories)
	var files []os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry)
		}
	}
	if len(files) == 0 {
		fmt.Println("No files found in directory.")
		return
	}

	if check := CheckifPng(filepath.Join(parts, files[0].Name())); !check {
		return
	}
	fmt.Println("....................................................")

	for _, srcEntry := range files {
		srcPath := filepath.Join(parts, srcEntry.Name())

		if strings.Contains(srcEntry.Name(), "playlist.txt") {
			continue
		}

		// 1. Copy original to backup
		backupPath := filepath.Join(backup, srcEntry.Name())
		Check(copyFile(srcPath, backupPath))

		// 2. Trim the first 8 bytes of the original file
		in, err := os.Open(srcPath)
		Check(err)
		defer in.Close()

		_, err = in.Seek(8, io.SeekStart)
		Check(err)

		tmpPath := srcPath + ".tmp"
		out, err := os.Create(tmpPath)
		Check(err)
		defer out.Close()

		_, err = io.Copy(out, in)
		Check(err)

		// Close files before renaming
		in.Close()
		out.Close()

		// Replace original with trimmed version
		Check(os.Rename(tmpPath, srcPath))
	}
}
func Decrypt(files *[]string) {

}

func DownloadUsingAria(inputFilename string) {

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

func JoinPlaylist() {

}

func Check(err error) {

	if err != nil {
		panic(err)
	}

}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func CopyNeighborIfDoesntExist(fileLinks *[]string, workspace string) {

	for k, item := range *fileLinks {

		if CheckIfLineStartsWithHash(item) {
			continue
		}

		filename := strings.TrimPrefix(item, "file ")
		path := filepath.Join(workspace, "parts", filename)

		if !fileExists(path) {
			fmt.Printf("%v doesnt exist \n", path)

			var prevItem string
			if k-1 < 0 {
				prevItem = (*fileLinks)[k+1]
			} else {
				prevItem = (*fileLinks)[k-1]
			}
			prevFilename := strings.TrimPrefix(prevItem, "file ")
			prevPath := filepath.Join(workspace, "parts", prevFilename)
			copyFile(prevPath, path)
		}
	}

}

func MergePlaylist(filename string) {

	playlist := filepath.Join(filename, "parts", fmt.Sprintf("%v.playlist.txt", filename))

	args := []string{
		"-f", "concat",
		"-safe", "0",
		"-i", playlist,
		"-fflags", "+genpts",
		"-r", "30",
		"-c:v", "copy",
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
