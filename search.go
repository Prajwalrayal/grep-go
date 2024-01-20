package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func SearchPattern(filePath, searchWord string) ([]pair, error) {
	// open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var matchingLines []pair

	// creating a file buffer
	scanner := bufio.NewScanner(file)
	const maxBufferSize = 1024 * 1024	// 10 MB buffer size
	buf := make([]byte, maxBufferSize)
	scanner.Buffer(buf, maxBufferSize)
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
