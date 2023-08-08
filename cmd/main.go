package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
}

func ProcessDir(startDir string, filename string) {
	fmt.Println("ProcessDir > start > startDir: ", startDir, "filename: ", filename)

	files, _ := ioutil.ReadDir(startDir)

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
	//fmt.Println("ProcessFile > start > startDir: ", startDir, "filename: ", filename)
	showDir := ""
	subtitleBasename := ""
	subtitleFullFileName := ""
	subsDir := ""

	// Check if the file ends in .mp4, in which case it could have subtitles
	if strings.HasSuffix(filename, ".mp4") {
		showDir = strings.Repeat(startDir, 1)
		subtitleBasename = strings.TrimSuffix(filename, ".mp4")
		subtitleFullFileName = fmt.Sprintf("%s/%s", showDir, subtitleBasename)
		fmt.Println(filename, " is a show that might have subtitles with subtitle base filename: ", subtitleBasename)

		// Check if there is a Subs directory to see if there are subtitles to process
		files, _ := ioutil.ReadDir(startDir)
		for _, file := range files {
			if file.IsDir() == true && file.Name() == "Subs" {
				subsDir = fmt.Sprintf("%s/%s", startDir, file.Name())
				ProcessSubsDir(subsDir, subtitleFullFileName, subtitleBasename)

			}
		}
	}

}

func ProcessSubsDir(startDir string, subtitleFilename string, subtitleBasename string) {
	//fmt.Println("ProcessSubsDir > start > startDir: ", startDir, "subtitleFilename: ", subtitleFilename, " subtitleBasename: ", subtitleBasename)

	srtFiles := []string{}

	// Check if we have directories or files in this subtitles folder
	files, _ := ioutil.ReadDir(startDir)

	for _, file := range files {

		// If we have a directory in the Subs folder, we must drill into the one that matches our subtitle base name
		if file.IsDir() == true {
			newSubsDir := fmt.Sprintf("%s/%s", startDir, subtitleBasename)
			ProcessSubsDir(newSubsDir, subtitleFilename, subtitleBasename)
			break
		} else {
			if strings.Contains(file.Name(), "English") || strings.Contains(file.Name(), "en.srt") {
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

	index := len(unforcedSrtFiles) - 1
	currentForcedSubtitleFilename := fmt.Sprintf("%s/%s", startDir, unforcedSrtFiles[index])
	newForcedSubtitleFilename := fmt.Sprintf("%s%s", subtitleFilename, ".en.srt")

	fmt.Println("About to rename ", currentForcedSubtitleFilename, " to ", newForcedSubtitleFilename)

	e := os.Rename(currentForcedSubtitleFilename, newForcedSubtitleFilename)
	if e != nil {
		log.Fatal(e)
	}

}
