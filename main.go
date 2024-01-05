package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// SearchResult represents the result of a search including the path and matching lines

type pair struct {
	LineNumber int
	Line       string
}
type SearchResult struct {
	Path          string
	MatchingLines []pair
}

// Global Flags
var recursive bool = false
var caseInSensitive = false
var path string
var searchWord string

// ANSI escape codes for color formatting
const (
	Reset = "\033[0m"
	Green = "\033[32m"
)

// Search in a Directory
func searchPatternDir(directoryPath, searchWord string, depth int, wg *sync.WaitGroup, resultChan chan<- SearchResult) {
	defer wg.Done()
	err := filepath.WalkDir(directoryPath, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// checking the depth to implement recursive/non-recursive search
		if dir.IsDir() && strings.Count(path, string(os.PathSeparator)) > depth {
			return fs.SkipDir
		}

		if !dir.IsDir() {
			// add a goroutine
			wg.Add(1)
			go func() {
				defer wg.Done()
				// concurrently opening files
				lines, err := searchPattern(path, searchWord)

				if err != nil {
					fmt.Println(err)
					return
				}
				// storing the result in channel
				if len(lines) != 0 {
					resultChan <- SearchResult{Path: path, MatchingLines: lines}
				}

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

func searchPattern(filePath, searchWord string) ([]pair, error) {
	// open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var matchingLines []pair

	// creating a file buffer
	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()

		// local copy of search word to handle case sensitive/insensitive search
		searchWordLocal := searchWord
		if caseInSensitive {
			line = strings.ToLower(line)
			searchWordLocal = strings.ToLower(searchWord)
		}

		// linear search for word in line
		if strings.Contains(line, searchWordLocal) {
			// Highlight the search word in the line
			regex := regexp.MustCompile(fmt.Sprintf(`%s`, searchWordLocal))
			line = regex.ReplaceAllStringFunc(line, func(match string) string {
				return Green + match + Reset
			})
			matchingLines = append(matchingLines, pair{LineNumber: lineNumber, Line: line})
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return matchingLines, nil
}

func main() {
	// Parsing the input flags
	flag.BoolVar(&recursive, "r", false, "Enable recursive search")
	flag.BoolVar(&caseInSensitive, "i", false, "Enable case-insensitive search")
	flag.StringVar(&path, "path", "", "Specify the file path")
	flag.StringVar(&searchWord, "word", "", "Specify the word to search")
	flag.Parse()

	// if either path or search word is missing
	if path == "" || searchWord == "" {
		fmt.Println("Both file path and search word are required.")
		fmt.Println("Check if path contains space then this won't work")
		return
	}
	var results []SearchResult
	var err error
	// create a channel and wait group
	resultChan := make(chan SearchResult)
	var wg sync.WaitGroup

	// get file info
	fileInfo, err := os.Stat(path)

	if err != nil {
		fmt.Println(err)
		return
	}
	// default depth is 0, if recursive then infinite depth
	depth := 0
	if recursive {
		depth = math.MaxInt64
	}
	if fileInfo.IsDir() {
		// add goroutine
		wg.Add(1)
		go searchPatternDir(path, searchWord, depth, &wg, resultChan)
	} else {
		lines, err := searchPattern(path, searchWord)
		if err != nil {
			return
		}
		results = append(results, SearchResult{Path: path, MatchingLines: lines})
	}
	// waiting for goroutines to end
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	// extracting results from channel
	for result := range resultChan {
		results = append(results, result)
	}

	if len(results) > 0 {
		fmt.Printf("Search results for the word '%s':\n", searchWord)
		for _, result := range results {
			fmt.Printf("File: %s\n", result.Path)
			for _, line := range result.MatchingLines {
				fmt.Printf("  %d %s\n", line.LineNumber, line.Line)
			}
		}
	} else {
		fmt.Printf("The word '%s' was not found.\n", searchWord)
	}
}
