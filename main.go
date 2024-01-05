package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Global Flags
var recursive bool = false
var caseInSensitive = false
var path string
var searchWord string

// Search in a Directory
func searchPattern_dir(directory_path, searchWord string, depth int, wg *sync.WaitGroup, resultChan chan<- []string) {
	defer wg.Done()
	err := filepath.WalkDir(directory_path, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// checking the depth to implement recursive/non recursive search
		if dir.IsDir() && strings.Count(path, string(os.PathSeparator)) > depth {
			return fs.SkipDir
		}

		if !dir.IsDir() {
			// add a go routine
			wg.Add(1)
			go func() {
				defer wg.Done()
				// concurently opening files
				lines, err := searchPattern(path, searchWord)

				if err != nil {
					fmt.Println(err)
					return
				}
				// storing the result in channel
				resultChan <- lines

			}()
		}
		return nil
	})
	// handle error
	if err != nil {
		fmt.Println(err)
		return
	}

}

func searchPattern(filePath, searchWord string) ([]string, error) {
	// open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var matchingLines []string

	// creating a file buffer
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// local copy of search word to handle case sensetive/insensitive search
		searchWord_local := searchWord
		if caseInSensitive {
			line = strings.ToLower(line)
			searchWord_local = strings.ToLower(searchWord)
		}
		// linear search for word in line
		if strings.Contains(line, searchWord_local) {
			matchingLines = append(matchingLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return matchingLines, nil
}

func main() {

	//Parsing the input flags
	flag.BoolVar(&recursive, "r", false, "Enable recursive search")
	flag.BoolVar(&caseInSensitive, "i", false, "Enable case insensitive search")
	flag.StringVar(&path, "path", "", "Specify the file path")
	flag.StringVar(&searchWord, "word", "", "Specify the word to search")
	flag.Parse()

	// if either path or searchword is missing
	if path == "" || searchWord == "" {
		fmt.Println("Both file path and search word are required.")
		return
	}
	var matchingLines []string
	var err error
	// create a channel and wait group
	resultChan := make(chan []string)
	var wg sync.WaitGroup

	// get file info
	fileInfo, err := os.Stat(path)

	if err != nil {
		fmt.Println(err)
		return
	}
	// default depth is 0 , if recursive then infinite depth
	depth := 0
	if recursive {
		depth = math.MaxInt64
	}
	if fileInfo.IsDir() {
		// add go routine
		wg.Add(1)
		go searchPattern_dir(path, searchWord, depth, &wg, resultChan)
	} else {
		matchingLines, err = searchPattern(path, searchWord)
	}
	// waiting for goroutines to end
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	// extracting results from channel
	for lines := range resultChan {
		matchingLines = append(matchingLines, lines...)
	}
	// display result
	if len(matchingLines) > 0 {
		fmt.Printf("The word '%s' was found in the following lines :\n", searchWord)
		for _, line := range matchingLines {
			fmt.Println(line)
		}
	} else {
		fmt.Printf("The word '%s' was not found'\n", searchWord)
	}
}
