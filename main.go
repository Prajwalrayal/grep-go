package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func searchWordInFile(filePath, searchWord string) ([]string, error) {
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
	var filePath string
	var searchWord string

	flag.StringVar(&filePath, "file", "", "Specify the file path")
	flag.StringVar(&searchWord, "word", "", "Specify the word to search")

	flag.Parse()

	if filePath == "" || searchWord == "" {
		fmt.Println("Both file path and search word are required.")
		return
	}

	matchingLines, err := searchWordInFile(filePath, searchWord)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(matchingLines) > 0 {
		fmt.Printf("The word '%s' was found in the following lines of the file '%s':\n", searchWord, filePath)
		for _, line := range matchingLines {
			fmt.Println(line)
		}
	} else {
		fmt.Printf("The word '%s' was not found in the file '%s'\n", searchWord, filePath)
	}
}
