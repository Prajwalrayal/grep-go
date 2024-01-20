package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Search in a Directory
func SearchPatternDir(directoryPath, searchWord string, depth int, wg *sync.WaitGroup, resultChan chan<- SearchResult) {
	defer wg.Done()
	err := filepath.WalkDir(directoryPath, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// checking the depth to implement recursive/non-recursive search
		if dir.IsDir() && strings.Count(path, string(os.PathSeparator)) > depth {
			return fs.SkipDir
		}

		if !dir.IsDir() && !isExecutable(path) {
			// add a goroutine
			sem <- struct{}{}
			wg.Add(1)
			go func() {
				defer func() {
					<-sem
					wg.Done()
				}()
				// concurrently opening files
				// search.go
				lines, err := SearchPattern(path, searchWord)

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
