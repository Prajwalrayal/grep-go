package main

import "flag"

func Parse_Flags() {
	// Parsing the input flags
	flag.BoolVar(&recursive, "r", false, "Enable recursive search")
	flag.BoolVar(&caseInSensitive, "i", false, "Enable case-insensitive search")
	flag.StringVar(&path, "path", "", "Specify the file path")
	flag.StringVar(&searchWord, "word", "", "Specify the word to search")
	flag.Parse()
}
