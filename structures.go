package main

type pair struct {
	LineNumber int
	Line       string
}
type SearchResult struct {
	Path          string
	MatchingLines []pair
}
