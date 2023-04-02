package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/bogem/id3v2/v2"
)

const chapterOffset = -4

func getChapter(fileName string, tagTitle string) (string, string, int, int, error) {
	filenameRe := regexp.MustCompile(`(.*) \(Chapter (\d+)\).mp3`)
	match := filenameRe.FindStringSubmatch(fileName)

	if len(match) < 2 {
		return "", "", 0, 0, fmt.Errorf("filename did not match expected pattern")
	}
	filenameChapterNumber, err := strconv.Atoi(match[2])
	filenamePrefix := match[1]
	if err != nil {
		return "", "", 0, 0, err
	}

	titleRe := regexp.MustCompile(`(.*) - (\d+)`)
	match = titleRe.FindStringSubmatch(tagTitle)

	if len(match) < 2 {
		return "", "", 0, 0, fmt.Errorf("filename did not match expected pattern")
	}
	titleChapterNumber, err := strconv.Atoi(match[2])
	titlePrefix := match[1]

	if err != nil {
		return "", "", 0, 0, err
	}

	if filenameChapterNumber != titleChapterNumber {
		return "", "", 0, 0, fmt.Errorf("filename chapter %d does not match title chapter %d", filenameChapterNumber, titleChapterNumber)
	}

	fixedChapterNum := filenameChapterNumber + chapterOffset
	newFileName := fmt.Sprintf("%s (Chapter %d).mp3", filenamePrefix, fixedChapterNum)
	newTitle := fmt.Sprintf("%s - Chapter %d", titlePrefix, fixedChapterNum)

	return newFileName, newTitle, filenameChapterNumber, filenameChapterNumber + chapterOffset, nil
}

// In my case, each file is a book chapter and has incorrect chapter numbers
// Each chapter is 4 too high, eg chapter 10 is really chapter 6.
// Two things need to be corrected
// 1. Each filename is "<bookname> (Chapter N)", fix to "<bookname> (Chapter N-4)"
// 2. Each id3v2 'Title' field is "<bookname> - N", fix to "<bookname> - Ch. N-4"
func processFile(absFilePath string) {
	log.Printf("Processing file %q\n", absFilePath)

	tag, err := id3v2.Open(absFilePath, id3v2.Options{Parse: true})
	if err != nil {
		log.Println("Error reading id3v2 from file", absFilePath, err)
		return
	}

	baseFileName := filepath.Base(absFilePath)
	baseDir := filepath.Dir(absFilePath)

	newFileName, newTitle, _, _, err := getChapter(baseFileName, tag.Title())
	if err != nil {
		log.Printf("Error getting fixed chapter info from %q: %v\n", baseFileName, err)
		return
	}
	newAbsFilePath := filepath.Join(baseDir, newFileName)
	log.Printf("New filename: %q New Title: %q", newAbsFilePath, newTitle)
	log.Printf("Existing IDv3: %q %q %q", tag.Title(), tag.Artist(), tag.Year())

	// Step one, save new file metadata in-place
	// (idv3 lib doesn't support saving to new filename)
	tag.SetTitle(newTitle)
	err = tag.Save()
	if err != nil {
		log.Printf("Error while saving new tag metadata to file %v\n", err)
		return
	}
	// rename old file to new filename
	err = os.Rename(absFilePath, newAbsFilePath)
	if err != nil {
		log.Printf("Error while renaming file from %q to %q: %v\n", absFilePath, newAbsFilePath, err)
	}
}

func main() {
	for _, arg := range os.Args[1:] {
		var basepath string
		basepath, pattern := doublestar.SplitPattern(arg)
		if basepath == "." {
			// no pattern found, this is just a filepath
			absFilePath, err := filepath.Abs(arg)
			if err != nil {
				log.Printf("Error while resolving argument %q to an absolute path. Skipping...\n", arg)
				continue
			}
			processFile(absFilePath)
		} else {
			fsys := os.DirFS(basepath)
			matches, err := doublestar.Glob(fsys, pattern)
			fmt.Printf("Arg %q is a glob, found %d matches, basepath is %q pattern is %q\n", arg, len(matches), basepath, pattern)
			if err != nil {
				log.Println("Err while globbing", err)
				continue
			}

			for _, match := range matches {
				absFilePath := filepath.Join(basepath, match)
				processFile(absFilePath)
			}
		}
	}
}
