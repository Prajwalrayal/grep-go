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
)

func searchPattern_dir(directory_path, searchWord string, depth int) ([]string, error) {
	var matchingLines []string

	err := filepath.WalkDir(directory_path, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if dir.IsDir() && strings.Count(path, string(os.PathSeparator)) > depth {
			return fs.SkipDir
		}

		if !dir.IsDir() {
			lines, err := searchPattern(path, searchWord)
			if err != nil {
				return err
			}
			matchingLines = append(matchingLines, lines...)
		}
		return nil
	})

	return matchingLines, err
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
		if strings.Contains(line, searchWord) {
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
	var recursive bool

	flag.BoolVar(&recursive, "r", false, "Enable recursive search")
	flag.StringVar(&path, "path", "", "Specify the file path")
	flag.StringVar(&searchWord, "word", "", "Specify the word to search")

	flag.Parse()
	if path == "" || searchWord == "" {
		fmt.Println("Both file path and search word are required.")
		return
	}

	var matchingLines []string
	var err error

	fileInfo, err := os.Stat(path)

	if err != nil {
		panic(err)
	}
	depth := 0
	if recursive {
		depth = math.MaxInt64
	}
	if fileInfo.IsDir() {
		matchingLines, err = searchPattern_dir(path, searchWord, depth)
	} else {
		matchingLines, err = searchPattern(path, searchWord)
	}

	if err != nil {
		panic(err)
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
