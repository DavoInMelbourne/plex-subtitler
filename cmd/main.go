package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Println("Hello and welcome to the Plex Subtitler\n\n")
	baseDir := flag.String("start", "default starting directory", "a string for the starting directory")
	flag.Parse()
	fmt.Println("startDir:", *baseDir)

	ProcessDir(*baseDir, "")
	//ProcessDir("/Users/pauldavies/code/sandpit/plex-subtitler/testFiles", "")
}

func ProcessDir(startDir string, filename string) {
	fmt.Println("ProcessDir > start > startDir: ", startDir, "filename: ", filename)

	files, _ := os.ReadDir(startDir)

	for _, file := range files {
		fullFilename := fmt.Sprintf("%s/%s", startDir, file.Name())
		// fmt.Printf("%s is dir: %v\n", fullFilename, file.IsDir())

		// Process sub directory using recursion, or if it's a file, check if it's 	 show file and process accordingly
		if file.IsDir() == true {
			ProcessDir(fullFilename, file.Name())
			//Delete the Subs dir
			subsDir := fmt.Sprintf("%s/%s", fullFilename, "Subs")
			os.RemoveAll(subsDir)
		} else {
			ProcessFile(startDir, file.Name())
		}
	}

}

func ProcessFile(startDir string, filename string) {
	fmt.Println("ProcessFile > start > startDir: ", startDir, "filename: ", filename)
	showDir := ""
	subtitleBasename := ""
	subtitleFullFileName := ""
	subsDir := ""
	showDir = strings.Repeat(startDir, 1)

	// Check if the file ends in .mp4, in which case it could have subtitles
	if strings.HasSuffix(filename, ".mp4") {
		subtitleBasename = strings.TrimSuffix(filename, ".mp4")
		subtitleFullFileName = fmt.Sprintf("%s/%s", showDir, subtitleBasename)
		fmt.Println(filename, " is a show that might have subtitles with subtitle base filename: ", subtitleBasename)

		// Check if we have a subtitle file in the same directory. Lots of movies have these
		// Replace .mp4 with .srt
		srtFile := strings.TrimSuffix(filename, ".mp4") + ".srt"
		fullSrtFile := fmt.Sprintf("%s/%s", showDir, srtFile)
		enSrtFile := strings.TrimSuffix(filename, ".mp4") + ".en.srt"
		fullEnSrtFile := fmt.Sprintf("%s/%s", showDir, enSrtFile)

		// If the .en.srt file already exists, we don't want to process it again
		// Check if the .srt file exists
		_, err := os.Stat(fullEnSrtFile)
		if err == nil {
			fmt.Println(srtFile, "English Subtitie already processed.")
		} else {
			// Check if the .srt file exists
			subtitleFileInfo, err := os.Stat(fullSrtFile)
			if err == nil {
				fmt.Println(srtFile, "Subtitie exists in same directory.")

				if subtitleFileInfo.Size() > 3*1024 {
					fmt.Println(srtFile, "is greater than 3 KB.")

					e := os.Rename(fullSrtFile, fullEnSrtFile)
					if e != nil {
						log.Fatal(e)
					}
				}
			}

			// Check if there is a Subs directory to see if there are subtitles to process
			files, _ := os.ReadDir(startDir)
			for _, file := range files {
				if file.IsDir() == true && file.Name() == "Subs" {
					subsDir = fmt.Sprintf("%s/%s", startDir, file.Name())
					ProcessSubsDir(subsDir, subtitleFullFileName, subtitleBasename)

				}
			}
		}
	} else if strings.HasSuffix(filename, ".jpg") ||
		strings.HasSuffix(filename, ".txt") ||
		strings.HasSuffix(filename, ".nfo") ||
		strings.HasSuffix(filename, ".jpeg") ||
		strings.HasSuffix(filename, ".png") {
		// Delete the file
		fullFileName := fmt.Sprintf("%s/%s", showDir, filename)
		os.Remove(fullFileName)
	}
}

func ProcessSubsDir(startDir string, subtitleFilename string, subtitleBasename string) {
	//fmt.Println("ProcessSubsDir > start > startDir: ", startDir, "subtitleFilename: ", subtitleFilename, " subtitleBasename: ", subtitleBasename)

	srtFiles := []string{}

	// Check if we have directories or files in this subtitles folder
	files, _ := os.ReadDir(startDir)

	for _, file := range files {

		// If we have a directory in the Subs folder, we must drill into the one that matches our subtitle base name
		if file.IsDir() == true {
			newSubsDir := fmt.Sprintf("%s/%s", startDir, subtitleBasename)
			ProcessSubsDir(newSubsDir, subtitleFilename, subtitleBasename)
			break
		} else {
			if strings.Contains(file.Name(), "English") || strings.Contains(file.Name(), "en.srt") || strings.Contains(file.Name(), "eng.srt") {
				srtFiles = append(srtFiles, file.Name())
			}
		}
	}

	// Process the subtitle file(s) depending on the length of the array
	if len(srtFiles) == 0 {
		return
	}

	// Check if any of the subtitle files are forced
	for _, srt := range srtFiles {
		if strings.Contains(srt, "Forced") || strings.Contains(srt, "forced") {
			currentForcedSubtitleFilename := fmt.Sprintf("%s/%s", startDir, srt)
			newForcedSubtitleFilename := fmt.Sprintf("%s%s", subtitleFilename, ".en.forced.srt")

			fmt.Println("About to rename ", currentForcedSubtitleFilename, " to ", newForcedSubtitleFilename)

			e := os.Rename(currentForcedSubtitleFilename, newForcedSubtitleFilename)
			if e != nil {
				log.Fatal(e)
			}

			break
		}
	}

	// create a new srt file array with the unforced ones
	unforcedSrtFiles := []string{}
	for _, srt := range srtFiles {
		if !strings.Contains(srt, "Forced") && !strings.Contains(srt, "forced") {
			unforcedSrtFiles = append(unforcedSrtFiles, srt)
		}
	}

	// If unforced array is empty return
	if len(unforcedSrtFiles) == 0 {
		return
	}

	var largestFile string
	var largestSize int64 = 0

	// Loop through the filenames and find the largest file
	for _, filename := range unforcedSrtFiles {
		currentCheckFile := fmt.Sprintf("%s/%s", startDir, filename)
		fileInfo, err := os.Stat(currentCheckFile)
		if err != nil {
			fmt.Println("Error getting file info for", filename, ":", err)
			continue
		}

		// Compare the current file's size with the largest size found so far
		if fileInfo.Size() >= largestSize {
			largestSize = fileInfo.Size()
			largestFile = filename
		}
	}

	currenUnforcedSubtitleFilename := fmt.Sprintf("%s/%s", startDir, largestFile)
	newUnforcedSubtitleFilename := fmt.Sprintf("%s%s", subtitleFilename, ".en.srt")

	fmt.Println("About to rename ", currenUnforcedSubtitleFilename, " to ", newUnforcedSubtitleFilename)

	e := os.Rename(currenUnforcedSubtitleFilename, newUnforcedSubtitleFilename)
	if e != nil {
		log.Fatal(e)
	}
}
