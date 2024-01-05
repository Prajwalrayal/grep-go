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

var recursive bool = false
var caseInSensitive = false

func searchPattern_dir(directory_path, searchWord string, depth int, wg *sync.WaitGroup, resultChan chan<- []string) {
	var matchingLines []string
	defer wg.Done()
	err := filepath.WalkDir(directory_path, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if dir.IsDir() && strings.Count(path, string(os.PathSeparator)) > depth {
			return fs.SkipDir
		}

		if !dir.IsDir() {
			wg.Add(1)
			go func() {
				defer wg.Done()
				lines, err := searchPattern(path, searchWord)

				if err != nil {
					fmt.Println(err)
					return
				}
				resultChan <- lines

			}()
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	resultChan <- matchingLines
}

func searchPattern(filePath, searchWord string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var matchingLines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		searchWord_local := searchWord
		if caseInSensitive {
			line = strings.ToLower(line)
			searchWord_local = strings.ToLower(searchWord)
		}
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
	var path string
	var searchWord string
	flag.BoolVar(&recursive, "r", false, "Enable recursive search")
	flag.BoolVar(&caseInSensitive, "i", false, "Enable case insensitive search")
	flag.StringVar(&path, "path", "", "Specify the file path")
	flag.StringVar(&searchWord, "word", "", "Specify the word to search")

	flag.Parse()
	if path == "" || searchWord == "" {
		fmt.Println("Both file path and search word are required.")
		return
	}
	var matchingLines []string
	var err error
	resultChan := make(chan []string)
	var wg sync.WaitGroup

	fileInfo, err := os.Stat(path)

	if err != nil {
		fmt.Println(err)
		return
	}
	depth := 0
	if recursive {
		depth = math.MaxInt64
	}
	if fileInfo.IsDir() {
		wg.Add(1)
		go searchPattern_dir(path, searchWord, depth, &wg, resultChan)
	} else {
		matchingLines, err = searchPattern(path, searchWord)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for lines := range resultChan {
		matchingLines = append(matchingLines, lines...)
	}

	if len(matchingLines) > 0 {
		fmt.Printf("The word '%s' was found in the following lines :\n", searchWord)
		for _, line := range matchingLines {
			fmt.Println(line)
		}
	} else {
		fmt.Printf("The word '%s' was not found'\n", searchWord)
	}
}
