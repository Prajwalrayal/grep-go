package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
)

// SearchResult represents the result of a search including the path and matching lines

func init() {
	sem = make(chan struct{}, 10)
}

func main() {

	// flags.go
	Parse_Flags()


	// utilities.go
	if !Validate_input(path, searchWord) {
		return
	}

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
	// default depth is the current path depth, if recursive then infinite depth
	depth := strings.Count(path, string(os.PathSeparator))
	if recursive {
		depth = math.MaxInt64
	}
	if fileInfo.IsDir() {
		// add goroutine
		wg.Add(1)
		// directory_search.go
		go SearchPatternDir(path, searchWord, depth, &wg, resultChan)
	} else {
		// search.go
		lines, err := SearchPattern(path, searchWord)
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
	// utilities.go
	Display_Results(resultChan)
}
