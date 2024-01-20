package main

import (
	"fmt"
	"os"
)

func isExecutable(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	// Check if the file has executable permission
	return fileInfo.Mode().IsRegular() && (fileInfo.Mode()&0111) != 0
}

func Validate_input(path string, searchWord string) bool {
	// if either path or search word is missing
	if path == "" || searchWord == "" {
		fmt.Println("Both file path and search word are required.")
		fmt.Println("Check if path contains space then this won't work")
		return false
	}
	return true
}
func Display_Results(resultChan <-chan SearchResult) {
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
