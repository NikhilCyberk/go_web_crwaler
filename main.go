// File: main.go
package main

import (
	"fmt"
	"strings"
	"webcrawler/crawler"
	"webcrawler/indexer"
)

func display() {
	fmt.Println("\nWelcome to the Web Crawler!")
	for {
		fmt.Println("\nOptions:")
		fmt.Println("1. Search for a word")
		fmt.Println("2. Exit")
		fmt.Print("\nEnter your choice (1-2): ")

		var choice string
		fmt.Scanln(&choice)                // Read user's choice
		choice = strings.TrimSpace(choice) // Trim any extra spaces

		switch choice {
		case "1":
			fmt.Print("Enter search term: ")
			var word_to_find string
			fmt.Scanln(&word_to_find)
			term := strings.TrimSpace(word_to_find) // Trim any extra spaces
			indexer.SearchWord(term)

		case "2":
			fmt.Println("Exiting...")
			return // Exit the function and stop the loop

		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func main() {
	wordCountChan := make(chan crawler.WordCount) // Channel for word count updates
	visitedChan := make(chan string)              // Channel for visited URLs
	doneChan := make(chan bool)                   // Channel to signal completion

	// Start the indexer
	go indexer.Start(wordCountChan, visitedChan, doneChan)

	// Start the web crawler
	go crawler.Start(crawler.AllowedDomain, wordCountChan, visitedChan, doneChan)

	// Display the menu
	display()
}
	