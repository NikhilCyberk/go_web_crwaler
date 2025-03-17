// File: indexer/indexer.go
package indexer

import (
	"fmt"

	"webcrawler/crawler"
)

// Global map to store word counts per URL
var wordIndex = make(map[string]map[string]int)

// Track visited pages
var visited = make(map[string]bool)

// Start begins the indexing process
func Start(wordCountChan <-chan crawler.WordCount, visitedChan <-chan string, doneChan <-chan bool) {
	for {
		select {
		case wc := <-wordCountChan:
			// Process word count
			if wc.Word != "" {
				if wordIndex[wc.Word] == nil {
					wordIndex[wc.Word] = make(map[string]int)
				}
				wordIndex[wc.Word][wc.URL] += wc.Count
			}
		case url := <-visitedChan:
			// Mark URL as visited
			visited[url] = true
		case <-doneChan:
			// Crawler has finished
			return
		}
	}
}

// SearchWord searches for a word in the index
func SearchWord(word string) {
	// Search for a word in the global map `wordIndex`
	fmt.Printf("\nResults for '%s':\n", word)

	if len(wordIndex[word]) == 0 {
		fmt.Printf("Word '%s' not found\n", word)
		return
	}

	var totalCount int
	urlCount := len(wordIndex[word])
	for url, count := range wordIndex[word] {
		totalCount += count
		fmt.Printf("%s: %d\n", url, count)
	}

	fmt.Printf("Total occurrences: %d\n", totalCount)
	fmt.Printf("Total URLs: %d\n", urlCount)
}
