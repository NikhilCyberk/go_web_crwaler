// File: crawler/crawler.go
package crawler

import (
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const AllowedDomain string = "https://usf-cs272-s25.github.io/top10/"
const BaseUrl string = "https://usf-cs272-s25.github.io"

type WordCount struct {
	Word  string
	URL   string
	Count int
}

// Start initiates the web crawler
func Start(startURL string, wordCountChan chan<- WordCount, visitedChan chan<- string, doneChan chan<- bool) {
	// Mark the URL as visited
	visitedChan <- startURL

	crawl(startURL, wordCountChan, visitedChan)

	// Signal completion
	doneChan <- true
}

// crawl processes a URL and its linked pages
func crawl(url string, wordCountChan chan<- WordCount, visitedChan chan<- string) {
	htmlData, err := fetchData(url)
	if err != nil {
		log.Println("Error fetching:", err)
		return
	}

	text, links := extractLinksAndText(htmlData)
	queue := completeUrl(links)

	// Process the text and send word counts
	processText(text, url, wordCountChan)

	// Process linked pages (only if they are under the BaseUrl domain)
	for _, link := range queue {
		if isAllowed(link) {
			// Check if the link has already been visited
			visitedChan <- link // This will be processed by the indexer

			// We'll receive a response that tells us if the URL was already visited
			crawl(link, wordCountChan, visitedChan)
		}
	}
}

// fetchData retrieves the HTML content from a URL
func fetchData(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// extractLinksAndText parses HTML to extract text and links
func extractLinksAndText(html string) (string, []string) {
	// Remove <style>...</style> blocks
	html = regexp.MustCompile(`(?is)<style.*?>.*?</style>`).ReplaceAllString(html, "")
	// Remove inline styles
	html = regexp.MustCompile(`(?i)style=["'][^"']*["']`).ReplaceAllString(html, "")
	// Extract <a href="..."> links
	linkRegex := regexp.MustCompile(`(?i)<a\s+[^>]*href=["']([^"']+)["']`)
	// Remove HTML tags
	text := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(html, " ")

	// Extract links
	var links []string
	matches := linkRegex.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) > 1 {
			link := match[1]
			links = append(links, link)
		}
	}

	// Normalize text (remove extra spaces)
	text = strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(text, " "))
	return text, links
}

// completeUrl converts relative URLs to absolute URLs
func completeUrl(url []string) []string {
	for i := 0; i < len(url); i++ {
		if !strings.HasPrefix(url[i], "http") { // Convert relative URLs
			url[i] = BaseUrl + "/" + strings.TrimPrefix(url[i], "/")
		}
	}
	return url
}

// isAllowed checks if a URL is in the allowed domain
func isAllowed(url string) bool {
	// Ensure the URL is under the BaseUrl domain
	return strings.HasPrefix(url, AllowedDomain)
}

// processText splits text into words and sends word counts
func processText(text string, url string, wordCountChan chan<- WordCount) {
	words := strings.Fields(text)
	for _, word := range words {
		word = strings.ToLower(word)
		word = cleanWord(word)
		if word != "" {
			wordCountChan <- WordCount{Word: word, URL: url, Count: 1}
		}
	}
}

// cleanWord removes punctuation and normalizes words
func cleanWord(word string) string {
	// Remove punctuation at beginning and end of words
	word = strings.Trim(word, ".,;-:!?\"'()[]{}<>—")
	// Remove non-alphanumeric characters from the beginning and end
	word = regexp.MustCompile(`^[^a-zA-Z0-9-]+|[^a-zA-Z0-9-]+$`).ReplaceAllString(word, "")
	// Replace em dashes, en dashes, and other special punctuation with a space or hyphen
	word = regexp.MustCompile(`[—–]`).ReplaceAllString(word, " ")
	word = regexp.MustCompile(`[-]`).ReplaceAllString(word, " ")
	// Remove any remaining unwanted punctuation
	word = regexp.MustCompile(`[^a-zA-Z0-9-]`).ReplaceAllString(word, "")
	return word
}
