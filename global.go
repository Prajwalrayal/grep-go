package main

// Global Flags
var recursive bool = false
var caseInSensitive = false
var path string
var searchWord string
var results []SearchResult

// semaphore
var sem chan struct{}

// ANSI escape codes for color formatting
const (
	Reset = "\033[0m"
	Green = "\033[32m"
)
